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
		gLog = []byte(fmt.Sprintf("{\"version\":1.1,\"host\":\"%s\",\"short_message\":%s,\"level\":5,\"_id\":\"%s\",\"_from\":\"%s\",\"_to\":\"%s\",\"_method\":\"%s\",\"_response\":\"%s\",\"_ua\":\"%s\",\"_srcIP\":\"%s\",\"_dstIP\":\"%s\",\"_srcPort\":%d,\"_dstPort\":%d,\"_intVlan\":%d,\"_udpLen\":%d}",
			*name, strconv.Quote(i.SIP.SipMsg.Msg), i.SIP.SipMsg.CallID, i.SIP.SipMsg.FromUser, i.SIP.SipMsg.ToUser,
			i.SIP.SipMsg.StartLine.Method, i.SIP.SipMsg.StartLine.Resp, i.SIP.SipMsg.UserAgent,
			stringIPv4(i.SIP.SrcIP), stringIPv4(i.SIP.DstIP), i.SIP.SrcPort, i.SIP.DstPort, i.SIP.IntVlan, i.SIP.UDPlen))

	case "QOS":
		gLog, err = json.Marshal(i.mapLogQOS())
		checkErr(err)
	}

	// Graylog frame delimiter
	data := append(gLog, '\n', byte(0))

	_, err = conn.Graylog.TCPConn.Write(data)
	if err != nil {
		gErrCnt++
		if gErrCnt%128 == 0 {
			log.Printf("[WARN] <%s> %s\n", *name, err)
			gErrCnt = 0
			conn.ReWrite(data)
		}
	}
}

func reconnect(conn *Connections) error {
	conn.Graylog.Lock()
	defer conn.Graylog.Unlock()

	raddr := conn.Graylog.TCPConn.RemoteAddr()
	gconn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
	if err != nil {
		return err
	}

	conn.Graylog.TCPConn.Close()
	conn.Graylog.TCPConn = gconn
	return nil
}

func (conn *Connections) ReWrite(b []byte) (int, error) {
	conn.Graylog.RLock()
	defer conn.Graylog.RUnlock()

	if conn.Graylog.disconnected {
		conn.Graylog.RUnlock()
		if err := reconnect(conn); err != nil {
			conn.Graylog.disconnected = true
			conn.Graylog.RLock()
			return -1, err
		}
		conn.Graylog.disconnected = false
		conn.Graylog.RLock()
	}
	n, err := conn.Graylog.TCPConn.Write(b)
	if err == nil {
		return n, err
	}
	conn.Graylog.disconnected = true
	return -1, err
}
