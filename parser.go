package main

import (
	"bytes"
	"encoding/binary"
	"net"
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
	ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
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
	ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewCallQualityStats(header []byte) *CallStats {
	var ipfix IPFIX
	var callstats CallStats
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpBytes)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &callstats.IncRtcpBytes)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpMaxLat)

	binary.Read(r, binary.BigEndian, &callstats.IncrVal)
	binary.Read(r, binary.BigEndian, &callstats.IncMos)

	binary.Read(r, binary.BigEndian, &callstats.OutRtpBytes)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &callstats.OutRtcpBytes)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpMaxLat)
	binary.Read(r, binary.BigEndian, &callstats.OutrVal)
	binary.Read(r, binary.BigEndian, &callstats.OutMos)

	binary.Read(r, binary.BigEndian, &callstats.Type)

	binary.Read(r, binary.BigEndian, &callstats.CallerIncSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerIncDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerIncSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CallerIncDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CalleeIncSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIncDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIncSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIncDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CallerOutSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerOutDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerOutSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CallerOutDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CalleeOutSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeOutDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeOutSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CalleeOutDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CallerIntSlot)
	binary.Read(r, binary.BigEndian, &callstats.CallerIntPort)
	binary.Read(r, binary.BigEndian, &callstats.CallerIntVlan)

	binary.Read(r, binary.BigEndian, &callstats.CalleeIntSlot)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIntPort)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIntVlan)

	binary.Read(r, binary.BigEndian, &callstats.BeginTimeSec)
	binary.Read(r, binary.BigEndian, &callstats.BeginTimeMic)

	binary.Read(r, binary.BigEndian, &callstats.EndTimeSec)
	binary.Read(r, binary.BigEndian, &callstats.EndinTimeMic)

	binary.Read(r, binary.BigEndian, &callstats.Seperator)

	binary.Read(r, binary.BigEndian, &callstats.IncRealmLen)
	callstats.IncRealm = make([]byte, callstats.IncRealmLen)
	binary.Read(r, binary.BigEndian, &callstats.IncRealm)
	binary.Read(r, binary.BigEndian, &callstats.IncRealmEnd)

	binary.Read(r, binary.BigEndian, &callstats.OutRealmLen)
	callstats.OutRealm = make([]byte, callstats.OutRealmLen)
	binary.Read(r, binary.BigEndian, &callstats.OutRealm)
	binary.Read(r, binary.BigEndian, &callstats.OutRealmEnd)

	binary.Read(r, binary.BigEndian, &callstats.IncCallIDLen)
	callstats.IncCallID = make([]byte, callstats.IncCallIDLen)
	binary.Read(r, binary.BigEndian, &callstats.IncCallID)
	binary.Read(r, binary.BigEndian, &callstats.IncCallIDEnd)

	binary.Read(r, binary.BigEndian, &callstats.OutCallIDLen)
	callstats.OutCallID = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &callstats.OutCallID)

	return &callstats
}

// HandShake writes the binary Handshake representation into the buffer
func (ipfix *IPFIX) SendHandShake() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, &ipfix.Header)
	binary.Write(b, binary.BigEndian, &ipfix.SetHeader)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake)
	return b.Bytes()
}

// NewHEPMsg writes the binary HEP representation into the buffer
func NewHEPMsg(msg []byte) []byte {

	b := bytes.NewBuffer(make([]byte, 6))
	binary.Write(b, binary.BigEndian, msg)
	packet := b.Bytes()
	binary.BigEndian.PutUint32(packet, uint32(0x48455033)) // ASCII "HEP3"
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))
	return packet
}

// NewHEPChunck constructs the HEP chunck
func (ipfix *IPFIX) NewHEPChunck(ChunckVen uint16, ChunckType uint16) []byte {

	b := bytes.NewBuffer(make([]byte, 6))
	switch ChunckType {
	case 0x0001:
		binary.Write(b, binary.BigEndian, 0x02)

	case 0x0002:
		binary.Write(b, binary.BigEndian, 0x11)

	case 0x0003:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SrcIP)

	case 0x0004:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.DstIP)

	case 0x0007:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SrcPort)

	case 0x0008:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.DstPort)

	case 0x0009:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.TimeSec)

	case 0x000a:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.TimeMic)

	case 0x000b:
		binary.Write(b, binary.BigEndian, 0x01)

	case 0x000c:
		binary.Write(b, binary.BigEndian, 0x000007D1)

	case 0x000f:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	}
	packet := b.Bytes()
	binary.BigEndian.PutUint16(packet, ChunckVen)
	binary.BigEndian.PutUint16(packet[2:4], ChunckType)
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))

	return packet
}

// SendHEP sends the HEP message
func SendHEP(p *IPFIX, c net.Conn) {
	bhep := new(bytes.Buffer)
	bhep.Write(p.NewHEPChunck(0x0000, 0x0001))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0002))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0003))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0004))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0007))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0008))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0009))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000a))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000b))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000c))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000f))

	//fmt.Printf("%s\n", hex.Dump(bhep.Bytes()))
	c.Write(NewHEPMsg(bhep.Bytes()))
}
