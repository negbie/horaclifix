package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// NewRecSipUDP fills the SipSet struct with
// the dataSet 258
func NewRecSipUDP(header []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("NewRecSipUDP.timetaken", time.Since(t))
		}()
	*/
	var i IPFIX

	i.Header.Version = binary.BigEndian.Uint16(header[:2])
	i.Header.Length = binary.BigEndian.Uint16(header[2:4])
	i.Header.ExportTime = binary.BigEndian.Uint32(header[4:8])
	i.Header.SeqNum = binary.BigEndian.Uint32(header[8:12])
	i.Header.ObservationID = binary.BigEndian.Uint32(header[12:16])
	i.SetHeader.ID = binary.BigEndian.Uint16(header[16:18])
	i.SetHeader.Length = binary.BigEndian.Uint16(header[18:20])

	i.Data.SIP.TimeSec = binary.BigEndian.Uint32(header[20:24])
	i.Data.SIP.TimeMic = binary.BigEndian.Uint32(header[24:28])
	i.Data.SIP.IntSlot = uint8(header[28])
	i.Data.SIP.IntPort = uint8(header[29])
	i.Data.SIP.IntVlan = binary.BigEndian.Uint16(header[30:32])
	i.Data.SIP.CallIDEnd = uint8(header[32])
	i.Data.SIP.IPlen = binary.BigEndian.Uint16(header[33:35])
	i.Data.SIP.VL = uint8(header[35])
	i.Data.SIP.TOS = uint8(header[36])
	i.Data.SIP.TLen = binary.BigEndian.Uint16(header[37:39])
	i.Data.SIP.TID = binary.BigEndian.Uint16(header[39:41])
	i.Data.SIP.TFlags = binary.BigEndian.Uint16(header[41:43])
	i.Data.SIP.TTL = uint8(header[43])
	i.Data.SIP.TProto = uint8(header[44])
	i.Data.SIP.TPos = binary.BigEndian.Uint16(header[45:47])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[47:51])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(header[51:55])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(header[55:57])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[57:59])
	i.Data.SIP.UDPlen = binary.BigEndian.Uint16(header[59:61])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[61:63])
	i.Data.SIP.SipMsg = header[63:]
	return &i
}

// NewSendSipUDP fills the SipSet struct with
// the dataSet 259
func NewSendSipUDP(header []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("NewSendSipUDP.timetaken", time.Since(t))
		}()
	*/
	var i IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &i.Header)
	binary.Read(r, binary.BigEndian, &i.SetHeader)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.CallIDLen)
	if i.Data.SIP.CallIDLen != 0xff {
		i.Data.SIP.CallID = make([]byte, i.Data.SIP.CallIDLen)
		binary.Read(r, binary.BigEndian, &i.Data.SIP.CallID)
		binary.Read(r, binary.BigEndian, &i.Data.SIP.CallIDEnd)
	}
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IPlen)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.VL)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TOS)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TLen)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TID)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TFlags)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TTL)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TProto)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TPos)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.UDPlen)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.MsgLen)
	i.Data.SIP.SipMsg = make([]byte, r.Len())
	err := binary.Read(r, binary.BigEndian, &i.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewSendSipUDP binary.Read failed:", err, r)
	}
	return &i
}

// NewRecSipTCP fills the SipSet struct with
// the dataSet 260
func NewRecSipTCP(header []byte) *IPFIX {
	var i IPFIX

	i.Header.Version = binary.BigEndian.Uint16(header[:2])
	i.Header.Length = binary.BigEndian.Uint16(header[2:4])
	i.Header.ExportTime = binary.BigEndian.Uint32(header[4:8])
	i.Header.SeqNum = binary.BigEndian.Uint32(header[8:12])
	i.Header.ObservationID = binary.BigEndian.Uint32(header[12:16])
	i.SetHeader.ID = binary.BigEndian.Uint16(header[16:18])
	i.SetHeader.Length = binary.BigEndian.Uint16(header[18:20])

	i.Data.SIP.TimeSec = binary.BigEndian.Uint32(header[20:24])
	i.Data.SIP.TimeMic = binary.BigEndian.Uint32(header[24:28])
	i.Data.SIP.IntSlot = uint8(header[28])
	i.Data.SIP.IntPort = uint8(header[29])
	i.Data.SIP.IntVlan = binary.BigEndian.Uint16(header[30:32])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(header[32:36])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[36:40])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(header[40:42])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[42:44])
	i.Data.SIP.Context = binary.BigEndian.Uint32(header[44:48])
	i.Data.SIP.CallIDEnd = uint8(header[48])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[49:51])
	i.Data.SIP.SipMsg = header[51:]
	return &i
}

// NewSendSipTCP fills the SipSet struct with
// the dataSet 261
func NewSendSipTCP(header []byte) *IPFIX {
	var i IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &i.Header)
	binary.Read(r, binary.BigEndian, &i.SetHeader)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.Context)
	binary.Read(r, binary.BigEndian, &i.Data.SIP.CallIDLen)
	if i.Data.SIP.CallIDLen != 0xff {
		i.Data.SIP.CallID = make([]byte, i.Data.SIP.CallIDLen)
		binary.Read(r, binary.BigEndian, &i.Data.SIP.CallID)
		binary.Read(r, binary.BigEndian, &i.Data.SIP.CallIDEnd)
	}
	binary.Read(r, binary.BigEndian, &i.Data.SIP.MsgLen)
	i.Data.SIP.SipMsg = make([]byte, r.Len())
	err := binary.Read(r, binary.BigEndian, &i.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewSendSipTCP binary.Read failed:", err, r)
	}
	return &i
}
