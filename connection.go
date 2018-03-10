package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func NewExtConns() *Connections {
	var err error
	conn := new(Connections)

	if *maddr != "" {
		conn.MySQL, err = newMySQLDB()
		checkCritErr(err)
	}
	if *iaddr != "" {
		iconn, err := NewInfluxClient(&InfluxClientConfig{
			Endpoint:     "http://" + *iaddr,
			Database:     "horaclifix",
			BatchSize:    300,
			FlushTimeout: 5 * time.Second,
			ErrorFunc:    checkErr,
		})
		checkCritErr(err)
		conn.Influx = iconn
	}
	if *gaddr != "" {
		tcpGAddr, err := net.ResolveTCPAddr("tcp", *gaddr)
		checkCritErr(err)

		gconn, err := net.DialTCP("tcp", nil, tcpGAddr)
		checkCritErr(err)
		conn.Graylog.TCPConn = gconn
		conn.Graylog.RWMutex = new(sync.RWMutex)
	}
	if *haddr != "" {
		udpHAddr, err := net.ResolveUDPAddr("udp", *haddr)
		checkCritErr(err)

		hconn, err := net.DialUDP("udp", nil, udpHAddr)
		checkCritErr(err)
		conn.Homer = hconn
	}
	if *saddr != "" {
		udpSAddr, err := net.ResolveUDPAddr("udp", *saddr)
		checkCritErr(err)

		sconn, err := net.DialUDP("udp", nil, udpSAddr)
		checkCritErr(err)
		conn.StatsD = sconn
	}

	return conn
}

func CloseExtConns(conn *Connections) {
	if *gaddr != "" {
		log.Printf("Close Graylog connection.\n")
		err := conn.Graylog.Close()
		checkErr(err)
	}
	if *maddr != "" {
		log.Printf("Close MySQL connection.\n")
		err := conn.MySQL.conn.Close()
		checkErr(err)
	}
	if *haddr != "" {
		log.Printf("Close Homer connection.\n")
		err := conn.Homer.Close()
		checkErr(err)
	}
	if *saddr != "" {
		log.Printf("Close StatsD connection.\n")
		err := conn.StatsD.Close()
		checkErr(err)
	}
}
