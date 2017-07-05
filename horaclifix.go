package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	_ "net/http/pprof"
	"os"
)

var (
	addr  = flag.String("l", ":4739", "Host ipfix listen address")
	haddr = flag.String("H", "127.0.0.1:9060", "Homer server address")
	debug = flag.Bool("d", false, "Debug output to stdout")
	gaddr = flag.String("g", "", "Graylog server address")
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

// Start handles incoming packets
func Start(conn *net.TCPConn) {
	fmt.Println("Handling new connection...")

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
		setLen := int(set.SetHeader.Length)
		setID := int(set.SetHeader.ID)

		if setID == 256 && version == 10 && dataLen > 20 {
			SendHandshake(conn, packet)
		}

		for len(packet) > 200 && dataLen-setLen == 16 {
			version = int(uint16(packet[0])<<8 + uint16(packet[1]))
			dataLen = int(uint16(packet[2])<<8 + uint16(packet[3]))
			setID = int(uint16(packet[16])<<8 + uint16(packet[17]))
			setLen = int(uint16(packet[18])<<8 + uint16(packet[19]))

			if *debug {
				fmt.Println("########################################################################################################################################")
				fmt.Printf("Inc: len(packet): %d, datalen: %d, setID: %d, version: %d\n", len(packet), dataLen, setID, version)
			}

			if setID > 280 || setID < 258 || version != 10 {
				break
			}

			if len(packet) < dataLen {
				if *debug {
					fmt.Printf("Out of sync: len(packet): %d, datalen: %d, setID: %d, version: %d\n", len(packet), dataLen, setID, version)
				}
				dataLen = len(packet)
			}

			// Create a new data packet with the header length. This is our first dataset
			data := packet[:dataLen]
			// Cut the first dataset from the original packet
			packet = packet[dataLen:]

			/*	version = int(uint16(data[0])<<8 + uint16(data[1]))
				dataLen = int(uint16(data[2])<<8 + uint16(data[3]))
				setID = int(uint16(data[16])<<8 + uint16(data[17]))
				setLen = int(uint16(data[18])<<8 + uint16(data[19]))
			*/

			if *debug {
				fmt.Printf("Out: len(packet): %d\n\n", len(packet))
				fmt.Println("Hexdump output:")
				fmt.Printf("%s\n", hex.Dump(data))
			}

			// Go through the set's and fill the right structs
			switch setID {
			case 258:
				msg := NewRecSipUDP(data)
				SendSipHEP(msg)

				if *gaddr != "" {
					LogSIP(msg)
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}

			case 259:
				msg := NewSendSipUDP(data)
				SendSipHEP(msg)

				if *gaddr != "" {
					LogSIP(msg)
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 260:
				msg := NewRecSipTCP(data)
				SendSipHEP(msg)

				if *gaddr != "" {
					LogSIP(msg)
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 261:
				msg := NewSendSipTCP(data)
				SendSipHEP(msg)

				if *gaddr != "" {
					LogSIP(msg)
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 268:
				msg := NewQosStats(data)
				SendQosHEP(msg)
				if *gaddr != "" {
					LogQOS(msg)
				}
			}
		}

	}
}

func main() {

	flag.Parse()
	fmt.Printf("IPFIX IP %v\nHomer IP %v\nGraylog IP %v\n", *addr, *haddr, *gaddr)

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
		go Start(conn)
	}

}
