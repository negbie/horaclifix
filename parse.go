package main

import (
	"bytes"
	"encoding/binary"
)

// ParseRecSipUDP fills the SipSet struct with
// the dataSet 258
func ParseRecSipUDP(msg []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("ParseRecSipUDP.timetaken", time.Since(t))
		}()
	*/
	i := &IPFIX{}

	i.Data.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.Data.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.Data.SIP.IntSlot = uint8(msg[8])
	i.Data.SIP.IntPort = uint8(msg[9])
	i.Data.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.Data.SIP.CallIDEnd = uint8(msg[12])
	i.Data.SIP.IPlen = binary.BigEndian.Uint16(msg[13:15])
	i.Data.SIP.VL = uint8(msg[15])
	i.Data.SIP.TOS = uint8(msg[16])
	i.Data.SIP.TLen = binary.BigEndian.Uint16(msg[17:19])
	i.Data.SIP.TID = binary.BigEndian.Uint16(msg[19:21])
	i.Data.SIP.TFlags = binary.BigEndian.Uint16(msg[21:23])
	i.Data.SIP.TTL = uint8(msg[23])
	i.Data.SIP.TProto = uint8(msg[24])
	i.Data.SIP.TPos = binary.BigEndian.Uint16(msg[25:27])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(msg[27:31])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(msg[31:35])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(msg[35:37])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(msg[37:39])
	i.Data.SIP.UDPlen = binary.BigEndian.Uint16(msg[39:41])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(msg[41:43])
	i.Data.SIP.SipMsg = msg[43:]
	return i
}

// ParseSendSipUDP fills the SipSet struct with
// the dataSet 259
func ParseSendSipUDP(msg []byte) *IPFIX {
	i := &IPFIX{}

	i.Data.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.Data.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.Data.SIP.IntSlot = uint8(msg[8])
	i.Data.SIP.IntPort = uint8(msg[9])
	i.Data.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.Data.SIP.CallIDLen = uint8(msg[12])

	pos := int(i.Data.SIP.CallIDLen)
	if pos == 0 {
		pos = 13
	} else {
		pos = pos + 13
		i.Data.SIP.CallID = []byte(msg[13:pos])
	}
	i.Data.SIP.CallIDEnd = uint8(msg[pos])
	i.Data.SIP.IPlen = binary.BigEndian.Uint16(msg[pos+1 : pos+3])
	i.Data.SIP.VL = uint8(msg[pos+3])
	i.Data.SIP.TOS = uint8(msg[pos+4])
	i.Data.SIP.TLen = binary.BigEndian.Uint16(msg[pos+5 : pos+7])
	i.Data.SIP.TID = binary.BigEndian.Uint16(msg[pos+7 : pos+9])
	i.Data.SIP.TFlags = binary.BigEndian.Uint16(msg[pos+9 : pos+11])
	i.Data.SIP.TTL = uint8(msg[pos+11])
	i.Data.SIP.TProto = uint8(msg[pos+12])
	i.Data.SIP.TPos = binary.BigEndian.Uint16(msg[pos+13 : pos+15])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(msg[pos+15 : pos+19])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(msg[pos+19 : pos+23])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(msg[pos+23 : pos+25])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(msg[pos+25 : pos+27])
	i.Data.SIP.UDPlen = binary.BigEndian.Uint16(msg[pos+27 : pos+29])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(msg[pos+29 : pos+31])
	i.Data.SIP.SipMsg = []byte(msg[pos+31:])
	return i
}

// ParseRecSipTCP fills the SipSet struct with
// the dataSet 260
func ParseRecSipTCP(msg []byte) *IPFIX {
	i := &IPFIX{}

	i.Data.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.Data.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.Data.SIP.IntSlot = uint8(msg[8])
	i.Data.SIP.IntPort = uint8(msg[9])
	i.Data.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(msg[12:16])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(msg[16:20])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(msg[20:22])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(msg[22:24])
	i.Data.SIP.Context = binary.BigEndian.Uint32(msg[24:28])
	i.Data.SIP.CallIDEnd = uint8(msg[28])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(msg[29:31])
	i.Data.SIP.SipMsg = msg[31:]
	return i
}

// ParseSendSipTCP fills the SipSet struct with
// the dataSet 261
func ParseSendSipTCP(msg []byte) *IPFIX {
	i := &IPFIX{}

	i.Data.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.Data.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.Data.SIP.IntSlot = uint8(msg[8])
	i.Data.SIP.IntPort = uint8(msg[9])
	i.Data.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(msg[12:16])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(msg[16:20])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(msg[20:22])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(msg[22:24])
	i.Data.SIP.Context = binary.BigEndian.Uint32(msg[24:28])
	i.Data.SIP.CallIDLen = uint8(msg[28])

	pos := int(i.Data.SIP.CallIDLen)
	if pos == 0 {
		pos = 29
	} else {
		pos = pos + 29
		i.Data.SIP.CallID = []byte(msg[29:pos])
	}
	i.Data.SIP.CallIDEnd = uint8(msg[pos])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(msg[pos+1 : pos+3])
	i.Data.SIP.SipMsg = []byte(msg[pos+3:])
	return i
}

