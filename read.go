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
func Read(ic *net.TCPConn) {
	var (
		version int
		dataLen int
		setID   int
		dataSet []byte
	)

	log.Printf("New IPFIX connection from %s at %v\n", *name, ic.RemoteAddr())
	conn := NewExtConns()

	// Close connections when this function ends
	defer func() {
		CloseExtConns(conn)
		log.Printf("Close IPFIX connection to %s at %v\n", *name, ic.RemoteAddr())
		err := ic.Close()
		checkErr(err)
	}()

	// Create a buffer for incoming packets
	r := bufio.NewReader(ic)
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
			// Check if we have a valid IPFIX packet
			if version != 10 || dataLen < headerLen {
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
				fmt.Printf("headerversion: %d, headerlength: %d, setID: %d\n", version, dataLen, setID)
				fmt.Println("raw header:", header)
			}
			if *debug {
				fmt.Println("hexdump output:")
				fmt.Printf("%s\n", hex.Dump(header))
				fmt.Printf("%s\n", hex.Dump(dataSet))
			}
			if *paddr != "" {
				prom.CounterVecMetrics["horaclifix_packets_total"].WithLabelValues(*name).Inc()
			}

			switch setID {
			case 256:
				// Append the dataSet to the header.
				data := append(header, dataSet...)
				SendHandshake(ic, data)
			case 258:
				msg := ParseRecSipUDP(dataSet)
				conn.Send(msg)
			case 259:
				msg := ParseSendSipUDP(dataSet)
				conn.Send(msg)
			case 260:
				msg := ParseRecSipTCP(dataSet)
				conn.Send(msg)
			case 261:
				msg := ParseSendSipTCP(dataSet)
				conn.Send(msg)
			case 268:
				msg := ParseQosStats(dataSet)
				conn.Send(msg)
			default:
				log.Printf("[WARN] Unhandled setID: %v\nfor dataset: %v\n", setID, dataSet)
			}
		} else if err != nil {
			if err != io.EOF {
				log.Printf("[WARN] Couldn't read dataset: %s\n", err)
			}
			break
		}
	}
}
