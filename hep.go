package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	chuncks := i.NewHEPChuncks(s)
	hepMsg := make([]byte, len(chuncks)+6)
	copy(hepMsg[6:], chuncks)
	binary.BigEndian.PutUint32(hepMsg[:4], uint32(0x48455033))   // ASCII "HEP3"
	binary.BigEndian.PutUint16(hepMsg[4:6], uint16(len(hepMsg))) // Total length
	_, err := conn.Homer.Write(hepMsg)
	checkErr(err)
}

// MakeChunck will construct the respective HEP chunck
func (i *IPFIX) MakeChunck(chunckVen uint16, chunckType uint16, payloadType string) []byte {
	var chunck []byte
	switch chunckType {
	// Chunk IP protocol family (0x02=IPv4)
	case 0x0001:
		chunck = make([]byte, 6+1)
		chunck[6] = 0x02

	// Chunk IP protocol ID (0x11=UDP)
	case 0x0002:
		chunck = make([]byte, 6+1)
		chunck[6] = 0x11

	// Chunk IPv4 source address
	case 0x0003:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.SrcIP)
		default:
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.CallerIncSrcIP)
		}

	// Chunk IPv4 destination address
	case 0x0004:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.DstIP)
		default:
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.CallerIncDstIP)
		}

	// Chunk IPv6 source address
	// case 0x0005:

	// Chunk IPv6 destination address
	// case 0x0006:

	// Chunk protocol source port
	case 0x0007:
		chunck = make([]byte, 6+2)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.SIP.SrcPort)
		default:
			binary.BigEndian.PutUint16(chunck[6:], i.Data.QOS.CallerIncSrcPort)
		}

	// Chunk destination source port
	case 0x0008:
		chunck = make([]byte, 6+2)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.SIP.DstPort)
		default:
			binary.BigEndian.PutUint16(chunck[6:], i.Data.QOS.CallerIncDstPort)
		}

	// Chunk unix timestamp, seconds
	case 0x0009:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.TimeSec)
		default:
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.EndTimeSec)
		}

	// Chunk unix timestamp, microseconds offset
	case 0x000a:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.TimeMic)
		default:
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.EndinTimeMic)
		}

	// Chunk protocol type (SIP/H323/RTP/MGCP/M2UA)
	case 0x000b:
		chunck = make([]byte, 6+1)
		switch payloadType {
		case "allQOS":
			chunck[6] = 38
		case "incQOS":
			chunck[6] = 34
		case "outQOS":
			chunck[6] = 34
		case "incMOS":
			chunck[6] = 35
		case "outMOS":
			chunck[6] = 35
		case "LOG":
			chunck[6] = 100
		default:
			chunck[6] = 1 // SIP
		}

	// Chunk capture agent ID
	case 0x000c:
		chunck = make([]byte, 6+4)
		binary.BigEndian.PutUint32(chunck[6:], 0x00000BEE)

	// case 0x000d:
	// Chunk keep alive timer

	// Chunk authenticate key (plain text / TLS connection)
	case 0x000e:
		chunck = make([]byte, len(*hepPw)+6)
		copy(chunck[6:], *hepPw)

	// Chunk captured packet payload
	case 0x000f:
		switch payloadType {
		case "allQOS":
			payload, _ := json.Marshal(i.mapAllQOS())
			chunck = make([]byte, len(payload)+6)
			copy(chunck[6:], payload)
		case "incQOS":
			payload, _ := json.Marshal(i.mapIncQOS())
			chunck = make([]byte, len(payload)+6)
			copy(chunck[6:], payload)
		case "outQOS":
			payload, _ := json.Marshal(i.mapOutQOS())
			chunck = make([]byte, len(payload)+6)
			copy(chunck[6:], payload)
		default:
			chunck = make([]byte, len(i.Data.SIP.SipMsg)+6)
			copy(chunck[6:], i.Data.SIP.SipMsg)
		}

	// case 0x0010:
	// Chunk captured compressed payload (gzip/inflate)

	// Chunk internal correlation id
	case 0x0011:
		if len(i.Data.QOS.IncCallID) > 0 {
			chunck = make([]byte, len(i.Data.QOS.IncCallID)+6)
			copy(chunck[6:], i.Data.QOS.IncCallID)
		} else {
			chunck = make([]byte, len(i.Data.QOS.OutCallID)+6)
			copy(chunck[6:], i.Data.QOS.OutCallID)
		}

	case 0x0020:
		chunck = make([]byte, 6+2)
		switch payloadType {
		case "incMOS":
			binary.BigEndian.PutUint16(chunck[6:], uint16(i.Data.QOS.IncMos))
		case "outMOS":
			binary.BigEndian.PutUint16(chunck[6:], uint16(i.Data.QOS.OutMos))
		}
	}

	binary.BigEndian.PutUint16(chunck[:2], chunckVen)
	binary.BigEndian.PutUint16(chunck[2:4], chunckType)
	binary.BigEndian.PutUint16(chunck[4:6], uint16(len(chunck)))
	return chunck
}

// NewHEPChuncks will fill a buffer with all the chuncks
func (i *IPFIX) NewHEPChuncks(s string) []byte {
	buf := new(bytes.Buffer)

	buf.Write(i.MakeChunck(0x0000, 0x0001, s))
	buf.Write(i.MakeChunck(0x0000, 0x0002, s))
	buf.Write(i.MakeChunck(0x0000, 0x0003, s))
	buf.Write(i.MakeChunck(0x0000, 0x0004, s))
	buf.Write(i.MakeChunck(0x0000, 0x0007, s))
	buf.Write(i.MakeChunck(0x0000, 0x0008, s))
	buf.Write(i.MakeChunck(0x0000, 0x0009, s))
	buf.Write(i.MakeChunck(0x0000, 0x000a, s))
	buf.Write(i.MakeChunck(0x0000, 0x000b, s))
	buf.Write(i.MakeChunck(0x0000, 0x000c, s))
	buf.Write(i.MakeChunck(0x0000, 0x000e, s))
	if s != "incMOS" || s != "outMOS" {
		buf.Write(i.MakeChunck(0x0000, 0x000f, s))
	}
	if s != "SIP" {
		buf.Write(i.MakeChunck(0x0000, 0x0011, s))
	}
	if s == "incMOS" || s == "outMOS" {
		buf.Write(i.MakeChunck(0x0000, 0x0020, s))
	}
	return buf.Bytes()
}
