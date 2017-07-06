package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// NewHeader fills the IpfixHeader struct with structured binary data from r
func NewHeader(header []byte) *IPFIX {
	var ipfix IPFIX

	ipfix.Header.Version = binary.BigEndian.Uint16(header[:2])
	ipfix.Header.Length = binary.BigEndian.Uint16(header[2:4])
	ipfix.Header.ExportTime = binary.BigEndian.Uint32(header[4:8])
	ipfix.Header.SeqNum = binary.BigEndian.Uint32(header[8:12])
	ipfix.Header.ObservationID = binary.BigEndian.Uint32(header[12:16])
	ipfix.SetHeader.ID = binary.BigEndian.Uint16(header[16:18])
	ipfix.SetHeader.Length = binary.BigEndian.Uint16(header[18:20])

	return &ipfix
}

func NewRecSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
	err = binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IPlen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.VL)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TOS)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TLen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TID)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TFlags)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TTL)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TProto)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TPos)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.UDPlen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewRecSipUDP binary.Read failed:", err, r)
	}

	return &ipfix
}

func NewSendSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
	err = binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDLen)
	if ipfix.Data.SIP.CallIDLen != 0xff {
		ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
		err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
		err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	}
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IPlen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.VL)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TOS)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TLen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TID)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TFlags)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TTL)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TProto)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TPos)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.UDPlen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewSendSipUDP binary.Read failed:", err, r)
	}

	return &ipfix
}

func NewRecSipTCP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
	err = binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.Context)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewRecSipTCP binary.Read failed:", err, r)
	}

	return &ipfix
}

func NewSendSipTCP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
	err = binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.Context)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDLen)
	if ipfix.Data.SIP.CallIDLen != 0xff {
		ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
		err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
		err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	}
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)
	if err != nil {
		log.Println("NewSendSipTCP binary.Read failed:", err, r)
	}

	return &ipfix
}
