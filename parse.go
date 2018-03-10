package main

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/negbie/sipparser"
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

	i.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.SIP.IntSlot = uint8(msg[8])
	i.SIP.IntPort = uint8(msg[9])
	i.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.SIP.CallIDEnd = uint8(msg[12])
	i.SIP.IPlen = binary.BigEndian.Uint16(msg[13:15])
	i.SIP.VL = uint8(msg[15])
	i.SIP.TOS = uint8(msg[16])
	i.SIP.TLen = binary.BigEndian.Uint16(msg[17:19])
	i.SIP.TID = binary.BigEndian.Uint16(msg[19:21])
	i.SIP.TFlags = binary.BigEndian.Uint16(msg[21:23])
	i.SIP.TTL = uint8(msg[23])
	i.SIP.TProto = uint8(msg[24])
	i.SIP.TPos = binary.BigEndian.Uint16(msg[25:27])
	i.SIP.SrcIP = binary.BigEndian.Uint32(msg[27:31])
	i.SIP.DstIP = binary.BigEndian.Uint32(msg[31:35])
	i.SIP.SrcPort = binary.BigEndian.Uint16(msg[35:37])
	i.SIP.DstPort = binary.BigEndian.Uint16(msg[37:39])
	i.SIP.UDPlen = binary.BigEndian.Uint16(msg[39:41])
	i.SIP.MsgLen = binary.BigEndian.Uint16(msg[41:43])
	i.SIP.RawMsg = msg[43:]

	if *gaddr != "" {
		err := i.parseSIP()
		checkErr(err)
	}
	return i
}

// ParseSendSipUDP fills the SipSet struct with
// the dataSet 259
func ParseSendSipUDP(msg []byte) *IPFIX {
	i := &IPFIX{}

	i.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.SIP.IntSlot = uint8(msg[8])
	i.SIP.IntPort = uint8(msg[9])
	i.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.SIP.CallIDLen = uint8(msg[12])

	pos := int(i.SIP.CallIDLen)
	if pos == 0 {
		pos = 13
	} else {
		pos = pos + 13
		i.SIP.CallID = []byte(msg[13:pos])
	}
	i.SIP.CallIDEnd = uint8(msg[pos])
	i.SIP.IPlen = binary.BigEndian.Uint16(msg[pos+1 : pos+3])
	i.SIP.VL = uint8(msg[pos+3])
	i.SIP.TOS = uint8(msg[pos+4])
	i.SIP.TLen = binary.BigEndian.Uint16(msg[pos+5 : pos+7])
	i.SIP.TID = binary.BigEndian.Uint16(msg[pos+7 : pos+9])
	i.SIP.TFlags = binary.BigEndian.Uint16(msg[pos+9 : pos+11])
	i.SIP.TTL = uint8(msg[pos+11])
	i.SIP.TProto = uint8(msg[pos+12])
	i.SIP.TPos = binary.BigEndian.Uint16(msg[pos+13 : pos+15])
	i.SIP.SrcIP = binary.BigEndian.Uint32(msg[pos+15 : pos+19])
	i.SIP.DstIP = binary.BigEndian.Uint32(msg[pos+19 : pos+23])
	i.SIP.SrcPort = binary.BigEndian.Uint16(msg[pos+23 : pos+25])
	i.SIP.DstPort = binary.BigEndian.Uint16(msg[pos+25 : pos+27])
	i.SIP.UDPlen = binary.BigEndian.Uint16(msg[pos+27 : pos+29])
	i.SIP.MsgLen = binary.BigEndian.Uint16(msg[pos+29 : pos+31])
	i.SIP.RawMsg = []byte(msg[pos+31:])

	if *gaddr != "" {
		err := i.parseSIP()
		checkErr(err)
	}
	return i
}

// ParseRecSipTCP fills the SipSet struct with
// the dataSet 260
func ParseRecSipTCP(msg []byte) *IPFIX {
	i := &IPFIX{}

	i.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.SIP.IntSlot = uint8(msg[8])
	i.SIP.IntPort = uint8(msg[9])
	i.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.SIP.DstIP = binary.BigEndian.Uint32(msg[12:16])
	i.SIP.SrcIP = binary.BigEndian.Uint32(msg[16:20])
	i.SIP.DstPort = binary.BigEndian.Uint16(msg[20:22])
	i.SIP.SrcPort = binary.BigEndian.Uint16(msg[22:24])
	i.SIP.Context = binary.BigEndian.Uint32(msg[24:28])
	i.SIP.CallIDEnd = uint8(msg[28])
	i.SIP.MsgLen = binary.BigEndian.Uint16(msg[29:31])
	i.SIP.RawMsg = msg[31:]

	if *gaddr != "" {
		err := i.parseSIP()
		checkErr(err)
	}
	return i
}

