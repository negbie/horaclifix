package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func NewExtConns() *Connections {
	var err error
	conns := new(Connections)

	if *maddr != "" {
		conns.MySQL, err = newMySQLDB()
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
		conns.Influx = iconn
	}

	if *gaddr != "" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *gaddr)
		checkCritErr(err)

		gconn, err := net.DialTCP("tcp", nil, tcpAddr)
		conns.Graylog.TCPConn = gconn
		conns.Graylog.RWMutex = new(sync.RWMutex)
		checkCritErr(err)
	}

	if *haddr != "" {
		hconn, err := net.Dial("udp", *haddr)
		checkCritErr(err)
		conns.Homer = hconn
	}
	if *saddr != "" {
		sconn, err := net.Dial("udp", *saddr)
		checkCritErr(err)
		conns.StatsD = sconn
	}
	return conns
}

func CloseExtConns(c *Connections) {
	if *gaddr != "" {
		log.Printf("Close Graylog connection.\n")
		err := c.Graylog.Close()
		checkErr(err)
	}
	if *maddr != "" {
		log.Printf("Close MySQL connection.\n")
		err := c.MySQL.conn.Close()
		checkErr(err)
	}
	if *haddr != "" {
		log.Printf("Close Homer connection.\n")
		err := c.Homer.Close()
		checkErr(err)
	}
	if *saddr != "" {
		log.Printf("Close StatsD connection.\n")
		err := c.StatsD.Close()
		checkErr(err)
	}
}
