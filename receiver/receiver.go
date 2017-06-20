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

var (
	hep3ID           = []byte{0x48, 0x45, 0x50, 0x33}
	iPProtocolFamily = []byte{0x02}
	protocolType     = []byte{0x01}
	captureAgentID   = []byte{0x07, 0xD1} //2001
)

type IpfixHeader struct {
	Version       uint16
	Length        uint16
	ExportTime    uint32
	SeqNum        uint32
	ObservationID uint32
}

// SetHeader represents the setheader fields
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
	HandShake  Hs
	RecSipUDP  Rsu
	SendSipUDP Ssu
	RecSipTCP  Rst
	SendSipTCP Sst
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

// IPFIX messages struct
type IPFIX struct {
	Header    IpfixHeader
	SetHeader IpfixSetHeader
	Data      DataSet
}

// Rsu holds the RecSipUDP dataset fields
type Rsu struct {
	TimeSec   uint32
	TimeMic   uint32
	IntSlot   uint8
	IntPort   uint8
	IntVlan   uint16
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
	UDPlen    uint16
	MsgLen    uint16
	SipMsg    []byte
}

// Ssu holds the SendSipUDP dataset fields
type Ssu struct {
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
	UDPlen    uint16
	MsgLen    uint16
	SipMsg    []byte
}

// Rst holds the RecSipTCP dataset fields
type Rst struct {
	TimeSec   uint32
	TimeMic   uint32
	IntSlot   uint8
	IntPort   uint8
	IntVlan   uint16
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
	UDPlen    uint16
	MsgLen    uint16
	SipMsg    []byte
}

// Sst holds the SendSipTCP dataset fields
type Sst struct {
	TimeSec   uint32
	TimeMic   uint32
	IntSlot   uint8
	IntPort   uint8
	IntVlan   uint16
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
	UDPlen    uint16
	MsgLen    uint16
	SipMsg    []byte
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

func NewRecSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.IPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.VL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TOS)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TFlags)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TTL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TProto)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.TPos)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.UDPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.MsgLen)

	/*	fmt.Println(r.Len())
		fmt.Println("headerlen:", len(header))
	*/
	ipfix.Data.RecSipUDP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.RecSipUDP.SipMsg)

	return &ipfix
}

func NewSendSipUDP(header []byte) *IPFIX {
	var ipfix IPFIX
	r := bytes.NewReader(header)

	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TimeSec)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TimeMic)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.IntSlot)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.IntPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.IntVlan)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.CallIDLen)
	ipfix.Data.SendSipUDP.CallID = make([]byte, ipfix.Data.SendSipUDP.CallIDLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.CallID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.CallIDEnd)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.IPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.VL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TOS)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TLen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TID)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TFlags)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TTL)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TProto)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.TPos)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.SrcIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.DstIP)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.DstPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.SrcPort)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.UDPlen)
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.MsgLen)

	/*	fmt.Println(r.Len())
		fmt.Println("headerlen:", len(header))*/

	ipfix.Data.SendSipUDP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SendSipUDP.SipMsg)

	return &ipfix
}

// HandShake writes the binary IPFIX representation into the buffer
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

func (ipfix *IPFIX) SendHep3() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, hep3ID)
	binary.Write(b, binary.BigEndian, len(hep3ID)+len(iPProtocolFamily))
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.TProto)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.SrcIP)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.DstIP)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.SrcPort)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.DstPort)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.TimeSec)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.TimeMic)
	binary.Write(b, binary.BigEndian, protocolType)
	binary.Write(b, binary.BigEndian, captureAgentID)
	binary.Write(b, binary.BigEndian, &ipfix.Data.SendSipUDP.SipMsg)

	return b.Bytes()
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

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		c.Close()
	}()

	//c.SetReadDeadline(time.Now())
	n, err := c.Read(hs)
	if err == io.EOF {
		fmt.Printf("EOF %v", err)

	} /*else {
		//var zero time.Time
		c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	}*/

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
		b := make([]byte, 4096)
		n, err := c.Read(b)
		if err == io.EOF {
			fmt.Printf("EOF %v", err)
			break

		}

		set := NewHeader(b[:n]).SetHeader.ID
		fmt.Println(set)

		switch set {
		case 258:
			packet := NewRecSipUDP(b[:n])
			fmt.Printf("%s", packet.Data.RecSipUDP)
		case 259:
			packet := NewSendSipUDP(b[:n])
			fmt.Printf("%s", packet.Data.SendSipUDP)

		}

		/*		// Read IPFIX Messages delimited by newline
				bytes, err := bufReader.ReadBytes('\n')
				checkError(err)*/

		//fmt.Printf("%s", (binary.BigEndian.Uint16(bytes)))

	}
}