// ParseSendSipTCP fills the SipSet struct with
// the dataSet 261
func ParseSendSipTCP(msg []byte) *IPFIX {
	i := &IPFIX{}

	i.SIP.TimeSec = binary.BigEndian.Uint32(msg[:4])
	i.SIP.TimeMic = binary.BigEndian.Uint32(msg[4:8])
	i.SIP.IntSlot = uint8(msg[8])
	i.SIP.IntPort = uint8(msg[9])
	i.SIP.IntVlan = binary.BigEndian.Uint16(msg[10:12])
	i.SIP.DstIP = binary.BigEndian.Uint32(msg[12:16])
	i.SIP.SrcIP = binary.BigEndian.Uint32(msg[16:20])
	i.SIP.DstPort = binary.BigEndian.Uint16(msg[20:22])
	i.SIP.SrcPort = binary.BigEndian.Uint16(msg[22:24])
	i.SIP.Context = binary.BigEndian.Uint32(msg[24:28])
	i.SIP.CallIDLen = uint8(msg[28])

	pos := int(i.SIP.CallIDLen)
	if pos == 0 {
		pos = 29
	} else {
		pos = pos + 29
		i.SIP.CallID = []byte(msg[29:pos])
	}
	i.SIP.CallIDEnd = uint8(msg[pos])
	i.SIP.MsgLen = binary.BigEndian.Uint16(msg[pos+1 : pos+3])
	i.SIP.RawMsg = []byte(msg[pos+3:])

	if *gaddr != "" {
		err := i.parseSIP()
		checkErr(err)
	}
	return i
}

// ParseQosStats fills the QosSet struct with
// the dataSet 268
func ParseQosStats(msg []byte) *IPFIX {
	var i IPFIX
	r := reader{r: bytes.NewReader(msg)}

	r.binRead(&i.QOS.IncRtpBytes)
	r.binRead(&i.QOS.IncRtpPackets)
	r.binRead(&i.QOS.IncRtpLostPackets)
	r.binRead(&i.QOS.IncRtpAvgJitter)
	r.binRead(&i.QOS.IncRtpMaxJitter)

	r.binRead(&i.QOS.IncRtcpBytes)
	r.binRead(&i.QOS.IncRtcpPackets)
	r.binRead(&i.QOS.IncRtcpLostPackets)
	r.binRead(&i.QOS.IncRtcpAvgJitter)
	r.binRead(&i.QOS.IncRtcpMaxJitter)
	r.binRead(&i.QOS.IncRtcpAvgLat)
	r.binRead(&i.QOS.IncRtcpMaxLat)

	r.binRead(&i.QOS.IncrVal)
	r.binRead(&i.QOS.IncMos)

	r.binRead(&i.QOS.OutRtpBytes)
	r.binRead(&i.QOS.OutRtpPackets)
	r.binRead(&i.QOS.OutRtpLostPackets)
	r.binRead(&i.QOS.OutRtpAvgJitter)
	r.binRead(&i.QOS.OutRtpMaxJitter)

	r.binRead(&i.QOS.OutRtcpBytes)
	r.binRead(&i.QOS.OutRtcpPackets)
	r.binRead(&i.QOS.OutRtcpLostPackets)
	r.binRead(&i.QOS.OutRtcpAvgJitter)
	r.binRead(&i.QOS.OutRtcpMaxJitter)
	r.binRead(&i.QOS.OutRtcpAvgLat)
	r.binRead(&i.QOS.OutRtcpMaxLat)

	r.binRead(&i.QOS.OutrVal)
	r.binRead(&i.QOS.OutMos)

	r.binRead(&i.QOS.Type)

	r.binRead(&i.QOS.CallerIncSrcIP)
	r.binRead(&i.QOS.CallerIncDstIP)
	r.binRead(&i.QOS.CallerIncSrcPort)
	r.binRead(&i.QOS.CallerIncDstPort)

	r.binRead(&i.QOS.CalleeIncSrcIP)
	r.binRead(&i.QOS.CalleeIncDstIP)
	r.binRead(&i.QOS.CalleeIncSrcPort)
	r.binRead(&i.QOS.CalleeIncDstPort)

	r.binRead(&i.QOS.CallerOutSrcIP)
	r.binRead(&i.QOS.CallerOutDstIP)
	r.binRead(&i.QOS.CallerOutSrcPort)
	r.binRead(&i.QOS.CallerOutDstPort)

	r.binRead(&i.QOS.CalleeOutSrcIP)
	r.binRead(&i.QOS.CalleeOutDstIP)
	r.binRead(&i.QOS.CalleeOutSrcPort)
	r.binRead(&i.QOS.CalleeOutDstPort)

	r.binRead(&i.QOS.CallerIntSlot)
	r.binRead(&i.QOS.CallerIntPort)
	r.binRead(&i.QOS.CallerIntVlan)

	r.binRead(&i.QOS.CalleeIntSlot)
	r.binRead(&i.QOS.CalleeIntPort)
	r.binRead(&i.QOS.CalleeIntVlan)

	r.binRead(&i.QOS.BeginTimeSec)
	r.binRead(&i.QOS.BeginTimeMic)

	r.binRead(&i.QOS.EndTimeSec)
	r.binRead(&i.QOS.EndinTimeMic)

	r.binRead(&i.QOS.Seperator)

	r.binRead(&i.QOS.IncRealmLen)
	i.QOS.IncRealm = make([]byte, i.QOS.IncRealmLen)
	r.binRead(&i.QOS.IncRealm)
	r.binRead(&i.QOS.IncRealmEnd)

	r.binRead(&i.QOS.OutRealmLen)
	i.QOS.OutRealm = make([]byte, i.QOS.OutRealmLen)
	r.binRead(&i.QOS.OutRealm)
	r.binRead(&i.QOS.OutRealmEnd)

	r.binRead(&i.QOS.IncCallIDLen)
	i.QOS.IncCallID = make([]byte, i.QOS.IncCallIDLen)
	r.binRead(&i.QOS.IncCallID)
	r.binRead(&i.QOS.IncCallIDEnd)

	r.binRead(&i.QOS.OutCallIDLen)
	i.QOS.OutCallID = make([]byte, i.QOS.OutCallIDLen)
	r.binRead(&i.QOS.OutCallID)

	return &i
}

