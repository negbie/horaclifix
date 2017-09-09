package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

// ParseRecSipUDP fills the SipSet struct with
// the dataSet 258
func ParseRecSipUDP(header []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("ParseRecSipUDP.timetaken", time.Since(t))
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
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[55:57])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(header[57:59])
	i.Data.SIP.UDPlen = binary.BigEndian.Uint16(header[59:61])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[61:63])
	i.Data.SIP.SipMsg = header[63:]
	return &i
}

// ParseSendSipUDP fills the SipSet struct with
// the dataSet 259
func ParseSendSipUDP(header []byte) *IPFIX {
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
	i.Data.SIP.CallIDLen = uint8(header[32])

	pos := int(i.Data.SIP.CallIDLen)
	if pos == 0 {
		pos = 33
	} else {
		pos = pos + 33
		i.Data.SIP.CallID = []byte(header[33:pos])
	}
	i.Data.SIP.CallIDEnd = uint8(header[pos])
	i.Data.SIP.IPlen = binary.BigEndian.Uint16(header[pos+1 : pos+3])
	i.Data.SIP.VL = uint8(header[pos+3])
	i.Data.SIP.TOS = uint8(header[pos+4])
	i.Data.SIP.TLen = binary.BigEndian.Uint16(header[pos+5 : pos+7])
	i.Data.SIP.TID = binary.BigEndian.Uint16(header[pos+7 : pos+9])
	i.Data.SIP.TFlags = binary.BigEndian.Uint16(header[pos+9 : pos+11])
	i.Data.SIP.TTL = uint8(header[pos+11])
	i.Data.SIP.TProto = uint8(header[pos+12])
	i.Data.SIP.TPos = binary.BigEndian.Uint16(header[pos+13 : pos+15])
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[pos+15 : pos+19])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(header[pos+19 : pos+23])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[pos+23 : pos+25])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(header[pos+25 : pos+27])
	i.Data.SIP.UDPlen = binary.BigEndian.Uint16(header[pos+27 : pos+29])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[pos+29 : pos+31])
	i.Data.SIP.SipMsg = []byte(header[pos+31:])
	return &i
}

// ParseRecSipTCP fills the SipSet struct with
// the dataSet 260
func ParseRecSipTCP(header []byte) *IPFIX {
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
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[32:36])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(header[36:40])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[40:42])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(header[42:44])
	i.Data.SIP.Context = binary.BigEndian.Uint32(header[44:48])
	i.Data.SIP.CallIDEnd = uint8(header[48])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[49:51])
	i.Data.SIP.SipMsg = header[51:]
	return &i
}

// ParseSendSipTCP fills the SipSet struct with
// the dataSet 261
func ParseSendSipTCP(header []byte) *IPFIX {
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
	i.Data.SIP.SrcIP = binary.BigEndian.Uint32(header[32:36])
	i.Data.SIP.DstIP = binary.BigEndian.Uint32(header[36:40])
	i.Data.SIP.SrcPort = binary.BigEndian.Uint16(header[40:42])
	i.Data.SIP.DstPort = binary.BigEndian.Uint16(header[42:44])
	i.Data.SIP.Context = binary.BigEndian.Uint32(header[44:48])
	i.Data.SIP.CallIDLen = uint8(header[48])

	pos := int(i.Data.SIP.CallIDLen)
	if pos == 0 {
		pos = 49
	} else {
		pos = pos + 49
		i.Data.SIP.CallID = []byte(header[49:pos])
	}
	i.Data.SIP.CallIDEnd = uint8(header[pos])
	i.Data.SIP.MsgLen = binary.BigEndian.Uint16(header[pos+1 : pos+3])
	i.Data.SIP.SipMsg = []byte(header[pos+3:])
	return &i
}

// ParseQosStats fills the QosSet struct with
// the dataSet 268
func ParseQosStats(header []byte) *IPFIX {
	var i IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &i.Header)
	binary.Read(r, binary.BigEndian, &i.SetHeader)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtpBytes)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtpPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtpLostPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpBytes)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRtcpMaxLat)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncrVal)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncMos)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtpBytes)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtpPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtpLostPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpBytes)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRtcpMaxLat)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutrVal)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutMos)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.Type)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIncSrcIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIncDstIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIncSrcPort)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIncDstPort)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIncSrcIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIncDstIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIncSrcPort)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIncDstPort)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerOutSrcIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerOutDstIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerOutSrcPort)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerOutDstPort)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeOutSrcIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeOutDstIP)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeOutSrcPort)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeOutDstPort)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIntSlot)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIntPort)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CallerIntVlan)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIntSlot)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIntPort)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.CalleeIntVlan)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.BeginTimeSec)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.BeginTimeMic)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.EndTimeSec)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.EndinTimeMic)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.Seperator)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRealmLen)
	i.Data.QOS.IncRealm = make([]byte, i.Data.QOS.IncRealmLen)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRealm)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncRealmEnd)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRealmLen)
	i.Data.QOS.OutRealm = make([]byte, i.Data.QOS.OutRealmLen)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRealm)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutRealmEnd)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncCallIDLen)
	i.Data.QOS.IncCallID = make([]byte, i.Data.QOS.IncCallIDLen)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncCallID)
	binary.Read(r, binary.BigEndian, &i.Data.QOS.IncCallIDEnd)

	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutCallIDLen)
	i.Data.QOS.OutCallID = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &i.Data.QOS.OutCallID)

	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return &i
}

// stringIPv4 converts a ipv4 unit32 into a string
func stringIPv4(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}

/*
// Template for a sync.Pool buffer
var buffers = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 65536)
	},
}
packet := buffers.Get().([]byte)
buffers.Put(packet)
*/
