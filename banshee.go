package main

import (
	"bytes"
	"fmt"
)

// SendBanshee creates a map with QOS or SIP stats which will
// be converted into statsd like strings seperated by '\n'
func (conn Connections) SendBanshee(i *IPFIX, s string) {
	var mapQOS map[string]interface{}
	buf := new(bytes.Buffer)
	switch s {
	case "QOS":
		mapQOS = map[string]interface{}{
			"timer.mean_90." + *name + ".inc.mos":              float32(i.Data.QOS.IncMos) / 100,
			"timer.mean_90." + *name + ".out.mos":              float32(i.Data.QOS.OutMos) / 100,
			"timer.mean_90." + *name + ".inc.rval":             float32(i.Data.QOS.IncrVal) / 100,
			"timer.mean_90." + *name + ".out.rval":             float32(i.Data.QOS.OutrVal) / 100,
			"timer.mean_90." + *name + ".inc.rtp.packets":      float32(i.Data.QOS.IncRtpPackets),
			"timer.mean_90." + *name + ".out.rtp.packets":      float32(i.Data.QOS.OutRtpPackets),
			"timer.mean_90." + *name + ".inc.rtcp.packets":     float32(i.Data.QOS.IncRtcpPackets),
			"timer.mean_90." + *name + ".out.rtcp.packets":     float32(i.Data.QOS.OutRtcpPackets),
			"timer.mean_90." + *name + ".inc.rtp.lostPackets":  float32(i.Data.QOS.IncRtpLostPackets),
			"timer.mean_90." + *name + ".out.rtp.lostPackets":  float32(i.Data.QOS.OutRtpLostPackets),
			"timer.mean_90." + *name + ".inc.rtcp.lostPackets": float32(i.Data.QOS.IncRtcpLostPackets),
			"timer.mean_90." + *name + ".out.rtcp.lostPackets": float32(i.Data.QOS.OutRtcpLostPackets),
			"timer.mean_90." + *name + ".inc.rtp.avgJitter":    float32(i.Data.QOS.IncRtpAvgJitter),
			"timer.mean_90." + *name + ".out.rtp.avgJitter":    float32(i.Data.QOS.OutRtpAvgJitter),
			"timer.mean_90." + *name + ".inc.rtp.maxJitter":    float32(i.Data.QOS.IncRtpMaxJitter),
			"timer.mean_90." + *name + ".out.rtp.maxJitter":    float32(i.Data.QOS.OutRtpMaxJitter),
			"timer.mean_90." + *name + ".inc.rtcp.avgJitter":   float32(i.Data.QOS.IncRtcpAvgJitter),
			"timer.mean_90." + *name + ".out.rtcp.avgJitter":   float32(i.Data.QOS.OutRtcpAvgJitter),
			"timer.mean_90." + *name + ".inc.rtcp.maxJitter":   float32(i.Data.QOS.IncRtcpMaxJitter),
			"timer.mean_90." + *name + ".out.rtcp.maxJitter":   float32(i.Data.QOS.OutRtcpMaxJitter),
			"timer.mean_90." + *name + ".inc.rtcp.avgLat":      float32(i.Data.QOS.IncRtcpAvgLat),
			"timer.mean_90." + *name + ".out.rtcp.avgLat":      float32(i.Data.QOS.OutRtcpAvgLat),
			"timer.mean_90." + *name + ".inc.rtcp.maxLat":      float32(i.Data.QOS.IncRtcpMaxLat),
			"timer.mean_90." + *name + ".out.rtcp.maxLat":      float32(i.Data.QOS.OutRtcpMaxLat),
		}
	}

	for metric, value := range mapQOS {
		buf.Write([]byte(fmt.Sprintf("%s %d %.2f\n", metric, i.Data.QOS.EndTimeSec, value)))
	}
	conn.Banshee.Write(buf.Bytes())
}
