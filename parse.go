package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strconv"

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
	i.SIP.VL = uint8(msg[15] >> 4)
	i.SIP.TOS = uint8(msg[16])
	i.SIP.TLen = binary.BigEndian.Uint16(msg[17:19])
	i.SIP.TID = binary.BigEndian.Uint16(msg[19:21])
	i.SIP.TFlags = binary.BigEndian.Uint16(msg[21:23])
	if i.SIP.VL != 6 {
		i.SIP.TTL = uint8(msg[23])
		i.SIP.TProto = uint8(msg[24])
		i.SIP.TPos = binary.BigEndian.Uint16(msg[25:27])
		i.SIP.SrcIP = net.IPv4(msg[27], msg[28], msg[29], msg[30])
		i.SIP.DstIP = net.IPv4(msg[31], msg[32], msg[33], msg[34])
		i.SIP.SrcPort = binary.BigEndian.Uint16(msg[35:37])
		i.SIP.DstPort = binary.BigEndian.Uint16(msg[37:39])
		i.SIP.UDPlen = binary.BigEndian.Uint16(msg[39:41])
		i.SIP.MsgLen = binary.BigEndian.Uint16(msg[41:43])
		i.SIP.RawMsg = msg[43:]
	} else {
		i.SIP.TProto = 17

		i.SIP.SrcIP = net.IP{msg[23], msg[24], msg[25], msg[26], msg[27], msg[28], msg[29],
			msg[30], msg[31], msg[32], msg[33], msg[34], msg[35], msg[36], msg[37], msg[38]}

		i.SIP.DstIP = net.IP{msg[39], msg[40], msg[41], msg[42], msg[43], msg[44], msg[45],
			msg[46], msg[47], msg[48], msg[49], msg[50], msg[51], msg[52], msg[53], msg[54]}

		i.SIP.SrcPort = binary.BigEndian.Uint16(msg[55:57])
		i.SIP.DstPort = binary.BigEndian.Uint16(msg[57:59])
		i.SIP.UDPlen = binary.BigEndian.Uint16(msg[59:61])
		i.SIP.MsgLen = binary.BigEndian.Uint16(msg[61:63])
		i.SIP.RawMsg = msg[63:]
	}

	i.SIP.SrcIPString = i.SIP.SrcIP.String()
	i.SIP.DstIPString = i.SIP.DstIP.String()

	if *gaddr != "" {
		err := i.parseSIP()
		if err != nil {
			log.Printf("Could not parse SIP msg: %s", strconv.Quote(string(i.SIP.RawMsg)))
		}
		if *filter != "" && i.SIP.SipMsg.CseqMethod == *filter {
			i.SIP.RawMsg = nil
		}
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
	i.SIP.VL = uint8(msg[pos+3] >> 4)
	i.SIP.TOS = uint8(msg[pos+4])
	i.SIP.TLen = binary.BigEndian.Uint16(msg[pos+5 : pos+7])
	i.SIP.TID = binary.BigEndian.Uint16(msg[pos+7 : pos+9])
	i.SIP.TFlags = binary.BigEndian.Uint16(msg[pos+9 : pos+11])
	if i.SIP.VL != 6 {
		i.SIP.TTL = uint8(msg[pos+11])
		i.SIP.TProto = uint8(msg[pos+12])
		i.SIP.TPos = binary.BigEndian.Uint16(msg[pos+13 : pos+15])
		i.SIP.SrcIP = net.IPv4(msg[pos+15], msg[pos+16], msg[pos+17], msg[pos+18])
		i.SIP.DstIP = net.IPv4(msg[pos+19], msg[pos+20], msg[pos+21], msg[pos+22])
		i.SIP.SrcPort = binary.BigEndian.Uint16(msg[pos+23 : pos+25])
		i.SIP.DstPort = binary.BigEndian.Uint16(msg[pos+25 : pos+27])
		i.SIP.UDPlen = binary.BigEndian.Uint16(msg[pos+27 : pos+29])
		i.SIP.MsgLen = binary.BigEndian.Uint16(msg[pos+29 : pos+31])
		i.SIP.RawMsg = []byte(msg[pos+31:])
	} else {
		i.SIP.TProto = 17

		i.SIP.SrcIP = net.IP{msg[pos+11], msg[pos+12], msg[pos+13], msg[pos+14], msg[pos+15], msg[pos+16], msg[pos+17],
			msg[pos+18], msg[pos+19], msg[pos+20], msg[pos+21], msg[pos+22], msg[pos+23], msg[pos+24], msg[pos+25], msg[pos+26]}

		i.SIP.DstIP = net.IP{msg[pos+27], msg[pos+28], msg[pos+29], msg[pos+30], msg[pos+31], msg[pos+32], msg[pos+33],
			msg[pos+34], msg[pos+35], msg[pos+36], msg[pos+37], msg[pos+38], msg[pos+39], msg[pos+40], msg[pos+41], msg[pos+42]}

		i.SIP.SrcPort = binary.BigEndian.Uint16(msg[pos+43 : pos+45])
		i.SIP.DstPort = binary.BigEndian.Uint16(msg[pos+45 : pos+47])
		i.SIP.UDPlen = binary.BigEndian.Uint16(msg[pos+47 : pos+49])
		i.SIP.MsgLen = binary.BigEndian.Uint16(msg[pos+49 : pos+51])
		i.SIP.RawMsg = []byte(msg[pos+51:])
	}

	i.SIP.SrcIPString = i.SIP.SrcIP.String()
	i.SIP.DstIPString = i.SIP.DstIP.String()

	if *gaddr != "" {
		err := i.parseSIP()
		if err != nil {
			log.Printf("Could not parse SIP msg: %s", strconv.Quote(string(i.SIP.RawMsg)))
		}
		if *filter != "" && i.SIP.SipMsg.CseqMethod == *filter {
			i.SIP.RawMsg = nil
		}
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
	i.SIP.DstIP = net.IPv4(msg[12], msg[13], msg[14], msg[15])
	i.SIP.SrcIP = net.IPv4(msg[16], msg[17], msg[18], msg[19])
	i.SIP.DstPort = binary.BigEndian.Uint16(msg[20:22])
	i.SIP.SrcPort = binary.BigEndian.Uint16(msg[22:24])
	i.SIP.Context = binary.BigEndian.Uint32(msg[24:28])
	i.SIP.CallIDEnd = uint8(msg[28])
	i.SIP.MsgLen = binary.BigEndian.Uint16(msg[29:31])
	i.SIP.RawMsg = msg[31:]

	i.SIP.SrcIPString = i.SIP.SrcIP.String()
	i.SIP.DstIPString = i.SIP.DstIP.String()

	if *gaddr != "" {
		err := i.parseSIP()
		if err != nil {
			log.Printf("Could not parse SIP msg: %s", strconv.Quote(string(i.SIP.RawMsg)))
		}
		if *filter != "" && i.SIP.SipMsg.CseqMethod == *filter {
			i.SIP.RawMsg = nil
		}
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
	i.SIP.DstIP = net.IPv4(msg[12], msg[13], msg[14], msg[15])
	i.SIP.SrcIP = net.IPv4(msg[16], msg[17], msg[18], msg[19])
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

	i.SIP.SrcIPString = i.SIP.SrcIP.String()
	i.SIP.DstIPString = i.SIP.DstIP.String()

	if *gaddr != "" {
		err := i.parseSIP()
		if err != nil {
			log.Printf("Could not parse SIP msg: %s", strconv.Quote(string(i.SIP.RawMsg)))
		}
		if *filter != "" && i.SIP.SipMsg.CseqMethod == *filter {
			i.SIP.RawMsg = nil
		}
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
	r.binRead(&i.QOS.EndTimeMic)

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

	i.normQOS()

	return &i
}

func (i *IPFIX) parseSIP() error {
	i.SIP.SipMsg = sipparser.ParseMsg(string(i.SIP.RawMsg))

	if i.SIP.SipMsg.StartLine == nil {
		i.SIP.SipMsg.StartLine = new(sipparser.StartLine)
	}
	if i.SIP.SipMsg.StartLine.URI == nil {
		i.SIP.SipMsg.StartLine.URI = new(sipparser.URI)
	}

	if i.SIP.SipMsg.Error != nil {
		return i.SIP.SipMsg.Error
	} else if len(i.SIP.SipMsg.CseqMethod) < 3 {
		return errors.New("Could not find a valid CSeq in packet")
	} else if len(i.SIP.SipMsg.CallID) < 1 {
		return errors.New("Could not find a valid Call-ID in packet")
	}

	return nil
}

func (i *IPFIX) normQOS() {
	i.QOS.IncRtpPackets = normMaxQOS(i.QOS.IncRtpPackets)
	i.QOS.IncRtpLostPackets = normMaxQOS(i.QOS.IncRtpLostPackets)
	i.QOS.IncRtpAvgJitter = normMaxQOS(i.QOS.IncRtpAvgJitter)
	i.QOS.IncRtpMaxJitter = normMaxQOS(i.QOS.IncRtpMaxJitter)
	i.QOS.IncRtcpPackets = normMaxQOS(i.QOS.IncRtcpPackets)
	i.QOS.IncRtcpLostPackets = normMaxQOS(i.QOS.IncRtcpLostPackets)
	i.QOS.IncRtcpAvgJitter = normMaxQOS(i.QOS.IncRtcpAvgJitter)
	i.QOS.IncRtcpMaxJitter = normMaxQOS(i.QOS.IncRtcpMaxJitter)
	i.QOS.IncRtcpAvgLat = normMaxQOS(i.QOS.IncRtcpAvgLat)
	i.QOS.IncRtcpMaxLat = normMaxQOS(i.QOS.IncRtcpMaxLat)

	i.QOS.OutRtpPackets = normMaxQOS(i.QOS.OutRtpPackets)
	i.QOS.OutRtpLostPackets = normMaxQOS(i.QOS.OutRtpLostPackets)
	i.QOS.OutRtpAvgJitter = normMaxQOS(i.QOS.OutRtpAvgJitter)
	i.QOS.OutRtpMaxJitter = normMaxQOS(i.QOS.OutRtpMaxJitter)
	i.QOS.OutRtcpPackets = normMaxQOS(i.QOS.OutRtcpPackets)
	i.QOS.OutRtcpLostPackets = normMaxQOS(i.QOS.OutRtcpLostPackets)
	i.QOS.OutRtcpAvgJitter = normMaxQOS(i.QOS.OutRtcpAvgJitter)
	i.QOS.OutRtcpMaxJitter = normMaxQOS(i.QOS.OutRtcpMaxJitter)
	i.QOS.OutRtcpAvgLat = normMaxQOS(i.QOS.OutRtcpAvgLat)
	i.QOS.OutRtcpMaxLat = normMaxQOS(i.QOS.OutRtcpMaxLat)
}
