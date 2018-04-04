package main

import (
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

var gErrCnt int

// SendLog will encode the QOS & SIP maps to json
// and send them over UDP to Graylog
func (conn *Connections) SendLog(i *IPFIX, s string) {
	var err error
	var gLog []byte

	switch s {
	case "SIP":
		bb := bytebufferpool.Get()
		defer bytebufferpool.Put(bb)
		//gLog, err = json.Marshal(i.mapLogSIP())
		gLog = formGelf(bb, i)
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

func formGelf(b *bytebufferpool.ByteBuffer, i *IPFIX) []byte {
	b.WriteString("{")
	b.WriteString("\"version\":1.1")
	b.WriteString(",\"host\":\"")
	b.WriteString(*name)
	b.WriteString("\",\"short_message\":\"")
	b.WriteString(i.SIP.SipMsg.Msg)
	b.WriteString("\",\"level\":5")
	b.WriteString(",\"_id\":\"")
	b.WriteString(i.SIP.SipMsg.CallID)
	b.WriteString("\",\"_from\":\"")
	b.WriteString(i.SIP.SipMsg.FromUser)
	b.WriteString("\",\"_to\":\"")
	b.WriteString(i.SIP.SipMsg.ToUser)
	b.WriteString("\",\"_paiUser\":\"")
	b.WriteString(i.SIP.SipMsg.PaiUser)
	b.WriteString("\",\"_paiHost\":\"")
	b.WriteString(i.SIP.SipMsg.PaiHost)
	b.WriteString("\",\"_method\":\"")
	b.WriteString(i.SIP.SipMsg.StartLine.Method)
	b.WriteString("\",\"_response\":\"")
	b.WriteString(i.SIP.SipMsg.StartLine.Resp)
	b.WriteString("\",\"_ua\":\"")
	b.WriteString(i.SIP.SipMsg.UserAgent)
	b.WriteString("\",\"_srcIP\":\"")
	b.WriteString(i.SIP.SrcIPString)
	b.WriteString("\",\"_dstIP\":\"")
	b.WriteString(i.SIP.DstIPString)
	b.WriteString("\",\"_intVlan\":")
	b.WriteString(strconv.Itoa(int(i.SIP.IntVlan)))
	b.WriteString(",\"_udpLen\":")
	b.WriteString(strconv.Itoa(int(i.SIP.UDPlen)))
	b.WriteString("}")
	return b.Bytes()
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
