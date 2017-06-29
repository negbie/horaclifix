package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	addr        = flag.String("l", ":4739", "Host ipfix listen address")
	haddr       = flag.String("H", "127.0.0.1:9060", "Homer server address")
	debug       = flag.Bool("d", false, "Debug output to stdout")
	graylogAddr = flag.String("g", "127.0.0.1:4488", "Graylog server address")
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// Start handles incoming packets
func Start(conn *net.TCPConn, haddr string, debug bool) {
	fmt.Println("Handling new connection...")
	NewGelfLogger()
	// UDP connection to Homer
	hconn, _ := net.Dial("udp", haddr)

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	byts := make([]byte, 32768)
	r := bufio.NewReader(conn)

	for {
		blen, err := r.Read(byts)
		buf := new(bytes.Buffer)
		buf.Write(byts[:blen])
		// Create a new buffer with the actual packet
		packet := buf.Bytes()

		// Check for EOF and go out of this loop. Don't cut the connection. Mby we just rebooted the sbc
		if err == io.EOF {
			fmt.Printf("EOF %v\n", err)
			break
		}

		set := NewHeader(packet)
		version := int(set.Header.Version)
		dataLen := int(set.Header.Length)
		setID := int(set.SetHeader.ID)
		if setID == 256 {
			SendHandshake(conn, packet)
		}

		for len(packet) > 200 {
			version = int(uint16(packet[0])<<8 + uint16(packet[1]))
			dataLen = int(uint16(packet[2])<<8 + uint16(packet[3]))
			setID = int(uint16(packet[16])<<8 + uint16(packet[17]))

			if setID > 280 || setID < 258 || version != 10 {
				break
			}

			if len(packet) < dataLen {
				fmt.Println("If this happen we are out of sync!")
				dataLen = len(packet)
			}

			// Create a new packet with the header length. This is our first dataset
			//fmt.Printf("%s\n", hex.Dump(packet))
			data := packet[:dataLen]
			if debug {
				fmt.Printf("%s\n", hex.Dump(data))
			}

			// Cut the first dataset from the original packet
			packet = packet[dataLen:]

			//fmt.Printf("%s\n", hex.Dump(packet))

			// Go through the set's and fill the right structs
			switch setID {
			case 258:
				dataSet := NewRecSipUDP(data)
				if debug {
					fmt.Printf("%s\n", dataSet.Data.SIP.SipMsg)
				}
				SendHEP(dataSet, hconn)
				LogSip(dataSet)
			case 259:
				dataSet := NewSendSipUDP(data)
				if debug {
					fmt.Printf("%s\n", dataSet.Data.SIP.SipMsg)
				}
				SendHEP(dataSet, hconn)
				LogSip(dataSet)
			case 260:
				dataSet := NewRecSipTCP(data)
				if debug {
					fmt.Printf("%s\n", dataSet.Data.SIP.SipMsg)
				}
				SendHEP(dataSet, hconn)
			case 261:
				dataSet := NewSendSipTCP(data)
				if debug {
					fmt.Printf("%s\n", dataSet.Data.SIP.SipMsg)
				}
				SendHEP(dataSet, hconn)
			case 262:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			case 263:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			case 264:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			case 265:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			case 266:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			case 267:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 268:
				// GOTCHA!!!!
				dataSet := NewQosStats(data)
				if debug {
					fmt.Printf("%s\n", dataSet.Data.QOS.IncCallID)
				}
				LogQos(dataSet)

			case 269:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 271:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 272:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 273:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 274:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 275:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 276:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 277:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			case 278:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))

			case 279:
				// Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				fmt.Printf("Unkown battlefield!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n\n%s\n", hex.Dump(data))
			}

		}

	}
}

/*
func handleConn(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		Start(conn, *haddr, *debug)
		out <- conn
	}
}
func closeConn(in <-chan *net.TCPConn) {
	for conn := range in {
		conn.Close()
	}
}
func main() {
	flag.Parse()
	fmt.Printf("Listening for IPFIX at: %v\n Send to Homer at: %v\n\n", *addr, *haddr)
	addr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)
	for i := 0; i < 5; i++ {
		go handleConn(pending, complete)
	}
	go closeConn(complete)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		pending <- conn
	}
}
*/

func main() {

	flag.Parse()
	fmt.Printf("Listening for IPFIX at: %v\nSend to Homer at: %v\n\n", *addr, *haddr)

	laddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		os.Exit(1)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			os.Exit(1)
		}
		go Start(conn, *haddr, *debug)
	}
}
