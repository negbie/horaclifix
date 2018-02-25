package main

import "fmt"

func (conn *Connections) Send(i *IPFIX, s string) {
	switch s {
	case "SIP":
		if *haddr != "" {
			conn.SendHep(i, "SIP")
		}
		if *gaddr != "" {
			conn.SendLog(i, "SIP")
		}
		if *debug {
			fmt.Println("SIP output:")
			fmt.Printf("%s\n", i.Data.SIP.SipMsg)
		}

	default:
		// Send only QOS stats with meaningful values
		if i.Data.QOS.IncMos > 0 || i.Data.QOS.OutMos > 0 {

			if *haddr != "" {
				if *hepicQOS {
					conn.SendHep(i, "incQOS")
					conn.SendHep(i, "outQOS")
					conn.SendHep(i, "incMOS")
					conn.SendHep(i, "outMOS")
				} else {
					conn.SendHep(i, "allQOS")
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
		}
		if *gaddr != "" {
			conn.SendLog(i, "QOS")
		}
		if *maddr != "" {
			conn.SendMySQL(i, "QOS")
		}
	}
}
