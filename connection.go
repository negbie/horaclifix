package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewExtConns() *Connections {
	var err error
	conn := new(Connections)

	if *maddr != "" {
		conn.MySQL, err = newMySQLDB()
		checkCritErr(err)
	}

	if *iaddr != "" {
		iconn, err := NewInfluxClient(&InfluxClientConfig{
			Endpoint:     "http://" + *iaddr,
			Database:     "horaclifix",
			BatchSize:    300,
			FlushTimeout: 5 * time.Second,
			ErrorFunc:    checkErr,
		})
		checkCritErr(err)
		conn.Influx = iconn
	}

	if *gaddr != "" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *gaddr)
		checkCritErr(err)

		gconn, err := net.DialTCP("tcp", nil, tcpAddr)
		checkCritErr(err)
		conn.Graylog.TCPConn = gconn
		conn.Graylog.RWMutex = new(sync.RWMutex)

	}

	if *haddr != "" {
		hconn, err := net.Dial("udp", *haddr)
		checkCritErr(err)
		conn.Homer = hconn
	}

	if *saddr != "" {
		sconn, err := net.Dial("udp", *saddr)
		checkCritErr(err)
		conn.StatsD = sconn
	}

	if *paddr != "" {
		conn.Prometheus.GaugeMetrics = map[string]prometheus.Gauge{}

		conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_mos"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_mos", Help: "Incoming RTP MOS"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtp_mos"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_mos", Help: "Outgoing RTP MOS"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_rval"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_rval", Help: "Incoming RTP rVal"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtp_rval"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_rval", Help: "Outgoing RTP rVal"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_packets", Help: "Incoming RTP packets"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_packets", Help: "Outgoing RTP packets"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_packets", Help: "Incoming RTCP packets"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_packets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_packets", Help: "Outgoing RTCP packets"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_lostPackets", Help: "Incoming RTP lostPackets"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_lostPackets", Help: "Outgoing RTP lostPackets"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_lostPackets", Help: "Incoming RTCP lostPackets"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_lostPackets"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_lostPackets", Help: "Outgoing RTCP lostPackets"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_avgJitter", Help: "Incoming RTP avgJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_avgJitter", Help: "Outgoing RTP avgJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtp_maxJitter", Help: "Incoming RTP maxJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtp_maxJitter", Help: "Outgoing RTP maxJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_avgJitter", Help: "Incoming RTCP avgJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_avgJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_avgJitter", Help: "Outgoing RTCP avgJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_maxJitter", Help: "Incoming RTCP maxJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_maxJitter"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_maxJitter", Help: "Outgoing RTCP maxJitter"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_avgLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_avgLat", Help: "Incoming RTCP avgLat"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_avgLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_avgLat", Help: "Outgoing RTCP avgLat"})
		conn.Prometheus.GaugeMetrics[*name+"_inc_rtcp_maxLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_inc_rtcp_maxLat", Help: "Incoming RTCP maxLat"})
		conn.Prometheus.GaugeMetrics[*name+"_out_rtcp_maxLat"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_out_rtcp_maxLat", Help: "Outgoing RTCP maxLat"})
		conn.Prometheus.GaugeMetrics[*name+"_duration"] = prometheus.NewGauge(prometheus.GaugeOpts{Name: *name + "_duration", Help: "Call duration"})
		fmt.Println("HUHU")
		for k := range conn.Prometheus.GaugeMetrics {
			log.Printf("register prometheus gaugeMetric %s", k)
			prometheus.MustRegister(conn.Prometheus.GaugeMetrics[k])
		}

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err = http.ListenAndServe(*paddr, nil)
			checkCritErr(err)
		}()
	}
	return conn
}

func CloseExtConns(conn *Connections) {
	if *gaddr != "" {
		log.Printf("Close Graylog connection.\n")
		err := conn.Graylog.Close()
		checkErr(err)
	}
	if *maddr != "" {
		log.Printf("Close MySQL connection.\n")
		err := conn.MySQL.conn.Close()
		checkErr(err)
	}
	if *haddr != "" {
		log.Printf("Close Homer connection.\n")
		err := conn.Homer.Close()
		checkErr(err)
	}
	if *saddr != "" {
		log.Printf("Close StatsD connection.\n")
		err := conn.StatsD.Close()
		checkErr(err)
	}
}
