package main

import (
	"encoding/json"
)

// SendLog will encode the QOS & SIP maps to json
// and send them over UDP to Graylog
func (conn *Connections) SendLog(i *IPFIX, s string) {
	var err error
	var gLog []byte

	switch s {
	case "SIP":
		gLog, err = json.Marshal(i.mapLogSIP())
		checkErr(err)
		/*	err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogSIP())
			checkErr(err)
		*/
	case "QOS":
		gLog, err = json.Marshal(i.mapLogQOS())
		checkErr(err)
		/*	err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogQOS()); err != nil {
			checkErr(err)
		*/
	}
	// Graylog frame delimiter
	data := append(gLog, '\n', byte(0))
	_, err = conn.Graylog.Write(data)
	checkErr(err)
}
