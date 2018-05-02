package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"

	proto "github.com/gogo/protobuf/proto"
)

var (
	hErrCnt  int
	hepVer   = []byte{0x48, 0x45, 0x50, 0x33} // "HEP3"
	hepLen   = []byte{0x00, 0x00}
	hepLen7  = []byte{0x00, 0x07}
	hepLen8  = []byte{0x00, 0x08}
	hepLen10 = []byte{0x00, 0x0a}
	chunck16 = []byte{0x00, 0x00}
	chunck32 = []byte{0x00, 0x00, 0x00, 0x00}
)

// SendHep will write the final HEP message into the buffer and send it to wire
// The first 4 bytes are the string "HEP3". The next 2 bytes are the length of the
// whole message (len("HEP3") + length of all the chucks we have). The next bytes
// are all the chuncks created by NewHEPChuncks()
// Bytes: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31......
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//        | "HEP3"|len|chuncks(0x0001|0x0002|0x0003|0x0004|0x0007|0x0008|0x0009|0x000a|0x000b|......)
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func (conn *Connections) SendHep(i *IPFIX, s string) {
	hep := &HEP{}
	var hepMsg []byte
	var err error
	if *protobuf {
		if s == "SIP" {
			hep = &HEP{
				Version:   2,
				Protocol:  uint32(i.SIP.TProto),
				SrcIP:     i.SIP.SrcIPString,
				DstIP:     i.SIP.DstIPString,
				SrcPort:   uint32(i.SIP.SrcPort),
				DstPort:   uint32(i.SIP.DstPort),
				Tsec:      i.SIP.TimeSec,
				Tmsec:     i.SIP.TimeMic,
				ProtoType: 1,
				NodeID:    uint32(*hepID),
				NodePW:    *hepPW,
				Payload:   unsafeBytesToStr(i.SIP.RawMsg),
				CID:       unsafeBytesToStr(i.SIP.CallID),
				Vlan:      uint32(i.SIP.IntVlan),
			}
		} else {
			payload := i.mapQOS()
			hep = &HEP{
				Version:   2,
				Protocol:  17,
				SrcIP:     stringIPv4(i.QOS.CallerIncSrcIP),
				DstIP:     stringIPv4(i.QOS.CallerIncDstIP),
				SrcPort:   uint32(i.QOS.CallerIncSrcPort),
				DstPort:   uint32(i.QOS.CallerIncDstPort),
				Tsec:      i.QOS.EndTimeSec,
				Tmsec:     i.QOS.EndTimeMic,
				ProtoType: 38,
				NodeID:    uint32(*hepID),
				NodePW:    *hepPW,
				Payload:   unsafeBytesToStr(payload),
				CID:       unsafeBytesToStr(i.QOS.IncCallID),
				Vlan:      uint32(i.QOS.CalleeIntVlan),
			}
		}
		hepMsg, err = proto.Marshal(hep)
		if err != nil {
			log.Printf("[WARN] <%s> %s\n", *name, err)
		}
	} else {
		hepMsg = makeHEPChuncks(i, s)
		binary.BigEndian.PutUint16(hepMsg[4:6], uint16(len(hepMsg)))
	}

	if *network == "udp" {
		_, err = conn.Homer.UDP.Write(hepMsg)
	} else if *network == "tcp" {
		_, err = conn.Homer.TCP.Write(hepMsg)
	} else if *network == "tls" {
		_, err = conn.Homer.TLS.Write(hepMsg)
	}

	if err != nil {
		hErrCnt++
		if hErrCnt%128 == 0 {
			hErrCnt = 0
			log.Printf("[WARN] <%s> %s\n", *name, err)
			if *network == "tcp" || *network == "tls" {
				reWriteHomer(conn, hepMsg)
			}
		}
	}
}

