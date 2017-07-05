package main

import (
	"bytes"
	"encoding/binary"
)

// NewHeader fills the IpfixHeader struct with structured binary data from r
func NewHeader(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	return &ipfix
}

// NewHandShake fills the IPFIX struct with structured binary data from r
func NewHandShake(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake)
	return &ipfix
}

func NewRecSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
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
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewSendSipUDP(header []byte) *IPFIX {
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
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewRecSipTCP(header []byte) *IPFIX {
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
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

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
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewQosStats(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpBytes)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpLostPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpBytes)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpMaxLat)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncrVal)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncMos)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpBytes)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpLostPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpBytes)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpMaxLat)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutrVal)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutMos)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.Type)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncSrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncDstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncSrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncDstPort)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncSrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncDstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncSrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncDstPort)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutSrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutDstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutSrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutDstPort)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutSrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutDstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutSrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutDstPort)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIntVlan)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIntVlan)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.BeginTimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.BeginTimeMic)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.EndTimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.EndinTimeMic)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.Seperator)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRealmLen)
	ipfix.Data.QOS.IncRealm = make([]byte, ipfix.Data.QOS.IncRealmLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRealm)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRealmEnd)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRealmLen)
	ipfix.Data.QOS.OutRealm = make([]byte, ipfix.Data.QOS.OutRealmLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRealm)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRealmEnd)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncCallIDLen)
	ipfix.Data.QOS.IncCallID = make([]byte, ipfix.Data.QOS.IncCallIDLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncCallID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncCallIDEnd)

	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutCallIDLen)
	ipfix.Data.QOS.OutCallID = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutCallID)

	return &ipfix
}
