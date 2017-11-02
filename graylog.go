package main

import (
	"encoding/json"
	"net"
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
		/*	err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogQOS())
			checkErr(err)
		*/
	}
	// Graylog frame delimiter
	data := append(gLog, '\n', byte(0))

	_, err = conn.writeTCP(data)
	checkErr(err)
}

func (conn *Connections) reconnect() error {
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

func (conn *Connections) writeTCP(b []byte) (int, error) {
	conn.Graylog.RLock()
	defer conn.Graylog.RUnlock()

	if conn.Graylog.disconnected {
		conn.Graylog.RUnlock()
		if err := conn.reconnect(); err != nil {
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
