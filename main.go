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
		port = flag.String("port", ":4739", "Port to use to receive IPFIX packets")
	)

	flag.Parse()
	listener, err := net.Listen("tcp", *port)
	checkError(err)
	for {
		conn, err := listener.Accept()
		checkError(err)
		go receiver.SyncClient(conn)
	}
}
