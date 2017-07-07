package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"net"
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
		log.Println("Warning:", err)
	}
}

func checkCritErr(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}

// Start handles incoming packets
func Start(conn *net.TCPConn) {
	log.Printf("Handling new connection under %v\n", *addr)

	// Close connection when this function ends
	defer func() {
		log.Printf("Closing connection under %v\n", *addr)
		conn.Close()
	}()

	byts := make([]byte, 32768)
	r := bufio.NewReader(conn)

	for {
		blen, err := r.Read(byts)
		// Check for EOF and go out of this loop. Don't cut the connection. Mby we just rebooted the sbc
		if err == io.EOF {
			break
		}
		buf := new(bytes.Buffer)
		_, err = buf.Write(byts[:blen])
		checkErr(err)
		// Create a new buffer with the actual packet
		packet := buf.Bytes()

		if len(packet) > 20 {

			version := int(uint16(packet[0])<<8 + uint16(packet[1]))
			dataLen := int(uint16(packet[2])<<8 + uint16(packet[3]))
			setID := int(uint16(packet[16])<<8 + uint16(packet[17]))
			setLen := int(uint16(packet[18])<<8 + uint16(packet[19]))

			if setID == 256 && version == 10 && dataLen > 20 {
				SendHandshake(conn, packet)
			}

			for len(packet) > 200 && dataLen-setLen == 16 && version == 10 {
				version = int(uint16(packet[0])<<8 + uint16(packet[1]))
				dataLen = int(uint16(packet[2])<<8 + uint16(packet[3]))
				setID = int(uint16(packet[16])<<8 + uint16(packet[17]))
				setLen = int(uint16(packet[18])<<8 + uint16(packet[19]))

				if *debug {
					log.Println("########################################################################################################################################")
					log.Printf("Inc: len(packet): %d, datalen: %d, setID: %d, version: %d\n", len(packet), dataLen, setID, version)
				}

				if setID > 280 || setID < 258 || version != 10 {
					break
				}

				if len(packet) < dataLen {
					if *debug {
						log.Printf("Out of sync: len(packet): %d, datalen: %d, setID: %d, version: %d\n", len(packet), dataLen, setID, version)
					}
					dataLen = len(packet)
				}

				// Create a new data packet with the header length. This is our first dataset
				data := packet[:dataLen]
				// Cut the first dataset from the original packet
				packet = packet[dataLen:]

				if *debug {
					log.Printf("Out: len(packet): %d\n\n", len(packet))
					log.Println("Hexdump output:")
					log.Printf("%s\n", hex.Dump(data))
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
						log.Println("SIP output:")
						log.Printf("%s\n", msg.Data.SIP.SipMsg)
					}

				case 259:
					msg := NewSendSipUDP(data)
					SendSipHEP(msg)

					if *gaddr != "" {
						LogSIP(msg)
					}
					if *debug {
						log.Println("SIP output:")
						log.Printf("%s\n", msg.Data.SIP.SipMsg)
					}
				case 260:
					msg := NewRecSipTCP(data)
					SendSipHEP(msg)

					if *gaddr != "" {
						LogSIP(msg)
					}
					if *debug {
						log.Println("SIP output:")
						log.Printf("%s\n", msg.Data.SIP.SipMsg)
					}
				case 261:
					msg := NewSendSipTCP(data)
					SendSipHEP(msg)

					if *gaddr != "" {
						LogSIP(msg)
					}
					if *debug {
						log.Println("SIP output:")
						log.Printf("%s\n", msg.Data.SIP.SipMsg)
					}
				case 268:
					msg := NewQosStats(data)
					/*
						SendQosHEPincRTP(msg)
						SendQosHEPincRTCP(msg)
						SendQosHEPoutRTP(msg)
						SendQosHEPoutRTCP(msg)
					*/
					if *gaddr != "" {
						LogQOS(msg)
					}
				}
			}
		}

	}
}

func main() {
	flag.Parse()
	var err error
	var f *os.File
	if !*debug {
		f, err = os.OpenFile("horaclifix.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		f, err = os.OpenFile("horaclifix.debug", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
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
		go Start(conn)
	}

}
