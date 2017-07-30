package main

import (
	"encoding/json"
	"log"
)

// SendLog will encode the QOS & SIP maps to json
// and send them over UDP to Graylog
func (conn Connections) SendLog(i *IPFIX, s string) {
	var gLog []byte
	var err error
	switch s {
	case "SIP":
		if gLog, err = json.Marshal(i.mapLogSIP()); err != nil {
			log.Println("SIP json.Marshal failed:", err, gLog)
		}
		/*	if err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogSIP()); err != nil {
			log.Println("SIP json.NewEncoder failed:", err, gLog)
		}*/
	case "QOS":
		if gLog, err = json.Marshal(i.mapLogQOS()); err != nil {
			log.Println("QOS json.Marshal failed:", err, gLog)
		}
		/*	if err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogQOS()); err != nil {
			log.Println("SIP json.NewEncoder failed:", err, gLog)
		}*/
	}
	// Graylog frame delimiter
	data := append(gLog, '\n', byte(0))
	conn.Graylog.Write(data)
}
