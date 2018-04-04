package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

var gErrCnt int

// SendLog will encode the QOS & SIP maps to json
// and send them over UDP to Graylog
func (conn *Connections) SendLog(i *IPFIX, s string) {
	var err error
	var gLog []byte

	switch s {
	case "SIP":
		//gLog, err = json.Marshal(i.mapLogSIP())
		gLog = []byte(fmt.Sprintf("{\"version\":1.1,\"host\":\"%s\",\"short_message\":%s,\"level\":5,\"_id\":\"%s\",\"_from\":\"%s\",\"_to\":\"%s\",\"_paiUser\":\"%s\",\"_paiHost\":\"%s\","+
			"\"_method\":\"%s\",\"_response\":\"%s\",\"_ua\":\"%s\",\"_srcIP\":\"%s\",\"_dstIP\":\"%s\",\"_srcPort\":%d,\"_dstPort\":%d,\"_intVlan\":%d,\"_udpLen\":%d}",
			*name, strconv.Quote(i.SIP.SipMsg.Msg), i.SIP.SipMsg.CallID, i.SIP.SipMsg.FromUser, i.SIP.SipMsg.ToUser, i.SIP.SipMsg.PaiUser, i.SIP.SipMsg.PaiHost,
			i.SIP.SipMsg.StartLine.Method, i.SIP.SipMsg.StartLine.Resp, i.SIP.SipMsg.UserAgent,
			stringIPv4(i.SIP.SrcIP), stringIPv4(i.SIP.DstIP), i.SIP.SrcPort, i.SIP.DstPort, i.SIP.IntVlan, i.SIP.UDPlen))

	case "QOS":
		gLog, err = json.Marshal(i.mapLogQOS())
		checkErr(err)
	}

	// Graylog frame delimiter
	data := append(gLog, '\n', byte(0))

	_, err = conn.Graylog.TCP.Write(data)
	if err != nil {
		gErrCnt++
		if gErrCnt%128 == 0 {
			log.Printf("[WARN] <%s> %s\n", *name, err)
			gErrCnt = 0
			reWriteGraylog(conn, data)
		}
	}
}

func reConnectGraylog(conn *Connections) error {
	conn.Graylog.Lock()
	defer conn.Graylog.Unlock()

	raddr := conn.Graylog.TCP.RemoteAddr()
	gconn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
	if err != nil {
		return err
	}
	conn.Graylog.TCP.Close()
	conn.Graylog.TCP = gconn
	return nil
}

func reWriteGraylog(conn *Connections, b []byte) error {
	conn.Graylog.RLock()
	defer conn.Graylog.RUnlock()

	if conn.Graylog.disconnected {
		conn.Graylog.RUnlock()
		if err := reConnectGraylog(conn); err != nil {
			conn.Graylog.disconnected = true
			conn.Graylog.RLock()
			return err
		}
		conn.Graylog.disconnected = false
		conn.Graylog.RLock()
	}
	_, err := conn.Graylog.TCP.Write(b)
	if err == nil {
		return err
	}
	conn.Graylog.disconnected = true
	return err
}
