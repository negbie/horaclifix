package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var (
	addr  = flag.String("l", ":4739", "Host ipfix listen address")
	haddr = flag.String("H", "127.0.0.1:9060", "Homer server address")
	saddr = flag.String("s", "127.0.0.1:8125", "StatsD server address")
	debug = flag.Bool("d", false, "Debug output to stdout")
	gaddr = flag.String("g", "", "Graylog server address")
)

func checkErr(err error) {
	if err != nil {
		log.Println("Warning:", err)
	}
}

func checkCritErr(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}

var buffers = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 65536)
	},
}

// start handles incoming packets
func start(conn *net.TCPConn) {
	log.Printf("Handling new connection under %v\n", *addr)

	// Close connection when this function ends
	defer func() {
		log.Printf("Closing connection under %v\n", *addr)
		conn.Close()
	}()

	byts := buffers.Get().([]byte)

	for {
		blen, err := bufio.NewReader(conn).Read(byts)
		// Check for EOF and go out of this loop. Don't cut the connection. Mby we just rebooted the sbc
		if err == io.EOF {
			break
		}
		checkErr(err)

		if blen > 20 {
			parse(conn, byts[:blen])
		}
		buffers.Put(byts)
	}
}

func main() {
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
		go start(conn)
	}

}
