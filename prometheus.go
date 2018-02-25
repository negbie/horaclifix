package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var prom Prometheus

func regProm() {
	if *paddr != "" {
		prom.GaugeMetrics = map[string]prometheus.Gauge{}
		prom.CounterMetrics = map[string]prometheus.Counter{}

		prom.CounterMetrics[*name+"_packets"] = prometheus.NewCounter(prometheus.CounterOpts{Name: *name + "_packets", Help: "Received packets"})
		prom.GaugeMetrics[*name+"_inc_rtp_mos"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_mos", Help: "Incoming RTP MOS"})
		prom.GaugeMetrics[*name+"_out_rtp_mos"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_mos", Help: "Outgoing RTP MOS"})
		prom.GaugeMetrics[*name+"_inc_rtp_rval"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_rval", Help: "Incoming RTP rVal"})
		prom.GaugeMetrics[*name+"_out_rtp_rval"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_rval", Help: "Outgoing RTP rVal"})
		prom.GaugeMetrics[*name+"_inc_rtp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_packets", Help: "Incoming RTP packets"})
		prom.GaugeMetrics[*name+"_out_rtp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_packets", Help: "Outgoing RTP packets"})
		prom.GaugeMetrics[*name+"_inc_rtcp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_packets", Help: "Incoming RTCP packets"})
		prom.GaugeMetrics[*name+"_out_rtcp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_packets", Help: "Outgoing RTCP packets"})
		prom.GaugeMetrics[*name+"_inc_rtp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_lostPackets", Help: "Incoming RTP lostPackets"})
		prom.GaugeMetrics[*name+"_out_rtp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_lostPackets", Help: "Outgoing RTP lostPackets"})
		prom.GaugeMetrics[*name+"_inc_rtcp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_lostPackets", Help: "Incoming RTCP lostPackets"})
		prom.GaugeMetrics[*name+"_out_rtcp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_lostPackets", Help: "Outgoing RTCP lostPackets"})
		prom.GaugeMetrics[*name+"_inc_rtp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_avgJitter", Help: "Incoming RTP avgJitter"})
		prom.GaugeMetrics[*name+"_out_rtp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_avgJitter", Help: "Outgoing RTP avgJitter"})
		prom.GaugeMetrics[*name+"_inc_rtp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_maxJitter", Help: "Incoming RTP maxJitter"})
		prom.GaugeMetrics[*name+"_out_rtp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_maxJitter", Help: "Outgoing RTP maxJitter"})
		prom.GaugeMetrics[*name+"_inc_rtcp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_avgJitter", Help: "Incoming RTCP avgJitter"})
		prom.GaugeMetrics[*name+"_out_rtcp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_avgJitter", Help: "Outgoing RTCP avgJitter"})
		prom.GaugeMetrics[*name+"_inc_rtcp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_maxJitter", Help: "Incoming RTCP maxJitter"})
		prom.GaugeMetrics[*name+"_out_rtcp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_maxJitter", Help: "Outgoing RTCP maxJitter"})
		prom.GaugeMetrics[*name+"_inc_rtcp_avgLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_avgLat", Help: "Incoming RTCP avgLat"})
		prom.GaugeMetrics[*name+"_out_rtcp_avgLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_avgLat", Help: "Outgoing RTCP avgLat"})
		prom.GaugeMetrics[*name+"_inc_rtcp_maxLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_maxLat", Help: "Incoming RTCP maxLat"})
		prom.GaugeMetrics[*name+"_out_rtcp_maxLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_maxLat", Help: "Outgoing RTCP maxLat"})
		prom.GaugeMetrics[*name+"_duration"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_duration", Help: "Call duration"})

		for k := range prom.GaugeMetrics {
			log.Printf("register prometheus gaugeMetric %s", k)
			prometheus.MustRegister(prom.GaugeMetrics[k])
		}
		for k := range prom.CounterMetrics {
			log.Printf("register prometheus counterMetric %s", k)
			prometheus.MustRegister(prom.CounterMetrics[k])
		}

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(*paddr, nil)
			checkCritErr(err)
		}()
	}
}

