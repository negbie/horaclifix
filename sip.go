package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func NewRecSipUDP(header []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("NewRecSipUDP.timetaken", time.Since(t))
		}()
	*/
	var ipfix IPFIX

	ipfix.Header.Version = binary.BigEndian.Uint16(header[:2])
	ipfix.Header.Length = binary.BigEndian.Uint16(header[2:4])
	ipfix.Header.ExportTime = binary.BigEndian.Uint32(header[4:8])
	ipfix.Header.SeqNum = binary.BigEndian.Uint32(header[8:12])
	ipfix.Header.ObservationID = binary.BigEndian.Uint32(header[12:16])
	ipfix.SetHeader.ID = binary.BigEndian.Uint16(header[16:18])
	ipfix.SetHeader.Length = binary.BigEndian.Uint16(header[18:20])

	ipfix.Data.SIP.TimeSec = binary.BigEndian.Uint32(header[20:24])
	ipfix.Data.SIP.TimeMic = binary.BigEndian.Uint32(header[24:28])
	ipfix.Data.SIP.IntSlot = uint8(header[28])
	ipfix.Data.SIP.IntPort = uint8(header[29])
	ipfix.Data.SIP.IntVlan = binary.BigEndian.Uint16(header[30:32])
	ipfix.Data.SIP.CallIDEnd = uint8(header[32])
	ipfix.Data.SIP.IPlen = binary.BigEndian.Uint16(header[33:35])
	ipfix.Data.SIP.VL = uint8(header[35])
	ipfix.Data.SIP.TOS = uint8(header[36])
	ipfix.Data.SIP.TLen = binary.BigEndian.Uint16(header[37:39])
	ipfix.Data.SIP.TID = binary.BigEndian.Uint16(header[39:41])
	ipfix.Data.SIP.TFlags = binary.BigEndian.Uint16(header[41:43])
	ipfix.Data.SIP.TTL = uint8(header[43])
	ipfix.Data.SIP.TProto = uint8(header[44])
	ipfix.Data.SIP.TPos = binary.BigEndian.Uint16(header[45:47])
	ipfix.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[47:51])
	ipfix.Data.SIP.DstIP = binary.BigEndian.Uint32(header[51:55])
	ipfix.Data.SIP.DstPort = binary.BigEndian.Uint16(header[55:57])
	ipfix.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[57:59])
	ipfix.Data.SIP.UDPlen = binary.BigEndian.Uint16(header[59:61])
	ipfix.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[61:63])
	ipfix.Data.SIP.SipMsg = header[63:]

	return &ipfix
}

func NewSendSipUDP(header []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("NewSendSipUDP.timetaken", time.Since(t))
		}()
	*/
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDLen)
	if ipfix.Data.SIP.CallIDLen != 0xff {
		ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
		binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
		binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	}
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.VL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TOS)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TFlags)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TTL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TProto)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TPos)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.UDPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	err := binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewSendSipUDP binary.Read failed:", err, r)
	}

	return &ipfix
}

func NewRecSipTCP(header []byte) *IPFIX {
	var ipfix IPFIX

	ipfix.Header.Version = binary.BigEndian.Uint16(header[:2])
	ipfix.Header.Length = binary.BigEndian.Uint16(header[2:4])
	ipfix.Header.ExportTime = binary.BigEndian.Uint32(header[4:8])
	ipfix.Header.SeqNum = binary.BigEndian.Uint32(header[8:12])
	ipfix.Header.ObservationID = binary.BigEndian.Uint32(header[12:16])
	ipfix.SetHeader.ID = binary.BigEndian.Uint16(header[16:18])
	ipfix.SetHeader.Length = binary.BigEndian.Uint16(header[18:20])

	ipfix.Data.SIP.TimeSec = binary.BigEndian.Uint32(header[20:24])
	ipfix.Data.SIP.TimeMic = binary.BigEndian.Uint32(header[24:28])
	ipfix.Data.SIP.IntSlot = uint8(header[28])
	ipfix.Data.SIP.IntPort = uint8(header[29])
	ipfix.Data.SIP.IntVlan = binary.BigEndian.Uint16(header[30:32])
	ipfix.Data.SIP.DstIP = binary.BigEndian.Uint32(header[32:36])
	ipfix.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[36:40])
	ipfix.Data.SIP.DstPort = binary.BigEndian.Uint16(header[40:42])
	ipfix.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[42:44])
	ipfix.Data.SIP.Context = binary.BigEndian.Uint32(header[44:48])
	ipfix.Data.SIP.CallIDEnd = uint8(header[48])
	ipfix.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[49:51])
	ipfix.Data.SIP.SipMsg = header[51:]

	return &ipfix

}

func NewSendSipTCP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.Context)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDLen)
	if ipfix.Data.SIP.CallIDLen != 0xff {
		ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
		binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
		binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	}
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	err := binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewSendSipTCP binary.Read failed:", err, r)
	}

	return &ipfix
}
