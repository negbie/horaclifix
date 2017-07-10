package main

import (
	"encoding/hex"
	"fmt"
	"net"
)

func parse(conn *net.TCPConn, packet []byte) {

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
			fmt.Println("################################################################################################")
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

		if *debug {
			fmt.Printf("Out: len(packet): %d\n\n", len(packet))
			fmt.Println("Hexdump output:")
			fmt.Printf("%s\n", hex.Dump(data))
		}

		// Go through the set's and fill the right structs
		switch setID {
		case 258:
			msg := NewRecSipUDP(data)
			NewSipHEP(msg)

			if *gaddr != "" {
				LogSIP(msg)
			}
			if *debug {
				fmt.Println("SIP output:")
				fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
			}

		case 259:
			msg := NewSendSipUDP(data)
			NewSipHEP(msg)

			if *gaddr != "" {
				LogSIP(msg)
			}
			if *debug {
				fmt.Println("SIP output:")
				fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
			}
		case 260:
			msg := NewRecSipTCP(data)
			NewSipHEP(msg)

			if *gaddr != "" {
				LogSIP(msg)
			}
			if *debug {
				fmt.Println("SIP output:")
				fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
			}
		case 261:
			msg := NewSendSipTCP(data)
			NewSipHEP(msg)

			if *gaddr != "" {
				LogSIP(msg)
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
			msg.SendStatsd("QOS")
			if *gaddr != "" {
				LogQOS(msg)
			}
		}
	}
}
