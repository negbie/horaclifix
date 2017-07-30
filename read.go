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
	log.Printf("New IPFIX connection from %s at %v\n", *name, conn.RemoteAddr())
	c := NewConnections()

	// Close connection when this function ends
	defer func() {
		if *aaddr != "" {
			c.Amqp.Close()
			c.AmqpChannel.Close()
		}
		if *baddr != "" {
			c.Banshee.Close()
		}
		if *gaddr != "" {
			c.Graylog.Close()
		}
		if *gtaddr != "" {
			c.GraylogTLS.Close()
		}
		if *haddr != "" {
			c.Homer.Close()
		}
		if *saddr != "" {
			c.StatsD.Close()
		}
		log.Printf("Close IPFIX connection to %s at %v\n", *name, conn.RemoteAddr())
		conn.Close()
	}()

	// Create a buffer for incoming packets
	r := bufio.NewReader(conn)
	header := make([]byte, headerLen)
	var (
		version int
		dataLen int
		setID   int
		dataSet []byte
	)

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
				log.Printf("Malformed IPFIX header: %v", header)
				break
			}
			// Create a buffer which holds exactly one
			// dataSet of the length of (dataLen-headerLen)
			// or in words the packet length minus the header length
			dataSet = make([]byte, dataLen-headerLen)

		} else if err != nil {
			if err != io.EOF {
				log.Printf("Header read error: %s", err)
			}
			break
		}
		// Read the next (dataLen-headerLen) bytes into
		// the dataSet buffer at once
		if _, err := io.ReadFull(r, dataSet); err == nil {
			// Now merge the header with the dataSet.
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
				msg := ParseRecSipUDP(data)
				c.Send(msg, "SIP")
			case 259:
				msg := ParseSendSipUDP(data)
				c.Send(msg, "SIP")
			case 260:
				msg := ParseRecSipTCP(data)
				c.Send(msg, "SIP")
			case 261:
				msg := ParseSendSipTCP(data)
				c.Send(msg, "SIP")
			case 268:
				msg := ParseQosStats(data)
				c.Send(msg, "QOS")
			default:
				log.Printf("Unhandled setID %v\n", setID)
			}
		} else if err != nil {
			if err != io.EOF {
				log.Printf("Dataset read error: %s", err)
			}
			break
		}
	}
}
