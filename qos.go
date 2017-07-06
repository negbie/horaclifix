package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func NewQosStats(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
	err = binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpBytes)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpLostPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpAvgJitter)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtpMaxJitter)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpBytes)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpLostPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpAvgJitter)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpMaxJitter)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpAvgLat)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRtcpMaxLat)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncrVal)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncMos)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpBytes)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpLostPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpAvgJitter)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtpMaxJitter)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpBytes)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpLostPackets)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpAvgJitter)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpMaxJitter)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpAvgLat)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRtcpMaxLat)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutrVal)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutMos)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.Type)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncSrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncDstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncSrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIncDstPort)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncSrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncDstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncSrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIncDstPort)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutSrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutDstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutSrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerOutDstPort)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutSrcIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutDstIP)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutSrcPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeOutDstPort)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIntSlot)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIntPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CallerIntVlan)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIntSlot)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIntPort)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.CalleeIntVlan)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.BeginTimeSec)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.BeginTimeMic)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.EndTimeSec)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.EndinTimeMic)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.Seperator)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRealmLen)
	ipfix.Data.QOS.IncRealm = make([]byte, ipfix.Data.QOS.IncRealmLen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRealm)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncRealmEnd)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRealmLen)
	ipfix.Data.QOS.OutRealm = make([]byte, ipfix.Data.QOS.OutRealmLen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRealm)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutRealmEnd)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncCallIDLen)
	ipfix.Data.QOS.IncCallID = make([]byte, ipfix.Data.QOS.IncCallIDLen)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncCallID)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.IncCallIDEnd)

	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutCallIDLen)
	ipfix.Data.QOS.OutCallID = make([]byte, r.Len())
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.QOS.OutCallID)

	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	return &ipfix
}
