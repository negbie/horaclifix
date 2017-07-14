package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// SendStatsd creates a map with QOS or SIP stats which will
// be converted into statsd compatible strings seperated by '\n'
func (i *IPFIX) SendStatsd(s string) {
	var metrics []string
	var mapQOS map[string]interface{}
	var mapCalls map[string]interface{}
	switch s {
	case "SIP":

	case "QOS":
		mapQOS = map[string]interface{}{
			"QOS.IncMos": i.Data.QOS.IncMos,
			"QOS.OutMos": i.Data.QOS.OutMos,

			"QOS.IncRtpPackets": i.Data.QOS.IncRtpPackets,
			"QOS.OutRtpPackets": i.Data.QOS.OutRtpPackets,

			"QOS.IncRtcpAvgLat": i.Data.QOS.IncRtcpAvgLat,
			"QOS.OutRtcpAvgLat": i.Data.QOS.OutRtcpAvgLat,

			"QOS.IncRtpAvgJitter": i.Data.QOS.IncRtpAvgJitter,
			"QOS.OutRtpAvgJitter": i.Data.QOS.OutRtpAvgJitter,

			"QOS.IncRtpMaxJitter": i.Data.QOS.IncRtpMaxJitter,
			"QOS.OutRtpMaxJitter": i.Data.QOS.OutRtpMaxJitter,

			"QOS.IncRtcpAvgJitter": i.Data.QOS.IncRtcpAvgJitter,
			"QOS.OutRtcpAvgJitter": i.Data.QOS.OutRtcpAvgJitter,

			"QOS.IncRtcpMaxJitter": i.Data.QOS.IncRtcpMaxJitter,
			"QOS.OutRtcpMaxJitter": i.Data.QOS.OutRtcpMaxJitter,

			"QOS.IncRtpLostPackets": i.Data.QOS.IncRtpLostPackets,
			"QOS.OutRtpLostPackets": i.Data.QOS.OutRtpLostPackets,

			"QOS.IncRtcpLostPackets": i.Data.QOS.IncRtcpLostPackets,
			"QOS.OutRtcpLostPackets": i.Data.QOS.OutRtcpLostPackets,
		}

		mapCalls = map[string]interface{}{
			"QOS.IncCallID": i.Data.QOS.IncCallID,
		}
	}

	for metric, value := range mapQOS {
		metrics = append(metrics, fmt.Sprintf("%s:%d|h", metric, value))
	}
	for metric, value := range mapCalls {
		metrics = append(metrics, fmt.Sprintf("%s:%d|s", metric, value))
	}
	stats := strings.Join(metrics, "\n")

	if conn, err := net.Dial("udp", *saddr); err == nil {
		io.WriteString(conn, stats)
		conn.Close()
	}
}

/*
var queue = make(chan string, 200)

func init() {
	go statsdSender()
}

// StatCount sends name:value|c where value is a positive or negative integer number.
// :value can be omitted and statsd will assume it is 1.
func StatCount(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|c", metric, value)
}

// StatMeter sends name:value|m where value is a positive or negative integer number.
// :value can be omitted and statsd will assume it is 1.
func StatMeter(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|m", metric, value)
}

// StatTime sends name:value|ms, where value a time in ms,
// statsd reports min, max, average, sum, average of 95th percentile, median and standard deviation
// and the total number of times it was updated (events).
func StatTime(metric string, duration time.Duration) {
	queue <- fmt.Sprintf("%s:%d|ms", metric, duration/1e6)
}

// StatHist sends name:value|h, where value is any decimal number,
// statsd reports min, max, average, sum, average of 95th percentile, median and standard deviation
// and the total number of times it was updated (events).
func StatHist(metric string, hist int) {
	queue <- fmt.Sprintf("%s:%d|h", metric, hist)
}

// StatGauge sends a constant data type. They are not subject to averaging
// and they dont change unless you change them. That is, once you set a gauge value,
// it will be a flat line on the graph until you change it again.
func StatGauge(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|g", metric, value)
}

// StatNumSet sends name:value|s where value is a positive or negative integer number.
func StatNumSet(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|s", metric, value)
}

// StatStrSet sends name:value|s where value is a string.
func StatStrSet(metric string, value string) {
	queue <- fmt.Sprintf("%s:%s|s", metric, value)
}

func statsdSender() {
	for s := range queue {
		if conn, err := net.Dial("udp", *saddr); err == nil {
			io.WriteString(conn, s)
			conn.Close()
		}
	}
}
*/
