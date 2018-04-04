package main

import "fmt"

func (conn *Connections) Send(i *IPFIX) {
	if i.SIP.RawMsg != nil {
		if *haddr != "" {
			conn.SendHep(i, "SIP")
		}
		if *gaddr != "" {
			conn.SendLog(i, "SIP")
		}
		if *debug {
			fmt.Println("SIP output:")
			fmt.Printf("%s\n", i.SIP.RawMsg)
		}
	} else if i.QOS.IncMos > 0 || i.QOS.OutMos > 0 {
		if *haddr != "" {
			if *hepicQOS {
				conn.SendHep(i, "incQOS")
				conn.SendHep(i, "outQOS")
				conn.SendHep(i, "incMOS")
				conn.SendHep(i, "outMOS")
			} else {
				conn.SendHep(i, "QOS")
			}
		}
		if *iaddr != "" {
			conn.Influx.Send(i, "QOS")
		}
		if *paddr != "" {
			conn.SendMetric(i, "QOS")
		}
		if *saddr != "" {
			conn.SendStatsD(i, "QOS")
		}
		if *gaddr != "" {
			conn.SendLog(i, "QOS")
		}
		if *maddr != "" {
			conn.SendMySQL(i, "QOS")
		}
	}
}
