package main

import (
	"bytes"
	"fmt"
)

// SendStatsD creates a map with QOS or SIP stats which will
// be converted into statsd compatible strings seperated by '\n'
func (conn Connections) SendStatsD(i *IPFIX, s string) {
	buf := new(bytes.Buffer)
	var mapQOS map[string]float32
	switch s {
	case "QOS":
		mapQOS = map[string]float32{
			"qos.inc.rtp.mos":          float32(i.Data.QOS.IncMos) / 100,
			"qos.out.rtp.mos":          float32(i.Data.QOS.OutMos) / 100,
			"qos.inc.rtp.rval":         float32(i.Data.QOS.IncrVal) / 100,
			"qos.out.rtp.rval":         float32(i.Data.QOS.OutrVal) / 100,
			"qos.inc.rtp.packets":      float32(i.Data.QOS.IncRtpPackets),
			"qos.out.rtp.packets":      float32(i.Data.QOS.OutRtpPackets),
			"qos.inc.rtcp.packets":     float32(i.Data.QOS.IncRtcpPackets),
			"qos.out.rtcp.packets":     float32(i.Data.QOS.OutRtcpPackets),
			"qos.inc.rtp.lostPackets":  float32(i.Data.QOS.IncRtpLostPackets),
			"qos.out.rtp.lostPackets":  float32(i.Data.QOS.OutRtpLostPackets),
			"qos.inc.rtcp.lostPackets": float32(i.Data.QOS.IncRtcpLostPackets),
			"qos.out.rtcp.lostPackets": float32(i.Data.QOS.OutRtcpLostPackets),
			"qos.inc.rtp.avgJitter":    float32(i.Data.QOS.IncRtpAvgJitter),
			"qos.out.rtp.avgJitter":    float32(i.Data.QOS.OutRtpAvgJitter),
			"qos.inc.rtp.maxJitter":    float32(i.Data.QOS.IncRtpMaxJitter),
			"qos.out.rtp.maxJitter":    float32(i.Data.QOS.OutRtpMaxJitter),
			"qos.inc.rtcp.avgJitter":   float32(i.Data.QOS.IncRtcpAvgJitter),
			"qos.out.rtcp.avgJitter":   float32(i.Data.QOS.OutRtcpAvgJitter),
			"qos.inc.rtcp.maxJitter":   float32(i.Data.QOS.IncRtcpMaxJitter),
			"qos.out.rtcp.maxJitter":   float32(i.Data.QOS.OutRtcpMaxJitter),
			"qos.inc.rtcp.avgLat":      float32(i.Data.QOS.IncRtcpAvgLat),
			"qos.out.rtcp.avgLat":      float32(i.Data.QOS.OutRtcpAvgLat),
			"qos.inc.rtcp.maxLat":      float32(i.Data.QOS.IncRtcpMaxLat),
			"qos.out.rtcp.maxLat":      float32(i.Data.QOS.OutRtcpMaxLat),
		}
		for metric, value := range mapQOS {
			buf.Write([]byte(fmt.Sprintf("%s.%s:%.2f|h\n", *name, metric, value)))
		}
	}
	conn.StatsD.Write(buf.Bytes())
}
