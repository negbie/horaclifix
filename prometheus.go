package main

func (conn *Connections) SendMetric(i *IPFIX, s string) {
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_mos"].Set(float64(i.Data.QOS.IncMos / 100))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtp_mos"].Set(float64(i.Data.QOS.OutMos / 100))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_rval"].Set(float64(i.Data.QOS.IncrVal / 100))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtp_rval"].Set(float64(i.Data.QOS.OutrVal / 100))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_packets"].Set(float64(i.Data.QOS.IncRtpPackets))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtp_packets"].Set(float64(i.Data.QOS.OutRtpPackets))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_packets"].Set(float64(i.Data.QOS.IncRtcpPackets))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_packets"].Set(float64(i.Data.QOS.OutRtcpPackets))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_lostPackets"].Set(float64(i.Data.QOS.IncRtpLostPackets))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtp_lostPackets"].Set(float64(i.Data.QOS.OutRtpLostPackets))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_lostPackets"].Set(float64(i.Data.QOS.IncRtcpLostPackets))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_lostPackets"].Set(float64(i.Data.QOS.OutRtcpLostPackets))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_avgJitter"].Set(float64(i.Data.QOS.IncRtpAvgJitter))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtp_avgJitter"].Set(float64(i.Data.QOS.OutRtpAvgJitter))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_maxJitter"].Set(float64(i.Data.QOS.IncRtpMaxJitter))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtp_maxJitter"].Set(float64(i.Data.QOS.OutRtpMaxJitter))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_avgJitter"].Set(float64(i.Data.QOS.IncRtcpAvgJitter))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_avgJitter"].Set(float64(i.Data.QOS.OutRtcpAvgJitter))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_maxJitter"].Set(float64(i.Data.QOS.IncRtcpMaxJitter))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_maxJitter"].Set(float64(i.Data.QOS.OutRtcpMaxJitter))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_avgLat"].Set(float64(i.Data.QOS.IncRtcpAvgLat))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_avgLat"].Set(float64(i.Data.QOS.OutRtcpAvgLat))
	conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_maxLat"].Set(float64(i.Data.QOS.IncRtcpMaxLat))
	conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_maxLat"].Set(float64(i.Data.QOS.OutRtcpMaxLat))
	if i.Data.QOS.EndTimeSec != 0 && i.Data.QOS.BeginTimeSec != 0 {
		conn.Prometheus.GaugeMetrics[*name+"_duration"].Set(float64(i.Data.QOS.EndTimeSec - i.Data.QOS.BeginTimeSec))
	}

}
