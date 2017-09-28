package main

import (
	"bytes"
	"fmt"
)

// SendStatsD creates a map with QOS or SIP stats which will
// be converted into statsd compatible strings seperated by '\n'
func (conn *Connections) SendStatsD(i *IPFIX, s string) {
	buf := new(bytes.Buffer)
	switch s {
	case "QOS":
		fields := i.mapMetricQOS()
		for metric, value := range fields {
			buf.Write([]byte(fmt.Sprintf("%s.%s:%.2f|h\n", *name, metric, value)))
		}
	}
	_, err := conn.StatsD.Write(buf.Bytes())
	checkErr(err)
}
