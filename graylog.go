package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
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
	if *aaddr != "" {
		err = conn.AmqpChannel.Publish(
			"",                  // exchange
			conn.AmqpQueue.Name, // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{Body: data},
		)
		checkErr(err)
	}
	conn.Graylog.Write(data)
}
