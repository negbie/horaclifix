package main

import (
	"crypto/tls"
	"fmt"
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
		conn.Graylog.TCP = gconn
		conn.Graylog.RWMutex = new(sync.RWMutex)
	}
	if *haddr != "" {
		if *network == "udp" {
			udpHAddr, err := net.ResolveUDPAddr("udp", *haddr)
			checkCritErr(err)

			hconn, err := net.DialUDP("udp", nil, udpHAddr)
			checkCritErr(err)
			conn.Homer.UDP = hconn
		} else if *network == "tcp" {
			tcpHAddr, err := net.ResolveTCPAddr("tcp", *haddr)
			checkCritErr(err)

			hconn, err := net.DialTCP("tcp", nil, tcpHAddr)
			checkCritErr(err)
			conn.Homer.TCP = hconn
			conn.Homer.RWMutex = new(sync.RWMutex)
		} else if *network == "tls" {
			hconn, err := tls.Dial("tcp", *haddr, &tls.Config{InsecureSkipVerify: true})
			checkCritErr(err)
			conn.Homer.TLS = hconn
			conn.Homer.RWMutex = new(sync.RWMutex)
		} else {
			checkCritErr(fmt.Errorf("Not supported network type %s", *network))
		}

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
		err := conn.Graylog.TCP.Close()
		checkErr(err)
	}
	if *maddr != "" {
		log.Printf("Close MySQL connection.\n")
		err := conn.MySQL.conn.Close()
		checkErr(err)
	}
	if *haddr != "" {
		log.Printf("Close Homer connection.\n")
		if *network == "udp" {
			err := conn.Homer.UDP.Close()
			checkErr(err)
		} else if *network == "tcp" {
			err := conn.Homer.TCP.Close()
			checkErr(err)
		} else if *network == "tls" {
			err := conn.Homer.TLS.Close()
			checkErr(err)
		}
	}
	if *saddr != "" {
		log.Printf("Close StatsD connection.\n")
		err := conn.StatsD.Close()
		checkErr(err)
	}
}
