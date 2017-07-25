package main

import (
	"bytes"
	"fmt"
)

// SendStatsD creates a map with QOS or SIP stats which will
// be converted into statsd compatible strings seperated by '\n'
func (conn Connections) SendStatsD(i *IPFIX, s string) {
	buf := new(bytes.Buffer)
	var mapQOS map[string]interface{}
	switch s {
	case "QOS":
		mapQOS = map[string]interface{}{
			"QOS.Inc.RTP.Mos":          float32(i.Data.QOS.IncMos) / 100,
			"QOS.Out.RTP.Mos":          float32(i.Data.QOS.OutMos) / 100,
			"QOS.Inc.RTCP.Mos":         float32(i.Data.QOS.IncMos) / 100,
			"QOS.Out.RTCP.Mos":         float32(i.Data.QOS.OutMos) / 100,
			"QOS.Inc.RTP.Rval":         float32(i.Data.QOS.IncrVal) / 100,
			"QOS.Out.RTP.Rval":         float32(i.Data.QOS.OutrVal) / 100,
			"QOS.Inc.RTCP.Rval":        float32(i.Data.QOS.IncrVal) / 100,
			"QOS.Out.RTCP.Rval":        float32(i.Data.QOS.OutrVal) / 100,
			"QOS.Inc.RTP.Packets":      float32(i.Data.QOS.IncRtpPackets),
			"QOS.Out.RTP.Packets":      float32(i.Data.QOS.OutRtpPackets),
			"QOS.Inc.RTCP.Packets":     float32(i.Data.QOS.IncRtcpPackets),
			"QOS.Out.RTCP.Packets":     float32(i.Data.QOS.OutRtcpPackets),
			"QOS.Inc.RTP.LostPackets":  float32(i.Data.QOS.IncRtpLostPackets),
			"QOS.Out.RTP.LostPackets":  float32(i.Data.QOS.OutRtpLostPackets),
			"QOS.Inc.RTCP.LostPackets": float32(i.Data.QOS.IncRtcpLostPackets),
			"QOS.Out.RTCP.LostPackets": float32(i.Data.QOS.OutRtcpLostPackets),
			"QOS.Inc.RTP.AvgJitter":    float32(i.Data.QOS.IncRtpAvgJitter),
			"QOS.Out.RTP.AvgJitter":    float32(i.Data.QOS.OutRtpAvgJitter),
			"QOS.Inc.RTP.MaxJitter":    float32(i.Data.QOS.IncRtpMaxJitter),
			"QOS.Out.RTP.MaxJitter":    float32(i.Data.QOS.OutRtpMaxJitter),
			"QOS.Inc.RTCP.AvgJitter":   float32(i.Data.QOS.IncRtcpAvgJitter),
			"QOS.Out.RTCP.AvgJitter":   float32(i.Data.QOS.OutRtcpAvgJitter),
			"QOS.Inc.RTCP.MaxJitter":   float32(i.Data.QOS.IncRtcpMaxJitter),
			"QOS.Out.RTCP.MaxJitter":   float32(i.Data.QOS.OutRtcpMaxJitter),
			"QOS.Inc.RTCP.AvgLat":      float32(i.Data.QOS.IncRtcpAvgLat),
			"QOS.Out.RTCP.AvgLat":      float32(i.Data.QOS.OutRtcpAvgLat),
			"QOS.Inc.RTCP.MaxLat":      float32(i.Data.QOS.IncRtcpMaxLat),
			"QOS.Out.RTCP.MaxLat":      float32(i.Data.QOS.OutRtcpMaxLat),
		}
		for metric, value := range mapQOS {
			buf.Write([]byte(fmt.Sprintf("%s.%s:%.2f|h\n", *name, metric, value)))
		}
	}
	conn.StatsD.Write(buf.Bytes())
}
