package main

import (
	"fmt"
)

func (conn *Connections) Send(ich chan *IPFIX) {
	var (
		msg *IPFIX
		ok  bool
	)

	for {
		msg, ok = <-ich
		if !ok {
			break
		}

		if msg.Data.SIP.CallID != nil {
			if *haddr != "" {
				conn.SendHep(msg, "SIP")
			}
			if *gaddr != "" {
				conn.SendLog(msg, "SIP")
			}
			if *debug {
				fmt.Println("SIP output:")
				fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
			}
		} else if msg.Data.QOS.IncMos > 0 || msg.Data.QOS.OutMos > 0 {
			if *haddr != "" {
				if *hepicQOS {
					conn.SendHep(msg, "incQOS")
					conn.SendHep(msg, "outQOS")
					conn.SendHep(msg, "incMOS")
					conn.SendHep(msg, "outMOS")
				} else {
					conn.SendHep(msg, "allQOS")
				}
			}
			if *iaddr != "" {
				conn.Influx.Send(msg, "QOS")
			}
			if *paddr != "" {
				conn.SendMetric(msg, "QOS")
			}
			if *saddr != "" {
				conn.SendStatsD(msg, "QOS")
			}

			if *gaddr != "" {
				conn.SendLog(msg, "QOS")
			}
			if *maddr != "" {
				conn.SendMySQL(msg, "QOS")
			}
		}
	}
}
