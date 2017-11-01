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
		/*	err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogQOS()); err != nil {
			checkErr(err)
		*/
	}
	// Graylog frame delimiter
	data := append(gLog, '\n', byte(0))

	_, err = conn.WriteTCP(data)
	checkErr(err)
}

func (c *Connections) reconnect() error {
	c.Graylog.Lock()
	defer c.Graylog.Unlock()

	raddr := c.Graylog.TCPConn.RemoteAddr()
	gconn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
	if err != nil {
		return err
	}

	c.Graylog.TCPConn.Close()
	c.Graylog.TCPConn = gconn
	return nil
}

func (c *Connections) WriteTCP(b []byte) (int, error) {
	c.Graylog.RLock()
	defer c.Graylog.RUnlock()

	if c.Graylog.disconnected {
		c.Graylog.RUnlock()
		if err := c.reconnect(); err != nil {
			c.Graylog.disconnected = true
			c.Graylog.RLock()
			return -1, err
		}
		c.Graylog.disconnected = false
		c.Graylog.RLock()
	}
	n, err := c.Graylog.TCPConn.Write(b)
	if err == nil {
		return n, err
	}
	c.Graylog.disconnected = true
	return -1, err
}
