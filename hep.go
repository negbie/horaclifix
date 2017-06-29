package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

// NewHEPMsg writes the binary HEP representation into the buffer
func NewHEPMsg(msg []byte) []byte {

	b := bytes.NewBuffer(make([]byte, 6))
	binary.Write(b, binary.BigEndian, msg)
	packet := b.Bytes()
	binary.BigEndian.PutUint32(packet, uint32(0x48455033)) // ASCII "HEP3"
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))
	return packet
}

// NewHEPChunck constructs the HEP chunck
func (ipfix *IPFIX) NewHEPChunck(ChunckVen uint16, ChunckType uint16) []byte {

	b := bytes.NewBuffer(make([]byte, 6))
	switch ChunckType {
	case 0x0001:
		binary.Write(b, binary.BigEndian, 0x02)

	case 0x0002:
		binary.Write(b, binary.BigEndian, 0x11)

	case 0x0003:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SrcIP)

	case 0x0004:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.DstIP)

	case 0x0007:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SrcPort)

	case 0x0008:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.DstPort)

	case 0x0009:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.TimeSec)

	case 0x000a:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.TimeMic)

	case 0x000b:
		binary.Write(b, binary.BigEndian, 0x01)

	case 0x000c:
		binary.Write(b, binary.BigEndian, 0x000007D1)

	case 0x000f:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	case 0x0020:
		// MOS
	}

	packet := b.Bytes()
	binary.BigEndian.PutUint16(packet, ChunckVen)
	binary.BigEndian.PutUint16(packet[2:4], ChunckType)
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))

	return packet
}

// SendHEP sends the HEP message
func SendHEP(i *IPFIX, c net.Conn) {
	bhep := new(bytes.Buffer)
	bhep.Write(i.NewHEPChunck(0x0000, 0x0001))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0002))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0003))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0004))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0007))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0008))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0009))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000a))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000b))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000c))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000f))

	//fmt.Printf("%s\n", hex.Dump(bhep.Bytes()))
	c.Write(NewHEPMsg(bhep.Bytes()))
}
