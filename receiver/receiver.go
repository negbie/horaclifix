package receiver

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
type IpfixHeader struct {
	Version       uint16
	Length        uint16
	ExportTime    uint32
	SeqNum        uint32
	ObservationID uint32
}

// SetHeader represents set header fields
type SetHeader struct {
	ID     uint16
	Length uint16
}

type DataSet struct {
	HandShake Hs
	SIP       Su
}

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

// IPFIX is the Oracle IPFIX Handshake Message
type IPFIX struct {
	Header IpfixHeader
	Set    SetHeader
	Data   DataSet
}

type Su struct {
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

// BufToInt8 casts []byte to []int8
func BufToInt8(b []byte) []int8 {
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}
	return bi
}

// NewIPFIX fills the IPFIX struct with structured binary data from r
func NewIPFIX(header []byte) *IPFIX {
	var ni IPFIX
	rni := bytes.NewReader(header)
	binary.Read(rni, binary.BigEndian, &ni.Header.Version)
	binary.Read(rni, binary.BigEndian, &ni.Header.Length)
	binary.Read(rni, binary.BigEndian, &ni.Header.ExportTime)
	binary.Read(rni, binary.BigEndian, &ni.Header.SeqNum)
	binary.Read(rni, binary.BigEndian, &ni.Header.ObservationID)
	binary.Read(rni, binary.BigEndian, &ni.Set.ID)
	binary.Read(rni, binary.BigEndian, &ni.Set.Length)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.MaVer)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.MiVer)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.CFlags1)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.CFlags2)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.SFlags)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.Timeout)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.SystemID)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.Product)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.SMaVer)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.SMiVer)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.Revision)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.HostnameLen)
	ni.Data.HandShake.Hostname = make([]byte, ni.Data.HandShake.HostnameLen)
	binary.Read(rni, binary.BigEndian, &ni.Data.HandShake.Hostname)
	return &ni
}

func NewSipRecUdp(header []byte) *IPFIX {
	var niu IPFIX
	rniu := bytes.NewReader(header)

	binary.Read(rniu, binary.BigEndian, &niu.Header.Version)
	binary.Read(rniu, binary.BigEndian, &niu.Header.Length)
	binary.Read(rniu, binary.BigEndian, &niu.Header.ExportTime)
	binary.Read(rniu, binary.BigEndian, &niu.Header.SeqNum)
	binary.Read(rniu, binary.BigEndian, &niu.Header.ObservationID)
	binary.Read(rniu, binary.BigEndian, &niu.Set.ID)
	binary.Read(rniu, binary.BigEndian, &niu.Set.Length)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TimeSec)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TimeMic)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.IntSlot)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.IntPort)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.IntVlan)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.CallIDLen)
	niu.Data.SIP.CallID = make([]byte, niu.Data.SIP.CallIDLen)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.CallID)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.CallIDEnd)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.IPlen)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.VL)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TOS)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TLen)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TID)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TFlags)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TTL)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TProto)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.TPos)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.SrcIP)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.DstIP)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.DstPort)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.SrcPort)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.UDPlen)
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.MsgLen)
	/*
		fmt.Println(rniu.Len())
		fmt.Println("headerlen:", len(header))
	*/
	niu.Data.SIP.SipMsg = make([]byte, rniu.Len())
	binary.Read(rniu, binary.BigEndian, &niu.Data.SIP.SipMsg)

	return &niu
}

// HandShake writes the binary IPFIX representation into the buffer
func (fi *IPFIX) HandShake() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &fi.Header.Version)
	binary.Write(buf, binary.BigEndian, &fi.Header.Length)
	binary.Write(buf, binary.BigEndian, &fi.Header.ExportTime)
	binary.Write(buf, binary.BigEndian, &fi.Header.SeqNum)
	binary.Write(buf, binary.BigEndian, &fi.Header.ObservationID)
	binary.Write(buf, binary.BigEndian, &fi.Set.ID)
	binary.Write(buf, binary.BigEndian, &fi.Set.Length)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.MaVer)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.MiVer)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.CFlags1)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.CFlags2)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.SFlags)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.Timeout)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.SystemID)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.Product)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.SMaVer)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.SMiVer)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.Revision)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.HostnameLen)
	binary.Write(buf, binary.BigEndian, &fi.Data.HandShake.Hostname)

	return buf.Bytes()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func SyncClient(conn net.Conn) {
	fmt.Println("Handling new connection...")

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	hs := make([]byte, 512)
	n, err := conn.Read(hs)
	checkError(err)

	packet := NewIPFIX(hs[:n])
	if packet.Set.ID == 256 {
		packet.Set.ID++
	}

	fmt.Printf("Write: %#v\n", BufToInt8(packet.HandShake()))

	if packet.Set.ID == 257 && packet.Header.Length > 20 {
		binary.Write(conn, binary.BigEndian, BufToInt8(packet.HandShake()))
	}

	//bufReader := bufio.NewReader(conn)

	for {
		b := make([]byte, 4096)
		n, err := conn.Read(b)
		checkError(err)

		data := NewSipRecUdp(b[:n])

		/*		// Read IPFIX Messages delimited by newline
				bytes, err := bufReader.ReadBytes('\n')
				checkError(err)*/

		//fmt.Printf("%s", (binary.BigEndian.Uint16(bytes)))
		fmt.Printf("%s", data.Data.SIP)

	}
}
