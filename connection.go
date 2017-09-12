package main

import (
	"log"
	"net"
	"time"
)

func NewExternalConnections() *Connections {
	conn := new(Connections)
	var err error

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
		gconn, err := net.Dial("udp", *gaddr)
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
		if err != nil {
			log.Printf("Close Graylog connection error: %s", err)
		}
	}
	if *maddr != "" {
		log.Println("Close MySQL connection")
		err := c.MySQL.conn.Close()
		if err != nil {
			log.Printf("Close MySQL connection error: %s", err)
		}
	}
	if *haddr != "" {
		log.Printf("Close Homer connection to %v\n", c.Homer.RemoteAddr())
		err := c.Homer.Close()
		if err != nil {
			log.Printf("Close Homer connection error: %s", err)
		}
	}
	if *saddr != "" {
		log.Printf("Close StatsD connection to %v\n", c.StatsD.RemoteAddr())
		err := c.StatsD.Close()
		if err != nil {
			log.Printf("Close StatsD connection error: %s", err)
		}
	}
}