func (conn *Connections) SendMetric(i *IPFIX, s string) {
	prom.GaugeMetrics[*name+"_inc_rtp_mos"].Set(float64(i.Data.QOS.IncMos / 100))
	prom.GaugeMetrics[*name+"_out_rtp_mos"].Set(float64(i.Data.QOS.OutMos / 100))
	prom.GaugeMetrics[*name+"_inc_rtp_rval"].Set(float64(i.Data.QOS.IncrVal / 100))
	prom.GaugeMetrics[*name+"_out_rtp_rval"].Set(float64(i.Data.QOS.OutrVal / 100))
	prom.GaugeMetrics[*name+"_inc_rtp_packets"].Set(float64(i.Data.QOS.IncRtpPackets))
	prom.GaugeMetrics[*name+"_out_rtp_packets"].Set(float64(i.Data.QOS.OutRtpPackets))
	prom.GaugeMetrics[*name+"_inc_rtcp_packets"].Set(float64(i.Data.QOS.IncRtcpPackets))
	prom.GaugeMetrics[*name+"_out_rtcp_packets"].Set(float64(i.Data.QOS.OutRtcpPackets))
	prom.GaugeMetrics[*name+"_inc_rtp_lostPackets"].Set(float64(i.Data.QOS.IncRtpLostPackets))
	prom.GaugeMetrics[*name+"_out_rtp_lostPackets"].Set(float64(i.Data.QOS.OutRtpLostPackets))
	prom.GaugeMetrics[*name+"_inc_rtcp_lostPackets"].Set(float64(i.Data.QOS.IncRtcpLostPackets))
	prom.GaugeMetrics[*name+"_out_rtcp_lostPackets"].Set(float64(i.Data.QOS.OutRtcpLostPackets))
	prom.GaugeMetrics[*name+"_inc_rtp_avgJitter"].Set(float64(i.Data.QOS.IncRtpAvgJitter))
	prom.GaugeMetrics[*name+"_out_rtp_avgJitter"].Set(float64(i.Data.QOS.OutRtpAvgJitter))
	prom.GaugeMetrics[*name+"_inc_rtp_maxJitter"].Set(float64(i.Data.QOS.IncRtpMaxJitter))
	prom.GaugeMetrics[*name+"_out_rtp_maxJitter"].Set(float64(i.Data.QOS.OutRtpMaxJitter))
	prom.GaugeMetrics[*name+"_inc_rtcp_avgJitter"].Set(float64(i.Data.QOS.IncRtcpAvgJitter))
	prom.GaugeMetrics[*name+"_out_rtcp_avgJitter"].Set(float64(i.Data.QOS.OutRtcpAvgJitter))
	prom.GaugeMetrics[*name+"_inc_rtcp_maxJitter"].Set(float64(i.Data.QOS.IncRtcpMaxJitter))
	prom.GaugeMetrics[*name+"_out_rtcp_maxJitter"].Set(float64(i.Data.QOS.OutRtcpMaxJitter))
	prom.GaugeMetrics[*name+"_inc_rtcp_avgLat"].Set(float64(i.Data.QOS.IncRtcpAvgLat))
	prom.GaugeMetrics[*name+"_out_rtcp_avgLat"].Set(float64(i.Data.QOS.OutRtcpAvgLat))
	prom.GaugeMetrics[*name+"_inc_rtcp_maxLat"].Set(float64(i.Data.QOS.IncRtcpMaxLat))
	prom.GaugeMetrics[*name+"_out_rtcp_maxLat"].Set(float64(i.Data.QOS.OutRtcpMaxLat))
	if i.Data.QOS.EndTimeSec != 0 && i.Data.QOS.BeginTimeSec != 0 {
		prom.GaugeMetrics[*name+"_duration"].Set(float64(i.Data.QOS.EndTimeSec - i.Data.QOS.BeginTimeSec))
	}
}
