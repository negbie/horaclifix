package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func NewQosStats(header []byte) *IPFIX {
	/*	t := time.Now()
		defer func() {
			StatTime("NewQosStats.timetaken", time.Since(t))
		}()
	*/
	var ipfix IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
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

	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return &ipfix
}
