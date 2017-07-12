package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	addr    = flag.String("l", ":4739", "Host ipfix listen address")
	haddr   = flag.String("H", "127.0.0.1:9060", "Homer server address")
	saddr   = flag.String("s", "127.0.0.1:8125", "StatsD server address")
	debug   = flag.Bool("d", false, "Debug output to stdout")
	verbose = flag.Bool("v", false, "Debug output to stdout")
	gaddr   = flag.String("g", "", "Graylog server address")
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

/*
// Template for a sync.Pool buffer
var buffers = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 65536)
	},
}
packet := buffers.Get().([]byte)
buffers.Put(packet)
*/

// start handles incoming packets
func start(conn *net.TCPConn) {
	log.Printf("Handling new connection under %v\n", *addr)

	// Close connection when this function ends
	defer func() {
		log.Printf("Closing connection under %v\n", *addr)
		conn.Close()
	}()

	r := bufio.NewReader(conn)
	header := make([]byte, 20)
	var (
		version int
		dataLen int
		setID   int
		dataSet []byte
	)

	for {
		if _, err := io.ReadFull(r, header); err == nil {
			version = int(uint16(header[0])<<8 + uint16(header[1]))
			dataLen = int(uint16(header[2])<<8 + uint16(header[3]))
			setID = int(uint16(header[16])<<8 + uint16(header[17]))
			dataSet = make([]byte, dataLen-len(header))

		} else {
			checkErr(err)
		}
		if _, err := io.ReadFull(r, dataSet); err == nil {
			data := append(header, dataSet...)

			if *verbose {
				fmt.Println("########################################################################")
				fmt.Printf("Headerversion: %d, Headerlength: %d, SetID: %d\n", version, dataLen, setID)
				fmt.Println("Header in raw:", header)
			}
			if *debug {
				fmt.Println("Hexdump output:")
				fmt.Printf("%s\n", hex.Dump(data))
			}
			switch setID {
			case 256:
				SendHandshake(conn, data)
			case 258:
				msg := NewRecSipUDP(data)
				NewSipHEP(msg)
				if *gaddr != "" {
					msg.SendLog("SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 259:
				msg := NewSendSipUDP(data)
				NewSipHEP(msg)
				if *gaddr != "" {
					msg.SendLog("SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 260:
				msg := NewRecSipTCP(data)
				NewSipHEP(msg)

				if *gaddr != "" {
					msg.SendLog("SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 261:
				msg := NewSendSipTCP(data)
				NewSipHEP(msg)

				if *gaddr != "" {
					msg.SendLog("SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 268:
				msg := NewQosStats(data)
				/*
					NewQosHEPincRTP(msg)
					NewQosHEPincRTCP(msg)
					NewQosHEPoutRTP(msg)
					NewQosHEPoutRTCP(msg)
				*/

				if *saddr != "" {
					if msg.Data.QOS.IncMos >= 100 && msg.Data.QOS.OutMos >= 100 {
						msg.SendStatsd("QOS")
					}
				}
				if *gaddr != "" {
					msg.SendLog("QOS")
				}
			}

		} else {
			checkErr(err)
		}

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
