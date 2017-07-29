package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/negbie/siprocket"
)

// SendLog will encode the QOS & SIP maps to json
// and send them over UDP to Graylog
func (conn Connections) SendLog(i *IPFIX, s string) {
	buf := new(bytes.Buffer)
	sipMSG := siprocket.Parse(i.Data.SIP.SipMsg)
	//siprocket.PrintSipStruct(&sipMSG)
	var gLog []byte
	var err error
	switch s {
	case "SIP":
		buf.Write([]byte(fmt.Sprintf("{\"%s\":\"%s\",", "version", "1.1")))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "host", *name)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "short_message", i.Data.SIP.SipMsg)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%d\",", "level", 5)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_id", sipMSG.CallId.Value)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_from", sipMSG.From.User)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_to", sipMSG.To.User)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_method", sipMSG.Req.Method)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_statusCode", sipMSG.Req.StatusCode)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_ua", sipMSG.Ua.Value)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_srcIp", stringIPv4(i.Data.SIP.SrcIP))))
		buf.Write([]byte(fmt.Sprintf("\"%s\":\"%s\",", "_dstIp", stringIPv4(i.Data.SIP.DstIP))))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_srcPort", i.Data.SIP.SrcPort)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_dstPort", i.Data.SIP.DstPort)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_ipLen", i.Data.SIP.IPlen)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_udpLen", i.Data.SIP.UDPlen)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_vl", i.Data.SIP.IntVlan)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_tos", i.Data.SIP.TOS)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_tlen", i.Data.SIP.TLen)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_tid", i.Data.SIP.TID)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_tflags", i.Data.SIP.TFlags)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d,", "_ttl", i.Data.SIP.TTL)))
		buf.Write([]byte(fmt.Sprintf("\"%s\":%d}", "_tproto", i.Data.SIP.TProto)))
		gLog = buf.Bytes()
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
