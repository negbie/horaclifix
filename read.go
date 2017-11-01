package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

const headerLen int = 20

// Read handles the incoming packets
func Read(conn *net.TCPConn) {
	var (
		version int
		dataLen int
		setID   int
		dataSet []byte
	)

	log.Printf("New IPFIX connection from %s at %v\n", *name, conn.RemoteAddr())
	c := NewExtConns()

	// Close connections when this function ends
	defer func() {
		CloseExtConns(c)
		log.Printf("Close IPFIX connection to %s at %v\n", *name, conn.RemoteAddr())
		err := conn.Close()
		checkErr(err)
	}()

	// Create a buffer for incoming packets
	r := bufio.NewReader(conn)
	header := make([]byte, headerLen)

	for {
		// Read the next 20 bytes which represent
		// the ipfix header + setheader at once
		if _, err := io.ReadFull(r, header); err == nil {
			// Extract the ipfix version
			version = int(uint16(header[0])<<8 + uint16(header[1]))
			// Extract the length of the whole packet
			dataLen = int(uint16(header[2])<<8 + uint16(header[3]))
			// Extract the setID
			setID = int(uint16(header[16])<<8 + uint16(header[17]))
			// Check if we have valid setID's and IPFIX packets
			if setID > 280 || setID < 256 || version != 10 {
				log.Printf("[WARN] Malformed IPFIX header: %v\n", header)
				break
			}
			// Create a buffer which holds exactly one
			// dataSet of the length of (dataLen-headerLen)
			// or in words the packet length minus the header length
			dataSet = make([]byte, dataLen-headerLen)

		} else if err != nil {
			if err != io.EOF {
				log.Printf("[WARN] Couldn't read header: %s\n", err)
			}
			break
		}
		// Read the next (dataLen-headerLen) bytes into
		// the dataSet buffer at once
		if _, err := io.ReadFull(r, dataSet); err == nil {

			if *verbose {
				fmt.Println("########################################################################")
				fmt.Printf("Headerversion: %d, Headerlength: %d, SetID: %d\n", version, dataLen, setID)
				fmt.Println("Header in raw:", header)
			}
			if *debug {
				fmt.Println("Hexdump output:")
				fmt.Printf("%s\n", hex.Dump(header))
				fmt.Printf("%s\n", hex.Dump(dataSet))
			}

			switch setID {
			case 256:
				// Merge the header with the dataSet.
				data := append(header, dataSet...)
				SendHandshake(conn, data)
			case 258:
				msg := ParseRecSipUDP(dataSet)
				c.Send(msg, "SIP")
			case 259:
				msg := ParseSendSipUDP(dataSet)
				c.Send(msg, "SIP")
			case 260:
				msg := ParseRecSipTCP(dataSet)
				c.Send(msg, "SIP")
			case 261:
				msg := ParseSendSipTCP(dataSet)
				c.Send(msg, "SIP")
			case 268:
				msg := ParseQosStats(dataSet)
				c.Send(msg, "QOS")
			default:
				log.Printf("[WARN] Unhandled SetID %v\n", setID)
			}
		} else if err != nil {
			if err != io.EOF {
				log.Printf("[WARN] Couldn't read dataset: %s\n", err)
			}
			break
		}
	}
}
