package main

import (
	"crypto/tls"
	"net"
)

// IPFIX holds the structure of one IPFIX packet
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

// DataSet holds different IPFIX datasets
type DataSet struct {
	Hs  HandShake
	SIP SipSet
	QOS QosSet
}

// HandShake holds the HandShake dataset fields
type HandShake struct {
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
	//Hostname    ByteString
}

// SipSet holds the SIP related fields
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

// QosSet holds the QoS related fields
type QosSet struct {
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

type Connections struct {
	Graylog    net.Conn
	Homer      net.Conn
	StatsD     net.Conn
	Banshee    net.Conn
	GraylogTLS *tls.Conn
}
