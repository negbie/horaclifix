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