func (i *IPFIX) parseSIP() error {

	i.SIP.SipMsg = sipparser.ParseMsg(string(i.SIP.RawMsg))

	if i.SIP.SipMsg.StartLine == nil {
		i.SIP.SipMsg.StartLine = new(sipparser.StartLine)
	}
	if i.SIP.SipMsg.StartLine.Method == "" {
		i.SIP.SipMsg.StartLine.Method = i.SIP.SipMsg.StartLine.Resp
	}
	if i.SIP.SipMsg.StartLine.URI == nil {
		i.SIP.SipMsg.StartLine.URI = new(sipparser.URI)
	}
	if i.SIP.SipMsg.From == nil {
		i.SIP.SipMsg.From = new(sipparser.From)
	}
	if i.SIP.SipMsg.From.URI == nil {
		i.SIP.SipMsg.From.URI = new(sipparser.URI)
	}
	if i.SIP.SipMsg.To == nil {
		i.SIP.SipMsg.To = new(sipparser.From)
	}
	if i.SIP.SipMsg.To.URI == nil {
		i.SIP.SipMsg.To.URI = new(sipparser.URI)
	}
	if i.SIP.SipMsg.Contact == nil {
		i.SIP.SipMsg.Contact = new(sipparser.From)
	}
	if i.SIP.SipMsg.Contact.URI == nil {
		i.SIP.SipMsg.Contact.URI = new(sipparser.URI)
	}
	if i.SIP.SipMsg.Authorization == nil {
		i.SIP.SipMsg.Authorization = new(sipparser.Authorization)
	}
	if i.SIP.SipMsg.Via == nil {
		i.SIP.SipMsg.Via = make([]*sipparser.Via, 1)
		i.SIP.SipMsg.Via[0] = new(sipparser.Via)
	}
	if i.SIP.SipMsg.Cseq == nil {
		i.SIP.SipMsg.Cseq = new(sipparser.Cseq)
	}
	if i.SIP.SipMsg.Reason == nil {
		i.SIP.SipMsg.Reason = new(sipparser.Reason)
	}

	if i.SIP.SipMsg.Error != nil {
		return i.SIP.SipMsg.Error
	} else if len(i.SIP.SipMsg.Cseq.Method) < 3 {
		return errors.New("Could not find a valid CSeq in packet")
	} else if len(i.SIP.SipMsg.CallId) < 3 {
		return errors.New("Could not find a valid Call-ID in packet")
	}

	return nil
}
