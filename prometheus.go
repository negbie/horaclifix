package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	GaugeVecMetrics   map[string]*prometheus.GaugeVec
	CounterVecMetrics map[string]*prometheus.CounterVec
	//CounterMetrics    map[string]prometheus.Counter
	//GaugeMetrics      map[string]prometheus.Gauge
}

var prom Prometheus

func regProm() {
	if *paddr != "" {
		prom.GaugeVecMetrics = map[string]*prometheus.GaugeVec{}
		prom.CounterVecMetrics = map[string]*prometheus.CounterVec{}

		prom.CounterVecMetrics["horaclifix_packets_total"] = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "horaclifix_packets_total", Help: "Total received packets"}, []string{"sbc_name"})

		prom.GaugeVecMetrics["horaclifix_inc_rtp_mos"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtp_mos", Help: "Incoming RTP MOS"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtp_rval"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtp_rval", Help: "Incoming RTP rVal"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtp_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtp_packets", Help: "Incoming RTP packets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtp_lost_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtp_lost_packets", Help: "Incoming RTP lostPackets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtp_avg_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtp_avg_jitter", Help: "Incoming RTP avgJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtp_max_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtp_max_jitter", Help: "Incoming RTP maxJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtcp_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtcp_packets", Help: "Incoming RTCP packets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtcp_lost_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtcp_lost_packets", Help: "Incoming RTCP lostPackets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtcp_avg_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtcp_avg_jitter", Help: "Incoming RTCP avgJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtcp_max_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtcp_max_jitter", Help: "Incoming RTCP maxJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtcp_avg_lat"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtcp_avg_lat", Help: "Incoming RTCP avgLat"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_inc_rtcp_max_lat"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_inc_rtcp_max_lat", Help: "Incoming RTCP maxLat"}, []string{"sbc_name", "inc_realm", "out_realm"})

		prom.GaugeVecMetrics["horaclifix_out_rtp_mos"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtp_mos", Help: "Outgoing RTP MOS"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtp_rval"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtp_rval", Help: "Outgoing RTP rVal"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtp_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtp_packets", Help: "Outgoing RTP packets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtp_lost_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtp_lost_packets", Help: "Outgoing RTP lostPackets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtp_avg_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtp_avg_jitter", Help: "Outgoing RTP avgJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtp_max_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtp_max_jitter", Help: "Outgoing RTP maxJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtcp_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtcp_packets", Help: "Outgoing RTCP packets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtcp_lost_packets"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtcp_lost_packets", Help: "Outgoing RTCP lostPackets"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtcp_avg_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtcp_avg_jitter", Help: "Outgoing RTCP avgJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtcp_max_jitter"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtcp_max_jitter", Help: "Outgoing RTCP maxJitter"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtcp_avg_lat"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtcp_avg_lat", Help: "Outgoing RTCP avgLat"}, []string{"sbc_name", "inc_realm", "out_realm"})
		prom.GaugeVecMetrics["horaclifix_out_rtcp_max_lat"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_out_rtcp_max_lat", Help: "Outgoing RTCP maxLat"}, []string{"sbc_name", "inc_realm", "out_realm"})

		prom.GaugeVecMetrics["horaclifix_duration"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "horaclifix_duration", Help: "Call duration"}, []string{"sbc_name", "inc_realm", "out_realm"})

		for k := range prom.CounterVecMetrics {
			log.Printf("register prometheus counterMetric %s", k)
			prometheus.MustRegister(prom.CounterVecMetrics[k])
		}

		for k := range prom.GaugeVecMetrics {
			log.Printf("register prometheus gaugeMetric %s", k)
			prometheus.MustRegister(prom.GaugeVecMetrics[k])
		}

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(*paddr, nil)
			checkCritErr(err)
		}()
	}
}

func (conn *Connections) SendMetric(i *IPFIX, s string) {
	incRealm := string(i.QOS.IncRealm)
	outRealm := string(i.QOS.OutRealm)

	prom.GaugeVecMetrics["horaclifix_inc_rtp_mos"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncMos / 100))
	prom.GaugeVecMetrics["horaclifix_inc_rtp_rval"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncrVal / 100))
	prom.GaugeVecMetrics["horaclifix_inc_rtp_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtpPackets))
	prom.GaugeVecMetrics["horaclifix_inc_rtp_lost_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtpLostPackets))
	prom.GaugeVecMetrics["horaclifix_inc_rtp_avg_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtpAvgJitter))
	prom.GaugeVecMetrics["horaclifix_inc_rtp_max_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtpMaxJitter))
	prom.GaugeVecMetrics["horaclifix_inc_rtcp_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtcpPackets))
	prom.GaugeVecMetrics["horaclifix_inc_rtcp_lost_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtcpLostPackets))
	prom.GaugeVecMetrics["horaclifix_inc_rtcp_avg_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtcpAvgJitter))
	prom.GaugeVecMetrics["horaclifix_inc_rtcp_max_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtcpMaxJitter))
	prom.GaugeVecMetrics["horaclifix_inc_rtcp_avg_lat"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtcpAvgLat))
	prom.GaugeVecMetrics["horaclifix_inc_rtcp_max_lat"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.IncRtcpMaxLat))

	prom.GaugeVecMetrics["horaclifix_out_rtp_mos"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutMos / 100))
	prom.GaugeVecMetrics["horaclifix_out_rtp_rval"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutrVal / 100))
	prom.GaugeVecMetrics["horaclifix_out_rtp_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtpPackets))
	prom.GaugeVecMetrics["horaclifix_out_rtp_lost_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtpLostPackets))
	prom.GaugeVecMetrics["horaclifix_out_rtp_avg_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtpAvgJitter))
	prom.GaugeVecMetrics["horaclifix_out_rtp_max_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtpMaxJitter))
	prom.GaugeVecMetrics["horaclifix_out_rtcp_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtcpPackets))
	prom.GaugeVecMetrics["horaclifix_out_rtcp_lost_packets"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtcpLostPackets))
	prom.GaugeVecMetrics["horaclifix_out_rtcp_avg_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtcpAvgJitter))
	prom.GaugeVecMetrics["horaclifix_out_rtcp_max_jitter"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtcpMaxJitter))
	prom.GaugeVecMetrics["horaclifix_out_rtcp_avg_lat"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtcpAvgLat))
	prom.GaugeVecMetrics["horaclifix_out_rtcp_max_lat"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.OutRtcpMaxLat))

	if i.QOS.EndTimeSec != 0 && i.QOS.BeginTimeSec != 0 {
		prom.GaugeVecMetrics["horaclifix_duration"].WithLabelValues(*name, incRealm, outRealm).Set(float64(i.QOS.EndTimeSec - i.QOS.BeginTimeSec))
	}
}
