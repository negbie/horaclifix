package receiver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

// Sync message from sbc
// echo -ne '\x00\x0A\x00\x30\x59\x41\x37\x38\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x20\x00\x01\x00\x02\x00\xFC\x77\x31\x00\x00\x00\x1E\x00\x00\x00\x00\x43\x5A\x07\x03\x00\x06\x65\x63\x7A\x37\x33\x30' | nc localhost 4739

// Iheader is the Oracle IPFIX Header
//
// Wire format:
//
// Bytes:  0                   1                   2                   3
// Bits:   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |             Version (uint16)          |            Length (uint16)            |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                   ExportTime (unit32)                         |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                 SequenceNumber (uint32)                       |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                 ObservationID (uint32)                        |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
//

// IPFIX messages struct
type IPFIX struct {
	Header    IpfixHeader
	SetHeader IpfixSetHeader
	Data      DataSet
}

// IpfixHeader holds the ipfix header fields
type IpfixHeader struct {
	Version       uint16
	Length        uint16
	ExportTime    uint32
	SeqNum        uint32
	ObservationID uint32
}

// IpfixSetHeader represents the setheader fields
type IpfixSetHeader struct {
	ID     uint16
	Length uint16
}

// DataSet holds multiple datasets with following SetID's:
// HandShake: 	257
// RecSipUDP: 	258
// SendSipUDP: 	259
// RecSipTCP: 	260
// SendSipTCP: 	261
type DataSet struct {
	HandShake Hs
	SIP       SipSet
}

// Hs holds the HandShake dataset fields
type Hs struct {
	MaVer       uint16
	MiVer       uint16
	CFlags1     uint16
	CFlags2     uint16
	SFlags      uint16
	Timeout     uint16
	SystemID    uint32
	Product     uint16
	SMaVer      uint8
	SMiVer      uint8
	Revision    uint8
	HostnameLen uint8
	Hostname    []byte
}

// SipUDP holds the SendSipUDP dataset fields
type SipSet struct {
	TimeSec   uint32
	TimeMic   uint32
	IntSlot   uint8
	IntPort   uint8
	IntVlan   uint16
	CallIDLen uint8
	CallID    []byte
	CallIDEnd uint8
	IPlen     uint16
	VL        uint8
	TOS       uint8
	TLen      uint16
	TID       uint16
	TFlags    uint16
	TTL       uint8
	TProto    uint8
	TPos      uint16
	SrcIP     uint32
	DstIP     uint32
	DstPort   uint16
	SrcPort   uint16
	Context   uint32
	UDPlen    uint16
	MsgLen    uint16
	SipMsg    []byte
}

// SendHandShake writes the binary Handshake representation into the buffer
func SendHEP(msg []byte) []byte {

	b := bytes.NewBuffer(make([]byte, 6))
	binary.Write(b, binary.BigEndian, msg)
	packet := b.Bytes()
	binary.BigEndian.PutUint32(packet, uint32(0x48455033)) // ASCII "HEP3"
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))

	return packet
}

// SendHandShake writes the binary Handshake representation into the buffer
func (ipfix *IPFIX) NewHEP(ChunckVen uint16, ChunckType uint16) []byte {

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
		binary.Write(b, binary.BigEndian, 0x000000E4)

	case 0x000f:
		binary.Write(b, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	}
	packet := b.Bytes()
	binary.BigEndian.PutUint16(packet, ChunckVen)
	binary.BigEndian.PutUint16(packet[2:4], ChunckType)
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))

	return packet
}

// BufToInt8 casts []byte to []int8
func BufToInt8(b []byte) []int8 {
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}
	return bi
}

// NewHeader fills the IpfixHeader struct with structured binary data from r
func NewHeader(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	return &ipfix
}

// NewHandShake fills the IPFIX struct with structured binary data from r
func NewHandShake(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.MaVer)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.MiVer)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.CFlags1)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.CFlags2)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.SFlags)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.Timeout)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.SystemID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.Product)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.SMaVer)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.SMiVer)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.Revision)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.HostnameLen)
	ipfix.Data.HandShake.Hostname = make([]byte, ipfix.Data.HandShake.HostnameLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake.Hostname)

	fmt.Printf("Header: %#v", ipfix)
	return &ipfix
}

// SendHandShake writes the binary Handshake representation into the buffer
func (ipfix *IPFIX) SendHandShake() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, &ipfix.Header)
	binary.Write(b, binary.BigEndian, &ipfix.SetHeader)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.MaVer)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.MiVer)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.CFlags1)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.CFlags2)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.SFlags)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.Timeout)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.SystemID)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.Product)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.SMaVer)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.SMiVer)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.Revision)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.HostnameLen)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake.Hostname)

	return b.Bytes()
}

func NewRecSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.VL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TOS)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TFlags)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TTL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TProto)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TPos)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.UDPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)

	/*	fmt.Println(r.Len())
		fmt.Println("headerlen:", len(header))
	*/
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewSendSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDLen)
	ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.VL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TOS)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TFlags)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TTL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TProto)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TPos)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.UDPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)

	/*	fmt.Println(r.Len())
		fmt.Println("headerlen:", len(header))*/

	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewRecSipTCP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.Context)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)

	/*	fmt.Println(r.Len())
		fmt.Println("headerlen:", len(header))
	*/
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewSendSipTCP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.Context)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDLen)
	ipfix.Data.SIP.CallID = make([]byte, ipfix.Data.SIP.CallIDLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.MsgLen)

	/*	fmt.Println(r.Len())
		fmt.Println("headerlen:", len(header))*/

	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func SyncClient(c net.Conn) {
	fmt.Println("Handling new connection...")

	hs := make([]byte, 512)

	con, _ := net.Dial("udp", "127.0.0.1:9060")

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		c.Close()
	}()

	n, err := c.Read(hs)
	if err == io.EOF {
		fmt.Printf("EOF %v", err)

	}

	s := NewHeader(hs[:n]).SetHeader.ID
	hlen := NewHeader(hs[:n]).SetHeader.Length

	if s == 256 && hlen > 20 {
		h := NewHandShake(hs[:n])
		h.SetHeader.ID++
		// fmt.Printf("Write: %#v\n", BufToInt8(h.SendHandShake()))
		binary.Write(c, binary.BigEndian, BufToInt8(h.SendHandShake()))
	}
	//bufReader := bufio.NewReader(conn)
	fmt.Println(s)
	for {
		/*		buf := bufio.NewReader(c)
				//Read IPFIX Messages delimited by newline
				_, err := buf.ReadBytes('\n')
				bi := new(bytes.Buffer)
				bi.ReadFrom(buf)
				b := bi.Bytes()
				n := len(b)
		*/

		b := make([]byte, 4096)
		n, err := c.Read(b)

		if err == io.EOF {
			fmt.Printf("EOF %v", err)
			break
		}

		set := NewHeader(b[:n]).SetHeader

		fmt.Println(set.ID)

		switch set.ID {
		case 258:
			packet := NewRecSipUDP(b[:n])
			fmt.Printf("%#v\n\n", packet.Data)
			fmt.Printf("%s\n\n", packet.Data.SIP.SipMsg)
		case 259:
			if n-int(set.Length) > 16 {
				fmt.Println("SetHeader<<<< Packetlen!!!!!!!!!!!")

				packet := NewSendSipUDP(b[:int(set.Length)+15])
				fmt.Printf("%#v\n\n", packet.Data)
				fmt.Printf("%s\n\n", packet.Data.SIP.SipMsg)

				packet = NewSendSipUDP(b[int(set.Length)+16 : n])
				fmt.Printf("%#v\n\n", packet.Data)
				fmt.Printf("%s\n\n", packet.Data.SIP.SipMsg)

			} else {
				packet := NewSendSipUDP(b[:n])
				fmt.Printf("%s\n", packet.Data.SIP.SipMsg)
				bhep := new(bytes.Buffer)
				bhep.Write(packet.NewHEP(0x0000, 0x0001))
				bhep.Write(packet.NewHEP(0x0000, 0x0002))
				bhep.Write(packet.NewHEP(0x0000, 0x0003))
				bhep.Write(packet.NewHEP(0x0000, 0x0004))
				bhep.Write(packet.NewHEP(0x0000, 0x0007))
				bhep.Write(packet.NewHEP(0x0000, 0x0008))
				bhep.Write(packet.NewHEP(0x0000, 0x0009))
				bhep.Write(packet.NewHEP(0x0000, 0x000a))
				bhep.Write(packet.NewHEP(0x0000, 0x000b))
				bhep.Write(packet.NewHEP(0x0000, 0x000c))
				bhep.Write(packet.NewHEP(0x0000, 0x000f))
				//fmt.Printf("HEEEP: %v\n", SendHEP(bhep.Bytes()))
				con.Write(SendHEP(bhep.Bytes()))

			}
		case 260:
			packet := NewRecSipTCP(b[:n])
			fmt.Printf("%#v\n\n", packet.Data)
			fmt.Printf("%s\n\n", packet.Data.SIP.SipMsg)
		case 261:
			packet := NewSendSipTCP(b[:n])
			fmt.Printf("%#v\n\n", packet.Data)
			fmt.Printf("%s\n\n", packet.Data.SIP.SipMsg)

		}
	}
}