// ParseQosStats fills the QosSet struct with
// the dataSet 268
func ParseQosStats(msg []byte) *IPFIX {
	var i IPFIX
	r := reader{r: bytes.NewReader(msg)}

	r.binRead(&i.Data.QOS.IncRtpBytes)
	r.binRead(&i.Data.QOS.IncRtpPackets)
	r.binRead(&i.Data.QOS.IncRtpLostPackets)
	r.binRead(&i.Data.QOS.IncRtpAvgJitter)
	r.binRead(&i.Data.QOS.IncRtpMaxJitter)

	r.binRead(&i.Data.QOS.IncRtcpBytes)
	r.binRead(&i.Data.QOS.IncRtcpPackets)
	r.binRead(&i.Data.QOS.IncRtcpLostPackets)
	r.binRead(&i.Data.QOS.IncRtcpAvgJitter)
	r.binRead(&i.Data.QOS.IncRtcpMaxJitter)
	r.binRead(&i.Data.QOS.IncRtcpAvgLat)
	r.binRead(&i.Data.QOS.IncRtcpMaxLat)

	r.binRead(&i.Data.QOS.IncrVal)
	r.binRead(&i.Data.QOS.IncMos)

	r.binRead(&i.Data.QOS.OutRtpBytes)
	r.binRead(&i.Data.QOS.OutRtpPackets)
	r.binRead(&i.Data.QOS.OutRtpLostPackets)
	r.binRead(&i.Data.QOS.OutRtpAvgJitter)
	r.binRead(&i.Data.QOS.OutRtpMaxJitter)

	r.binRead(&i.Data.QOS.OutRtcpBytes)
	r.binRead(&i.Data.QOS.OutRtcpPackets)
	r.binRead(&i.Data.QOS.OutRtcpLostPackets)
	r.binRead(&i.Data.QOS.OutRtcpAvgJitter)
	r.binRead(&i.Data.QOS.OutRtcpMaxJitter)
	r.binRead(&i.Data.QOS.OutRtcpAvgLat)
	r.binRead(&i.Data.QOS.OutRtcpMaxLat)

	r.binRead(&i.Data.QOS.OutrVal)
	r.binRead(&i.Data.QOS.OutMos)

	r.binRead(&i.Data.QOS.Type)

	r.binRead(&i.Data.QOS.CallerIncSrcIP)
	r.binRead(&i.Data.QOS.CallerIncDstIP)
	r.binRead(&i.Data.QOS.CallerIncSrcPort)
	r.binRead(&i.Data.QOS.CallerIncDstPort)

	r.binRead(&i.Data.QOS.CalleeIncSrcIP)
	r.binRead(&i.Data.QOS.CalleeIncDstIP)
	r.binRead(&i.Data.QOS.CalleeIncSrcPort)
	r.binRead(&i.Data.QOS.CalleeIncDstPort)

	r.binRead(&i.Data.QOS.CallerOutSrcIP)
	r.binRead(&i.Data.QOS.CallerOutDstIP)
	r.binRead(&i.Data.QOS.CallerOutSrcPort)
	r.binRead(&i.Data.QOS.CallerOutDstPort)

	r.binRead(&i.Data.QOS.CalleeOutSrcIP)
	r.binRead(&i.Data.QOS.CalleeOutDstIP)
	r.binRead(&i.Data.QOS.CalleeOutSrcPort)
	r.binRead(&i.Data.QOS.CalleeOutDstPort)

	r.binRead(&i.Data.QOS.CallerIntSlot)
	r.binRead(&i.Data.QOS.CallerIntPort)
	r.binRead(&i.Data.QOS.CallerIntVlan)

	r.binRead(&i.Data.QOS.CalleeIntSlot)
	r.binRead(&i.Data.QOS.CalleeIntPort)
	r.binRead(&i.Data.QOS.CalleeIntVlan)

	r.binRead(&i.Data.QOS.BeginTimeSec)
	r.binRead(&i.Data.QOS.BeginTimeMic)

	r.binRead(&i.Data.QOS.EndTimeSec)
	r.binRead(&i.Data.QOS.EndinTimeMic)

	r.binRead(&i.Data.QOS.Seperator)

	r.binRead(&i.Data.QOS.IncRealmLen)
	i.Data.QOS.IncRealm = make([]byte, i.Data.QOS.IncRealmLen)
	r.binRead(&i.Data.QOS.IncRealm)
	r.binRead(&i.Data.QOS.IncRealmEnd)

	r.binRead(&i.Data.QOS.OutRealmLen)
	i.Data.QOS.OutRealm = make([]byte, i.Data.QOS.OutRealmLen)
	r.binRead(&i.Data.QOS.OutRealm)
	r.binRead(&i.Data.QOS.OutRealmEnd)

	r.binRead(&i.Data.QOS.IncCallIDLen)
	i.Data.QOS.IncCallID = make([]byte, i.Data.QOS.IncCallIDLen)
	r.binRead(&i.Data.QOS.IncCallID)
	r.binRead(&i.Data.QOS.IncCallIDEnd)

	r.binRead(&i.Data.QOS.OutCallIDLen)
	//i.Data.QOS.OutCallID = make([]byte, r.Len())
	i.Data.QOS.IncCallID = make([]byte, i.Data.QOS.OutCallIDLen)
	r.binRead(&i.Data.QOS.OutCallID)

	return &i
}
