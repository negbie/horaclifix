package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	//"net"
	//_ "net/http/pprof"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const version = "horaclifix 1.1.3"

var (
	addr        = flag.String("l", ":4739", "IPFIX TCP listen address")
	gaddr       = flag.String("g", "", "Graylog gelf TCP server address")
	maddr       = flag.String("m", "", "MySQL TCP server address")
	muser       = flag.String("mu", "", "MySQL user")
	mpass       = flag.String("mp", "", "MySQL password")
	haddr       = flag.String("H", "", "Homer UDP server address")
	hepicQOS    = flag.Bool("HQ", false, "Send hepic QOS Stats")
	iaddr       = flag.String("I", "", "InfluxDB HTTP server address")
	saddr       = flag.String("s", "", "StatsD UDP server address")
	paddr       = flag.String("p", "", "Prometheus address")
	name        = flag.String("n", "sbc", "SBC name")
	network     = flag.String("nt", "udp", "Network types are [udp, tcp, tls]")
	protobuf    = flag.Bool("protobuf", false, "Use Protobuf on wire")
	filter      = flag.String("di", "", "Discard SIP method")
	hepID       = flag.Int("N", 2004, "HEP capture node ID")
	hepPW       = flag.String("P", "myhep", "HEP capture password")
	debug       = flag.Bool("d", false, "Debug output to stdout")
	verbose     = flag.Bool("v", false, "Verbose output to stdout")
	showVersion = flag.Bool("V", false, "Show version")
)

func main() {
	//go http.ListenAndServe(":8080", http.DefaultServeMux)
	//trace.Start(os.Stdout)
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	ex, err := os.Executable()
	checkCritErr(err)
	exPath := filepath.Dir(ex)
	exPathName := exPath + "/" + "horaclifix_" + *name + ".log"
	f, err := os.OpenFile(exPathName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	checkCritErr(err)
	// Close file when this function ends
	defer func() {
		err := f.Close()
		checkErr(err)
	}()

	log.SetOutput(&lumberjack.Logger{
		Filename:   exPathName,
		MaxSize:    10, // mb
		MaxBackups: 7,
		MaxAge:     28, //days
		Compress:   true,
	})

	log.Printf("Start horaclifix interfaces IPFIX:%v Homer:%v Graylog:%v\n", *addr, *haddr, *gaddr)

	laddr, err := net.ResolveTCPAddr("tcp", *addr)
	checkCritErr(err)

	listener, err := net.ListenTCP("tcp", laddr)
	checkCritErr(err)

	regProm()

	for {
		conn, err := listener.AcceptTCP()
		checkCritErr(err)
		go Read(conn)
	}
}
