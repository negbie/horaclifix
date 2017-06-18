package receiver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
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
type iheader struct {
	Version       uint16
	Length        uint16
	ExportTime    uint32
	SeqNum        uint32
	ObservationID uint32
	SetID         uint16
	SetLen        uint16
	MaVer         uint16
	MiVer         uint16
	CFlags1       uint16
	CFlags2       uint16
	SFlags        uint16
	Timeout       uint16
	SystemID      uint32
	Product       uint16
	SMaVer        uint8
	SMiVer        uint8
	Revision      uint8
	Hostname      []byte
}

// bufToInt8 casts byte to int8
func bufToInt8(b []byte) []int8 {

	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}
	return bi
}

// Get fills the Iheader struct with structured binary data from r
func get(header []byte) *iheader {
	var ih iheader
	r := bytes.NewReader(header)
	binary.Read(r, binary.BigEndian, &ih.Version)
	binary.Read(r, binary.BigEndian, &ih.Length)
	binary.Read(r, binary.BigEndian, &ih.ExportTime)
	binary.Read(r, binary.BigEndian, &ih.SeqNum)
	binary.Read(r, binary.BigEndian, &ih.ObservationID)
	binary.Read(r, binary.BigEndian, &ih.SetID)
	binary.Read(r, binary.BigEndian, &ih.SetLen)
	binary.Read(r, binary.BigEndian, &ih.MaVer)
	binary.Read(r, binary.BigEndian, &ih.MiVer)
	binary.Read(r, binary.BigEndian, &ih.CFlags1)
	binary.Read(r, binary.BigEndian, &ih.CFlags2)
	binary.Read(r, binary.BigEndian, &ih.SFlags)
	binary.Read(r, binary.BigEndian, &ih.Timeout)
	binary.Read(r, binary.BigEndian, &ih.SystemID)
	binary.Read(r, binary.BigEndian, &ih.Product)
	binary.Read(r, binary.BigEndian, &ih.SMaVer)
	binary.Read(r, binary.BigEndian, &ih.SMiVer)
	binary.Read(r, binary.BigEndian, &ih.Revision)
	ih.Hostname = header[41:len(header)]
	return &ih
}

// Set writes the binary Iheader representation into the buffer
func (ih *iheader) set() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &ih.Version)
	binary.Write(buf, binary.BigEndian, &ih.Length)
	binary.Write(buf, binary.BigEndian, &ih.ExportTime)
	binary.Write(buf, binary.BigEndian, &ih.SeqNum)
	binary.Write(buf, binary.BigEndian, &ih.ObservationID)
	binary.Write(buf, binary.BigEndian, &ih.SetID)
	binary.Write(buf, binary.BigEndian, &ih.SetLen)
	binary.Write(buf, binary.BigEndian, &ih.MaVer)
	binary.Write(buf, binary.BigEndian, &ih.MiVer)
	binary.Write(buf, binary.BigEndian, &ih.CFlags1)
	binary.Write(buf, binary.BigEndian, &ih.CFlags2)
	binary.Write(buf, binary.BigEndian, &ih.SFlags)
	binary.Write(buf, binary.BigEndian, &ih.Timeout)
	binary.Write(buf, binary.BigEndian, &ih.SystemID)
	binary.Write(buf, binary.BigEndian, &ih.Product)
	binary.Write(buf, binary.BigEndian, &ih.SMaVer)
	binary.Write(buf, binary.BigEndian, &ih.SMiVer)
	binary.Write(buf, binary.BigEndian, &ih.Revision)
	binary.Write(buf, binary.BigEndian, &ih.Hostname)

	out := buf.Bytes()
	return out
}

func SyncClient(conn net.Conn) {
	defer conn.Close()
	ib := make([]byte, 1024)
	n, _ := conn.Read(ib)

	sync := get(ib[:n])
	if sync.SetID == 256 {
		sync.SetID++
	}
	fmt.Printf("msg to send: %#v", bufToInt8(sync.set()))
	if sync.SetID == 257 {
		binary.Write(conn, binary.BigEndian, bufToInt8(sync.set()))
	}
}
