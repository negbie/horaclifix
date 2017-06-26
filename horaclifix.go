package main

import (
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
	addr  = flag.String("h", ":4739", "Host ipfix listen address")
	haddr = flag.String("H", "127.0.0.1:9060", "Homer server address")
	debug = flag.Bool("d", false, "Debug output to stdout")
)

// IPFIX holds the structure of one IPFIX packet
//
// Wire format:
//
// Bytes:  0                   1                   2                   3
// Bits:   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |             Version (uint16)          |            Length (uint16)            |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                 ExportTime     (unit32)                       |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                 SequenceNumber (uint32)                       |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                 ObservationID  (uint32)                       |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |           SetHeader ID (uint16)       |       SetHeader Length (uint16)       |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//         |                                 Dataset..................                     |
//         +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
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
	MaVer    uint16
	MiVer    uint16
	CFlags1  uint16
	CFlags2  uint16
	SFlags   uint16
	Timeout  uint16
	SystemID uint32
	Product  uint16
	SMaVer   uint8
	SMiVer   uint8
	Revision uint8
	//HostnameLen uint8
	//Hostname    []byte
}

// SipSet holds the SIP dataset fields
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

type CallStats struct {
	IncRtpBytes       uint32
	IncRtpPackets     uint32
	IncRtpLostPackets uint32
	IncRtpAvgJitter   uint32
	IncRtpMaxJitter   uint32

	IncRtcpBytes       uint32
	IncRtcpPackets     uint32
	IncRtcpLostPackets uint32
	IncRtcpAvgJitter   uint32
	IncRtcpMaxJitter   uint32
	IncRtcpAvgLat      uint32
	IncRtcpMaxLat      uint32

	IncrVal uint32
	IncMos  uint32

	OutRtpBytes       uint32
	OutRtpPackets     uint32
	OutRtpLostPackets uint32
	OutRtpAvgJitter   uint32
	OutRtpMaxJitter   uint32

	OutRtcpBytes       uint32
	OutRtcpPackets     uint32
	OutRtcpLostPackets uint32
	OutRtcpAvgJitter   uint32
	OutRtcpMaxJitter   uint32
	OutRtcpAvgLat      uint32
	OutRtcpMaxLat      uint32

	OutrVal uint32
	OutMos  uint32

	Type uint8

	CallerIncSrcIP   uint32
	CallerIncDstIP   uint32
	CallerIncSrcPort uint16
	CallerIncDstPort uint16

	CalleeIncSrcIP   uint32
	CalleeIncDstIP   uint32
	CalleeIncSrcPort uint16
	CalleeIncDstPort uint16

	CallerOutSrcIP   uint32
	CallerOutDstIP   uint32
	CallerOutSrcPort uint16
	CallerOutDstPort uint16

	CalleeOutSrcIP   uint32
	CalleeOutDstIP   uint32
	CalleeOutSrcPort uint16
	CalleeOutDstPort uint16

	CallerIntSlot uint8
	CallerIntPort uint8
	CallerIntVlan uint16

	CalleeIntSlot uint8
	CalleeIntPort uint8
	CalleeIntVlan uint16

	BeginTimeSec uint32
	BeginTimeMic uint32

	EndTimeSec   uint32
	EndinTimeMic uint32

	Seperator uint8

	IncRealmLen uint16
	IncRealm    []byte
	IncRealmEnd uint8

	OutRealmLen uint16
	OutRealm    []byte
	OutRealmEnd uint8

	IncCallIDLen uint16
	IncCallID    []byte
	IncCallIDEnd uint8

	OutCallIDLen uint16
	OutCallID    []byte
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
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake)
	return &ipfix
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
	ipfix.Data.SIP.SipMsg = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &ipfix.Data.SIP.SipMsg)

	return &ipfix
}

