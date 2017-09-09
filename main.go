package main

import (
	"flag"
	"log"
	"net"
	/*
		"net"
		_ "net/http/pprof"
	*/
	"os"
)

var (
	addr     = flag.String("l", ":4739", "IPFIX listen address")
	gaddr    = flag.String("g", "", "Graylog gelf UDP server address")
	maddr    = flag.String("m", "", "MySQL server address")
	muser    = flag.String("mu", "", "MySQL user")
	mpass    = flag.String("mp", "", "MySQL password")
	mdb      = flag.String("md", "", "MySQL database")
	haddr    = flag.String("H", "", "Homer UDP server address")
	hepicQOS = flag.Bool("HQ", false, "Send hepic QOS Stats")
	iaddr    = flag.String("I", "", "InfluxDB http server address")
	saddr    = flag.String("s", "", "StatsD UDP server address")
	name     = flag.String("n", "sbc", "SBC name")
	hepPw    = flag.String("P", "myhep", "HEP capture password")
	debug    = flag.Bool("d", false, "Debug output to stdout")
	verbose  = flag.Bool("v", false, "Verbose output to stdout")
)

func main() {
	//go http.ListenAndServe(":8080", http.DefaultServeMux)
	flag.Parse()
	f, err := os.OpenFile("horaclifix.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	checkCritErr(err)
	defer f.Close()
	log.SetOutput(f)

	log.Printf("Start horaclifix interfaces IPFIX:%v Homer:%v Graylog:%v\n", *addr, *haddr, *gaddr)

	laddr, err := net.ResolveTCPAddr("tcp", *addr)
	checkCritErr(err)

	listener, err := net.ListenTCP("tcp", laddr)
	checkCritErr(err)

	for {
		conn, err := listener.AcceptTCP()
		checkCritErr(err)
		go Read(conn)
	}
}