// makeHEPChuncks will construct the respective HEP chunck
func makeHEPChuncks(i *IPFIX, payloadType string) []byte {
	w := new(bytes.Buffer)
	w.Write(hepVer)
	// hepMsg length placeholder. Will be written later
	w.Write(hepLen)

	// Chunk IP protocol family (0x02=IPv4, 0x0a=IPv6)
	w.Write([]byte{0x00, 0x00, 0x00, 0x01})
	w.Write(hepLen7)
	w.WriteByte(0x02)

	// Chunk IP protocol ID (0x06=TCP, 0x11=UDP)
	w.Write([]byte{0x00, 0x00, 0x00, 0x02})
	w.Write(hepLen7)
	w.WriteByte(0x11)

	// Chunk IPv4 source address
	w.Write([]byte{0x00, 0x00, 0x00, 0x03})
	w.Write(hepLen10)
	if payloadType == "SIP" {
		binary.BigEndian.PutUint32(chunck32, i.SIP.SrcIP)
	} else {
		binary.BigEndian.PutUint32(chunck32, i.QOS.CallerIncSrcIP)
	}
	w.Write(chunck32)

	// Chunk IPv4 destination address
	w.Write([]byte{0x00, 0x00, 0x00, 0x04})
	w.Write(hepLen10)
	if payloadType == "SIP" {
		binary.BigEndian.PutUint32(chunck32, i.SIP.DstIP)
	} else {
		binary.BigEndian.PutUint32(chunck32, i.QOS.CallerIncDstIP)
	}
	w.Write(chunck32)

	// Chunk protocol source port
	w.Write([]byte{0x00, 0x00, 0x00, 0x07})
	w.Write(hepLen8)
	if payloadType == "SIP" {
		binary.BigEndian.PutUint16(chunck16, i.SIP.SrcPort)
	} else {
		binary.BigEndian.PutUint16(chunck16, i.QOS.CallerIncSrcPort)
	}
	w.Write(chunck16)

	// Chunk protocol destination port
	w.Write([]byte{0x00, 0x00, 0x00, 0x08})
	w.Write(hepLen8)
	if payloadType == "SIP" {
		binary.BigEndian.PutUint16(chunck16, i.SIP.DstPort)
	} else {
		binary.BigEndian.PutUint16(chunck16, i.QOS.CallerIncDstPort)
	}
	w.Write(chunck16)

	// Chunk unix timestamp, seconds
	w.Write([]byte{0x00, 0x00, 0x00, 0x09})
	w.Write(hepLen10)
	if payloadType == "SIP" {
		binary.BigEndian.PutUint32(chunck32, i.SIP.TimeSec)
	} else {
		binary.BigEndian.PutUint32(chunck32, i.QOS.EndTimeSec)
	}
	w.Write(chunck32)

	// Chunk unix timestamp, microseconds offset
	w.Write([]byte{0x00, 0x00, 0x00, 0x0a})
	w.Write(hepLen10)
	if payloadType == "SIP" {
		binary.BigEndian.PutUint32(chunck32, i.SIP.TimeMic)
	} else {
		binary.BigEndian.PutUint32(chunck32, i.QOS.EndTimeMic)
	}
	w.Write(chunck32)

	// Chunk protocol type (DNS, LOG, RTCP, SIP)
	w.Write([]byte{0x00, 0x00, 0x00, 0x0b})
	w.Write(hepLen7)
	switch payloadType {
	case "SIP":
		w.WriteByte(1)
	case "QOS":
		w.WriteByte(38)
	case "incQOS":
		w.WriteByte(34)
	case "outQOS":
		w.WriteByte(34)
	case "incMOS":
		w.WriteByte(35)
	case "outMOS":
		w.WriteByte(35)
	case "LOG":
		w.WriteByte(100)
	}

	// Chunk capture agent ID
	w.Write([]byte{0x00, 0x00, 0x00, 0x0c})
	w.Write(hepLen10)
	binary.BigEndian.PutUint32(chunck32, uint32(*hepID))
	w.Write(chunck32)

	// Chunk keep alive timer
	//w.Write([]byte{0x00, 0x00, 0x00, 0x0d})

	// Chunk authenticate key (plain text / TLS connection)
	w.Write([]byte{0x00, 0x00, 0x00, 0x0e})
	binary.BigEndian.PutUint16(hepLen, 6+uint16(len(*hepPW)))
	w.Write(hepLen)
	w.Write([]byte(*hepPW))

	if payloadType != "incMOS" && payloadType != "outMOS" {
		// Chunk captured packet payload
		w.Write([]byte{0x00, 0x00, 0x00, 0x0f})
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(hepLen, 6+uint16(len(i.SIP.RawMsg)))
			w.Write(hepLen)
			w.Write(i.SIP.RawMsg)
		case "QOS":
			payload := i.mapQOS()
			binary.BigEndian.PutUint16(hepLen, 6+uint16(len(payload)))
			w.Write(hepLen)
			w.Write(payload)
		case "incQOS":
			payload, err := json.Marshal(i.mapIncQOS())
			checkErr(err)
			binary.BigEndian.PutUint16(hepLen, 6+uint16(len(payload)))
			w.Write(hepLen)
			w.Write(payload)
		case "outQOS":
			payload, err := json.Marshal(i.mapOutQOS())
			checkErr(err)
			binary.BigEndian.PutUint16(hepLen, 6+uint16(len(payload)))
			w.Write(hepLen)
			w.Write(payload)
		}
	}

	// Chunk captured compressed payload (gzip/inflate)
	//w.Write([]byte{0x00,0x00, 0x00,0x10})

	if payloadType != "SIP" {
		// Chunk internal correlation id
		w.Write([]byte{0x00, 0x00, 0x00, 0x11})
		if len(i.QOS.IncCallID) > 0 {
			binary.BigEndian.PutUint16(hepLen, 6+uint16(len(i.QOS.IncCallID)))
			w.Write(hepLen)
			w.Write(i.QOS.IncCallID)
		} else {
			binary.BigEndian.PutUint16(hepLen, 6+uint16(len(i.QOS.OutCallID)))
			w.Write(hepLen)
			w.Write(i.QOS.OutCallID)
		}
	}

	if payloadType == "incMOS" || payloadType == "outMOS" {
		// Chunk MOS only
		w.Write([]byte{0x00, 0x00, 0x00, 0x20})
		w.Write(hepLen8)
		if payloadType == "incMOS" {
			binary.BigEndian.PutUint16(chunck16, uint16(i.QOS.IncMos))
		} else {
			binary.BigEndian.PutUint16(chunck16, uint16(i.QOS.OutMos))
		}
		w.Write(chunck16)
	}
	return w.Bytes()
}

func reConnectHomer(conn *Connections) error {
	conn.Homer.Lock()
	defer conn.Homer.Unlock()

	if *network == "tcp" {
		raddr := conn.Homer.TCP.RemoteAddr()
		hconn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
		if err != nil {
			return err
		}
		conn.Homer.TCP.Close()
		conn.Homer.TCP = hconn
	} else if *network == "tls" {
		hconn, err := tls.Dial("tcp", *haddr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		conn.Homer.TLS.Close()
		conn.Homer.TLS = hconn
	}
	return nil
}

func reWriteHomer(conn *Connections, b []byte) error {
	var err error
	conn.Homer.RLock()
	defer conn.Homer.RUnlock()

	if conn.Homer.disconnected {
		conn.Homer.RUnlock()
		if err := reConnectHomer(conn); err != nil {
			conn.Homer.disconnected = true
			conn.Homer.RLock()
			return err
		}
		conn.Homer.disconnected = false
		conn.Homer.RLock()
	}
	if *network == "tcp" {
		_, err = conn.Homer.TCP.Write(b)
	} else if *network == "tls" {
		_, err = conn.Homer.TLS.Write(b)
	}
	if err == nil {
		return err
	}
	conn.Homer.disconnected = true
	return err
}
