package main

import (
	"crypto/tls"
	"database/sql"
	"net"
	"sync"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/negbie/sipparser"
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

// HandShake holds the HandShake header and dataset fields
type HandShake struct {
	IpfixHeader
	IpfixSetHeader
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

// IPFIX messages struct
type IPFIX struct {
	DataSet
}

// DataSet holds different IPFIX datasets
type DataSet struct {
	SIP SipSet
	QOS QosSet
}

// SipSet holds the SIP related fields
type SipSet struct {
	TimeSec     uint32
	TimeMic     uint32
	IntSlot     uint8
	IntPort     uint8
	IntVlan     uint16
	CallIDLen   uint8
	CallID      []byte
	CallIDEnd   uint8
	IPlen       uint16
	VL          uint8
	TOS         uint8
	TLen        uint16
	TID         uint16
	TFlags      uint16
	TTL         uint8
	TProto      uint8
	TPos        uint16
	SrcIP       net.IP
	SrcIPString string
	DstIP       net.IP
	DstIPString string
	DstPort     uint16
	SrcPort     uint16
	Context     uint32
	UDPlen      uint16
	MsgLen      uint16
	RawMsg      []byte
	SipMsg      *sipparser.SipMsg
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

	EndTimeSec uint32
	EndTimeMic uint32

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

	CallerSsrc  uint32
	CallerSeqNo uint32
	CallerPt    uint8

	CalleeSsrc  uint32
	CalleeSeqNo uint32
	CalleePt    uint8

	TransID uint32
}

type Connections struct {
	Graylog
	Homer
	StatsD *net.UDPConn
	MySQL  *mysqlDB
	Influx *InfluxClient
}

type Graylog struct {
	TCP *net.TCPConn
	UDP *net.UDPConn
	*sync.RWMutex
	disconnected bool
}

type Homer struct {
	TCP *net.TCPConn
	UDP *net.UDPConn
	TLS *tls.Conn
	*sync.RWMutex
	disconnected bool
}

type InfluxMetric struct {
	measurement string
	tags        map[string]string
	fields      map[string]interface{}
	time        time.Time
}

type InfluxClient struct {
	client        influx.Client
	database      string
	batchSize     int
	pointsChannel chan *influx.Point
	batchConfig   influx.BatchPointsConfig
	errorFunc     func(err error)
}

type InfluxClientConfig struct {
	Endpoint     string
	Database     string
	BatchSize    int
	FlushTimeout time.Duration
	ErrorFunc    func(err error)
}

type mysqlDB struct {
	conn   *sql.DB
	insert *sql.Stmt
}