func NewCallQualityStats(header []byte) *CallStats {
	var ipfix IPFIX
	var callstats CallStats
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpBytes)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.IncRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &callstats.IncRtcpBytes)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &callstats.IncRtcpMaxLat)

	binary.Read(r, binary.BigEndian, &callstats.IncrVal)
	binary.Read(r, binary.BigEndian, &callstats.IncMos)

	binary.Read(r, binary.BigEndian, &callstats.OutRtpBytes)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.OutRtpMaxJitter)

	binary.Read(r, binary.BigEndian, &callstats.OutRtcpBytes)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpLostPackets)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpAvgJitter)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpMaxJitter)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpAvgLat)
	binary.Read(r, binary.BigEndian, &callstats.OutRtcpMaxLat)
	binary.Read(r, binary.BigEndian, &callstats.OutrVal)
	binary.Read(r, binary.BigEndian, &callstats.OutMos)

	binary.Read(r, binary.BigEndian, &callstats.Type)

	binary.Read(r, binary.BigEndian, &callstats.CallerIncSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerIncDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerIncSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CallerIncDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CalleeIncSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIncDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIncSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIncDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CallerOutSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerOutDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CallerOutSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CallerOutDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CalleeOutSrcIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeOutDstIP)
	binary.Read(r, binary.BigEndian, &callstats.CalleeOutSrcPort)
	binary.Read(r, binary.BigEndian, &callstats.CalleeOutDstPort)

	binary.Read(r, binary.BigEndian, &callstats.CallerIntSlot)
	binary.Read(r, binary.BigEndian, &callstats.CallerIntPort)
	binary.Read(r, binary.BigEndian, &callstats.CallerIntVlan)

	binary.Read(r, binary.BigEndian, &callstats.CalleeIntSlot)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIntPort)
	binary.Read(r, binary.BigEndian, &callstats.CalleeIntVlan)

	binary.Read(r, binary.BigEndian, &callstats.BeginTimeSec)
	binary.Read(r, binary.BigEndian, &callstats.BeginTimeMic)

	binary.Read(r, binary.BigEndian, &callstats.EndTimeSec)
	binary.Read(r, binary.BigEndian, &callstats.EndinTimeMic)

	binary.Read(r, binary.BigEndian, &callstats.Seperator)

	binary.Read(r, binary.BigEndian, &callstats.IncRealmLen)
	callstats.IncRealm = make([]byte, callstats.IncRealmLen)
	binary.Read(r, binary.BigEndian, &callstats.IncRealm)
	binary.Read(r, binary.BigEndian, &callstats.IncRealmEnd)

	binary.Read(r, binary.BigEndian, &callstats.OutRealmLen)
	callstats.OutRealm = make([]byte, callstats.OutRealmLen)
	binary.Read(r, binary.BigEndian, &callstats.OutRealm)
	binary.Read(r, binary.BigEndian, &callstats.OutRealmEnd)

	binary.Read(r, binary.BigEndian, &callstats.IncCallIDLen)
	callstats.IncCallID = make([]byte, callstats.IncCallIDLen)
	binary.Read(r, binary.BigEndian, &callstats.IncCallID)
	binary.Read(r, binary.BigEndian, &callstats.IncCallIDEnd)

	binary.Read(r, binary.BigEndian, &callstats.OutCallIDLen)
	callstats.OutCallID = make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, &callstats.OutCallID)

	return &callstats
}

// HandShake writes the binary Handshake representation into the buffer
func (ipfix *IPFIX) SendHandShake() []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, &ipfix.Header)
	binary.Write(b, binary.BigEndian, &ipfix.SetHeader)
	binary.Write(b, binary.BigEndian, &ipfix.Data.HandShake)
	return b.Bytes()
}

// BufToInt8 casts []byte to []int8
func BufToInt8(b []byte) []int8 {
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}
	return bi
}

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

	}
	packet := b.Bytes()
	binary.BigEndian.PutUint16(packet, ChunckVen)
	binary.BigEndian.PutUint16(packet[2:4], ChunckType)
	binary.BigEndian.PutUint16(packet[4:], uint16(len(packet)))

	return packet
}

// SendHEP sends the HEP message
func SendHEP(p *IPFIX, c net.Conn) {
	bhep := new(bytes.Buffer)
	bhep.Write(p.NewHEPChunck(0x0000, 0x0001))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0002))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0003))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0004))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0007))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0008))
	bhep.Write(p.NewHEPChunck(0x0000, 0x0009))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000a))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000b))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000c))
	bhep.Write(p.NewHEPChunck(0x0000, 0x000f))

	//fmt.Printf("%s\n", hex.Dump(bhep.Bytes()))
	c.Write(NewHEPMsg(bhep.Bytes()))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// Start handles incoming packets
func Start(conn *net.TCPConn, haddr string, debug bool) {
	fmt.Println("Handling new connection...")

	hconn, _ := net.Dial("udp", haddr)

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	for {
		// Create a bytes holder and read in the bytes from the network
		byts := make([]byte, 32768)
		blen, err := conn.Read(byts)

		buf := new(bytes.Buffer)
		buf.Write(byts[:blen])
		// Create a new buffer with the actual packet
		packet := buf.Bytes()

		// Check for EOF and go out of this loop. Don't cut the connection. Mby we just rebooted the sbc
		if err == io.EOF {
			fmt.Printf("EOF %v\n", err)
			break
		}

		// Check if we have atleast the bytes needed to parse the header
		if len(packet) > 20 {
			// Create a new header struct to get the header length & ID
			set := NewHeader(packet)
			dataLen := int(set.Header.Length)
			setID := int(set.SetHeader.ID)

			// Check if the packet is larger than the header length. If so we have multiple datasets inside one packet
			// Check for known setID's only
			for len(packet) >= dataLen && setID > 255 && setID < 280 {
				// Get the header length from the packet at position 2&3
				dataLen = int(uint16(packet[2])<<8 + uint16(packet[3]))
				// Create a new packet with the header length. This is our first dataset
				data := packet[:dataLen]
				// Cut the first dataset from the original packet
				packet = packet[dataLen:]

				if debug {
					fmt.Println("####################################################################")
					fmt.Printf("Length of incoming packet: %d\n", len(data))
					fmt.Printf("Length from header: %d\n", dataLen)
					fmt.Printf("SetID: %d\n\n", setID)
					fmt.Printf("%s\n", hex.Dump(data))
				}
				// Go through the set's and fill the right structs
				setID = int(uint16(data[16])<<8 + uint16(data[17]))
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
				case 259:
					dataSet := NewSendSipUDP(data)
					if debug {
						fmt.Printf("%s\n", dataSet.Data.SIP.SipMsg)
					}
					SendHEP(dataSet, hconn)
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
					dataSet := NewCallQualityStats(data)

					fmt.Printf("%#v\n", dataSet)
					fmt.Printf("%d\n", dataSet.OutMos)
					fmt.Printf("%d\n", dataSet.IncMos)
					fmt.Printf("%s\n", dataSet.IncCallID)
					fmt.Printf("%s\n", dataSet.OutCallID)

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
