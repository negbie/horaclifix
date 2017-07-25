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
			"timer.mean_90." + *name + ".IncMos":             float32(i.Data.QOS.IncMos) / 100,
			"timer.mean_90." + *name + ".OutMos":             float32(i.Data.QOS.OutMos) / 100,
			"timer.mean_90." + *name + ".IncrVal":            float32(i.Data.QOS.IncrVal) / 100,
			"timer.mean_90." + *name + ".OutrVal":            float32(i.Data.QOS.OutrVal) / 100,
			"timer.mean_90." + *name + ".IncRtpPackets":      float32(i.Data.QOS.IncRtpPackets),
			"timer.mean_90." + *name + ".OutRtpPackets":      float32(i.Data.QOS.OutRtpPackets),
			"timer.mean_90." + *name + ".IncRtcpPackets":     float32(i.Data.QOS.IncRtcpPackets),
			"timer.mean_90." + *name + ".OutRtcpPackets":     float32(i.Data.QOS.OutRtcpPackets),
			"timer.mean_90." + *name + ".IncRtpLostPackets":  float32(i.Data.QOS.IncRtpLostPackets),
			"timer.mean_90." + *name + ".OutRtpLostPackets":  float32(i.Data.QOS.OutRtpLostPackets),
			"timer.mean_90." + *name + ".IncRtcpLostPackets": float32(i.Data.QOS.IncRtcpLostPackets),
			"timer.mean_90." + *name + ".OutRtcpLostPackets": float32(i.Data.QOS.OutRtcpLostPackets),
			"timer.mean_90." + *name + ".IncRtpAvgJitter":    float32(i.Data.QOS.IncRtpAvgJitter),
			"timer.mean_90." + *name + ".OutRtpAvgJitter":    float32(i.Data.QOS.OutRtpAvgJitter),
			"timer.mean_90." + *name + ".IncRtpMaxJitter":    float32(i.Data.QOS.IncRtpMaxJitter),
			"timer.mean_90." + *name + ".OutRtpMaxJitter":    float32(i.Data.QOS.OutRtpMaxJitter),
			"timer.mean_90." + *name + ".IncRtcpAvgJitter":   float32(i.Data.QOS.IncRtcpAvgJitter),
			"timer.mean_90." + *name + ".OutRtcpAvgJitter":   float32(i.Data.QOS.OutRtcpAvgJitter),
			"timer.mean_90." + *name + ".IncRtcpMaxJitter":   float32(i.Data.QOS.IncRtcpMaxJitter),
			"timer.mean_90." + *name + ".OutRtcpMaxJitter":   float32(i.Data.QOS.OutRtcpMaxJitter),
			"timer.mean_90." + *name + ".IncRtcpAvgLat":      float32(i.Data.QOS.IncRtcpAvgLat),
			"timer.mean_90." + *name + ".OutRtcpAvgLat":      float32(i.Data.QOS.OutRtcpAvgLat),
			"timer.mean_90." + *name + ".IncRtcpMaxLat":      float32(i.Data.QOS.IncRtcpMaxLat),
			"timer.mean_90." + *name + ".OutRtcpMaxLat":      float32(i.Data.QOS.OutRtcpMaxLat),
		}
	}

	for metric, value := range mapQOS {
		buf.Write([]byte(fmt.Sprintf("%s %d %.2f\n", metric, i.Data.QOS.EndTimeSec, value)))
	}
	conn.Banshee.Write(buf.Bytes())
}
