package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

func toString(bs []byte) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

// BufToInt8 casts []byte to []int8
func BufToInt8(b []byte) []int8 {
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}
	return bi
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

	byts := make([]byte, 65535)
	r := bufio.NewReader(conn)

	for {
		// Create a bytes holder and read in the bytes from the network
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

		// Create a new header struct to get the header length & ID
		set := NewHeader(packet)
		dataLen := int(set.Header.Length)
		setLen := int(set.SetHeader.Length)
		setID := int(set.SetHeader.ID)

		for len(packet) > 20 && len(packet) >= dataLen && setID > 255 && setID < 280 && dataLen-setLen == 16 {

			// Get the header length from the packet at position 2&3
			dataLen = int(uint16(packet[2])<<8 + uint16(packet[3]))
			setLen = int(uint16(packet[18])<<8 + uint16(packet[19]))

			if debug {
				fmt.Println("####################################################################")
				fmt.Printf("Length of incoming packet: %d\n", len(packet))
				fmt.Printf("Length of following Data: %d\n", dataLen)
			}

			if len(packet) < dataLen {
				fmt.Println("If this happen we are out of sync!")
				dataLen = len(packet)
			}
			// Create a new packet with the header length. This is our first dataset
			data := packet[:dataLen]
			// Cut the first dataset from the original packet
			packet = packet[dataLen:]
			setID = int(uint16(data[16])<<8 + uint16(data[17]))
			if debug {
				fmt.Printf("Length of incoming packet: %d\n", len(data))
				fmt.Printf("Length from header: %d\n", dataLen)
				fmt.Printf("SetID: %d\n\n", setID)
				fmt.Printf("%s\n", hex.Dump(data))
			}
			// Go through the set's and fill the right structs

			switch setID {

			case 0:
				// Timeout packets
			case 256:
				h := NewHandShake(data)
				h.SetHeader.ID++
				// Disable timeout
				h.Data.HandShake.Timeout = 0
				binary.Write(conn, binary.BigEndian, BufToInt8(h.SendHandShake()))
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
		receiver.Start(conn, *haddr, *debug)
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
