package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/negbie/horaclifix/receiver"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {

	var (
		addr = flag.String("addr", ":4739", "Port to use to receive IPFIX packets")
	)

	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		os.Exit(1)
	}

	checkError(err)
	for {
		conn, err := listener.Accept()
		checkError(err)
		go receiver.Start(conn)
	}
}
