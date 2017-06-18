package main

import (
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

	listener, err := net.Listen("tcp", ":4739")
	checkError(err)
	for {
		conn, err := listener.Accept()
		checkError(err)
		go receiver.SyncClient(conn)
	}
}
