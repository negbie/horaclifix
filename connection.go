package main

import (
	"log"
	"net"
	"time"
)

func NewExternalConnections() *Connections {
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
			BatchSize:    256,
			FlushTimeout: 5 * time.Second,
			ErrorFunc:    checkErr,
		})
		checkCritErr(err)
		conn.Influx = iconn
	}

	if *gaddr != "" {
		gconn, err := net.Dial("tcp", *gaddr)
		checkCritErr(err)
		conn.Graylog = gconn
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

func CloseExternalConnections(c *Connections) {
	if *gaddr != "" {
		log.Printf("Close Graylog connection to %v\n", c.Graylog.RemoteAddr())
		err := c.Graylog.Close()
		checkErr(err)
	}
	if *maddr != "" {
		log.Printf("Close MySQL connection")
		err := c.MySQL.conn.Close()
		checkErr(err)
	}
	if *haddr != "" {
		log.Printf("Close Homer connection to %v\n", c.Homer.RemoteAddr())
		err := c.Homer.Close()
		checkErr(err)
	}
	if *saddr != "" {
		log.Printf("Close StatsD connection to %v\n", c.StatsD.RemoteAddr())
		err := c.StatsD.Close()
		checkErr(err)
	}
}
