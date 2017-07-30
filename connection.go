package main

import (
	"crypto/tls"
	"encoding/binary"
	"net"
	"time"
)

func NewConnections() *Connections {
	conn := new(Connections)
	if *baddr != "" {
		sconn, err := net.Dial("tcp", *baddr)
		checkCritErr(err)
		conn.Banshee = sconn
	}
	if *gaddr != "" {
		gconn, err := net.Dial("udp", *gaddr)
		checkCritErr(err)
		conn.Graylog = gconn
	}
	if *gtaddr != "" {
		gtconn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", *gtaddr, nil)
		checkCritErr(err)
		conn.GraylogTLS = gtconn
	}
	if *haddr != "" {
		hconn, err := net.Dial("udp", *haddr)
		checkCritErr(err)
		conn.Homer = hconn
	}
	if *saddr != "" {
		sconn, err := net.Dial("udp", *saddr)
		checkCritErr(err)
		conn.StatsD = sconn
	}
	return conn
}

// stringIPv4 converts a ipv4 unit32 into a string
func stringIPv4(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}
