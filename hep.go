package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

// SendHEPMsg writes the final HEP packet into the buffer
func SendHEPMsg(msg []byte) {
	var packet []byte
	packet = make([]byte, len(msg)+6)

	copy(packet[6:], msg)
	binary.BigEndian.PutUint32(packet[:4], uint32(0x48455033))   // ASCII "HEP3"
	binary.BigEndian.PutUint16(packet[4:6], uint16(len(packet))) // Total packet length

	if conn, err := net.Dial("udp", *haddr); err == nil {
		conn.Write(packet)
		conn.Close()
	}
}

// NewHEPChunck constructs the HEP chunck
func (ipfix *IPFIX) NewHEPChunck(ChunckVen uint16, ChunckType uint16, payloadType string) []byte {

	var packet []byte

	switch ChunckType {
	// Chunk IP protocol family (0x02=IPv4)
	case 0x0001:
		packet = make([]byte, 6+1)
		packet[6] = 0x02

	// Chunk IP protocol ID (0x11=UDP)
	case 0x0002:
		packet = make([]byte, 6+1)
		packet[6] = 0x11

	// Chunk IPv4 source address
	case 0x0003:
		packet = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.SIP.SrcIP)
		case "incRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CallerIncSrcIP)
		case "outRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CalleeIncSrcIP)
		case "incRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CallerIncSrcIP)
		case "outRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CalleeIncSrcIP)
		}

	// Chunk IPv4 destination address
	case 0x0004:
		packet = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.SIP.DstIP)
		case "incRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CallerIncDstIP)
		case "outRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CalleeIncDstIP)
		case "incRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CallerIncDstIP)
		case "outRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.CalleeIncDstIP)
		}

	// Chunk IPv6 source address
	// case 0x0005:

	// Chunk IPv6 destination address
	// case 0x0006:

	// Chunk protocol source port
	case 0x0007:
		packet = make([]byte, 6+2)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.SIP.SrcPort)
		case "incRTP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CallerIncSrcPort)
		case "outRTP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CalleeIncSrcPort)
		case "incRTCP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CallerIncSrcPort)
		case "outRTCP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CalleeIncSrcPort)
		}

	// Chunk destination source port
	case 0x0008:
		packet = make([]byte, 6+2)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.SIP.DstPort)
		case "incRTP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CallerIncDstPort)
		case "outRTP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CalleeIncDstPort)
		case "incRTCP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CallerIncDstPort)
		case "outRTCP":
			binary.BigEndian.PutUint16(packet[6:], ipfix.Data.QOS.CalleeIncDstPort)
		}

	// Chunk unix timestamp, seconds
	case 0x0009:
		packet = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.SIP.TimeSec)
		case "incRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeSec)
		case "outRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeSec)
		case "incRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeSec)
		case "outRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeSec)
		}

	// Chunk unix timestamp, microseconds offset
	case 0x000a:
		packet = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.SIP.TimeMic)
		case "incRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeMic)
		case "outRTP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeMic)
		case "incRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeMic)
		case "outRTCP":
			binary.BigEndian.PutUint32(packet[6:], ipfix.Data.QOS.BeginTimeMic)
		}

	// Chunk protocol type (SIP/H323/RTP/MGCP/M2UA)
	case 0x000b:
		packet = make([]byte, 6+1)
		switch payloadType {
		case "SIP":
			packet[6] = 1
		case "incRTP":
			packet[6] = 34
		case "outRTP":
			packet[6] = 34
		case "incRTCP":
			packet[6] = 35
		case "outRTCP":
			packet[6] = 35
		}

	// Chunk capture agent ID
	case 0x000c:
		packet = make([]byte, 6+4)
		binary.BigEndian.PutUint32(packet[6:], 0x00000BEE) // Node homer01:3054

	// case 0x000d:
	// Chunk keep alive timer

	// case 0x000e:
	// Chunk authenticate key (plain text / TLS connection)

	// Chunk captured packet payload
	case 0x000f:
		packet = make([]byte, len(ipfix.Data.SIP.SipMsg)+6)
		switch payloadType {
		case "SIP":
			copy(packet[6:], ipfix.Data.SIP.SipMsg)
		case "incRTP":
			payload, _ := ipfix.PrepIncRtp()
			copy(packet[6:], payload)
		case "outRTP":
			payload, _ := ipfix.PrepOutRtp()
			copy(packet[6:], payload)
		case "incRTCP":
			payload, _ := ipfix.PrepIncRtcp()
			copy(packet[6:], payload)
		case "outRTCP":
			payload, _ := ipfix.PrepIncRtcp()
			copy(packet[6:], payload)
		}

	// case 0x0010:
	// Chunk captured compressed payload (gzip/inflate)

	// Chunk internal correlation id
	case 0x0011:
		packet = make([]byte, len(ipfix.Data.QOS.IncCallID)+6)
		copy(packet[6:], ipfix.Data.QOS.IncCallID)
		log.Printf("Case 0x0011 Packet: %s", ipfix.Data.QOS.IncCallID)

	}

	binary.BigEndian.PutUint16(packet[:2], ChunckVen)
	binary.BigEndian.PutUint16(packet[2:4], ChunckType)
	binary.BigEndian.PutUint16(packet[4:6], uint16(len(packet)))

	return packet
}

// SendHEP sends the HEP message
func NewSipHEP(i *IPFIX) {
	bhep := new(bytes.Buffer)

	bhep.Write(i.NewHEPChunck(0x0000, 0x0001, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0002, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0003, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0004, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0007, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0008, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0009, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000a, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000b, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000c, "SIP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000f, "SIP"))
	SendHEPMsg(bhep.Bytes())
}

func NewQosHEPincRTP(i *IPFIX) {
	bhep := new(bytes.Buffer)

	bhep.Write(i.NewHEPChunck(0x0000, 0x0001, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0002, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0003, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0004, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0007, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0008, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0009, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000a, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000b, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000c, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000f, "incRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0011, "incRTP"))
	SendHEPMsg(bhep.Bytes())

}

func NewQosHEPoutRTP(i *IPFIX) {
	bhep := new(bytes.Buffer)

	bhep.Write(i.NewHEPChunck(0x0000, 0x0001, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0002, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0003, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0004, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0007, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0008, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0009, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000a, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000b, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000c, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000f, "outRTP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0011, "outRTP"))
	SendHEPMsg(bhep.Bytes())

}

func NewQosHEPincRTCP(i *IPFIX) {
	bhep := new(bytes.Buffer)

	bhep.Write(i.NewHEPChunck(0x0000, 0x0001, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0002, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0003, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0004, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0007, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0008, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0009, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000a, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000b, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000c, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000f, "incRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0011, "incRTCP"))
	SendHEPMsg(bhep.Bytes())

}
func NewQosHEPoutRTCP(i *IPFIX) {
	bhep := new(bytes.Buffer)

	bhep.Write(i.NewHEPChunck(0x0000, 0x0001, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0002, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0003, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0004, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0007, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0008, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0009, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000a, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000b, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000c, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x000f, "outRTCP"))
	bhep.Write(i.NewHEPChunck(0x0000, 0x0011, "outRTCP"))
	SendHEPMsg(bhep.Bytes())

}

// PrepIncRtp
func (ipfix *IPFIX) PrepIncRtp() ([]byte, error) {
	mapIncRtp := map[string]interface{}{

		"CORRELATION_ID":  string(ipfix.Data.QOS.IncCallID),
		"RTP_SIP_CALL_ID": string(ipfix.Data.QOS.IncCallID),
		"REPORT_TS":       ipfix.Data.QOS.BeginTimeSec,
		"TL_BYTE":         ipfix.Data.QOS.IncRtpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        ipfix.Data.QOS.IncRtpPackets,
		"EXPECTED_PK":     (ipfix.Data.QOS.IncRtpPackets + ipfix.Data.QOS.IncRtpLostPackets),
		"PACKET_LOSS":     ipfix.Data.QOS.IncRtpLostPackets,
		"SEQ":             0,
		"JITTER":          ipfix.Data.QOS.IncRtpAvgJitter,
		"MAX_JITTER":      ipfix.Data.QOS.IncRtpMaxJitter,
		"MEAN_JITTER":     ipfix.Data.QOS.IncRtpAvgJitter,
		"DELTA":           0.000,
		"MAX_DELTA":       0.000,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         ipfix.Data.QOS.IncMos,
		"MEAN_MOS":        ipfix.Data.QOS.IncMos,
		"MOS":             ipfix.Data.QOS.IncMos,
		"RFACTOR":         ipfix.Data.QOS.IncrVal,
		"MIN_RFACTOR":     ipfix.Data.QOS.IncrVal,
		"MEAN_RFACTOR":    ipfix.Data.QOS.IncrVal,
		"SRC_IP":          ipfix.Data.QOS.CallerIncSrcIP,
		"SRC_PORT":        ipfix.Data.QOS.CallerIncSrcPort,
		"DST_IP":          ipfix.Data.QOS.CallerIncDstIP,
		"DST_PORT":        ipfix.Data.QOS.CallerIncSrcPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        ipfix.Data.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      ipfix.Data.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(ipfix.Data.QOS.IncRealm),
		"PARTY":           0,
		"TYPE":            "PERIODIC",
	}
	j, err := json.Marshal(mapIncRtp)

	log.Printf("%s\n", j)

	return j, err

}

func (ipfix *IPFIX) PrepOutRtp() ([]byte, error) {
	mapOutRtp := map[string]interface{}{

		"CORRELATION_ID":  string(ipfix.Data.QOS.OutCallID),
		"RTP_SIP_CALL_ID": string(ipfix.Data.QOS.OutCallID),
		"REPORT_TS":       ipfix.Data.QOS.BeginTimeSec,
		"TL_BYTE":         ipfix.Data.QOS.OutRtpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        ipfix.Data.QOS.OutRtpPackets,
		"EXPECTED_PK":     (ipfix.Data.QOS.OutRtpPackets + ipfix.Data.QOS.OutRtpLostPackets),
		"PACKET_LOSS":     ipfix.Data.QOS.OutRtpLostPackets,
		"SEQ":             0,
		"JITTER":          ipfix.Data.QOS.OutRtpAvgJitter,
		"MAX_JITTER":      ipfix.Data.QOS.OutRtpMaxJitter,
		"MEAN_JITTER":     ipfix.Data.QOS.OutRtpAvgJitter,
		"DELTA":           0.000,
		"MAX_DELTA":       0.000,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         ipfix.Data.QOS.OutMos,
		"MEAN_MOS":        ipfix.Data.QOS.OutMos,
		"MOS":             ipfix.Data.QOS.OutMos,
		"RFACTOR":         ipfix.Data.QOS.OutrVal,
		"MIN_RFACTOR":     ipfix.Data.QOS.OutrVal,
		"MEAN_RFACTOR":    ipfix.Data.QOS.OutrVal,
		"SRC_IP":          ipfix.Data.QOS.CalleeOutSrcIP,
		"SRC_PORT":        ipfix.Data.QOS.CalleeOutSrcPort,
		"DST_IP":          ipfix.Data.QOS.CalleeOutDstIP,
		"DST_PORT":        ipfix.Data.QOS.CalleeOutSrcPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        ipfix.Data.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      ipfix.Data.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(ipfix.Data.QOS.OutRealm),
		"PARTY":           1,
		"TYPE":            "PERIODIC",
	}

	j, err := json.Marshal(mapOutRtp)

	log.Printf("%s\n", j)

	return j, err

}

func (ipfix *IPFIX) PrepIncRtcp() ([]byte, error) {
	mapIncRtcp := map[string]interface{}{

		"CORRELATION_ID":  string(ipfix.Data.QOS.IncCallID),
		"RTP_SIP_CALL_ID": string(ipfix.Data.QOS.IncCallID),
		"REPORT_TS":       ipfix.Data.QOS.BeginTimeSec,
		"TL_BYTE":         ipfix.Data.QOS.IncRtcpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        ipfix.Data.QOS.IncRtcpPackets,
		"EXPECTED_PK":     (ipfix.Data.QOS.IncRtcpPackets + ipfix.Data.QOS.IncRtcpLostPackets),
		"PACKET_LOSS":     ipfix.Data.QOS.IncRtcpLostPackets,
		"SEQ":             0,
		"JITTER":          ipfix.Data.QOS.IncRtcpAvgJitter,
		"MAX_JITTER":      ipfix.Data.QOS.IncRtcpMaxJitter,
		"MEAN_JITTER":     ipfix.Data.QOS.IncRtcpAvgJitter,
		"DELTA":           ipfix.Data.QOS.IncRtcpAvgLat,
		"MAX_DELTA":       ipfix.Data.QOS.IncRtcpMaxLat,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         ipfix.Data.QOS.IncMos,
		"MEAN_MOS":        ipfix.Data.QOS.IncMos,
		"MOS":             ipfix.Data.QOS.IncMos,
		"RFACTOR":         ipfix.Data.QOS.IncrVal,
		"MIN_RFACTOR":     ipfix.Data.QOS.IncrVal,
		"MEAN_RFACTOR":    ipfix.Data.QOS.IncrVal,
		"SRC_IP":          ipfix.Data.QOS.CallerIncSrcIP,
		"SRC_PORT":        ipfix.Data.QOS.CallerIncSrcPort,
		"DST_IP":          ipfix.Data.QOS.CallerIncDstIP,
		"DST_PORT":        ipfix.Data.QOS.CallerIncSrcPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        ipfix.Data.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      ipfix.Data.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(ipfix.Data.QOS.IncRealm),
		"PARTY":           0,
		"TYPE":            "PERIODIC",
	}

	return json.Marshal(mapIncRtcp)

}

func (ipfix *IPFIX) PrepOutRtcp() ([]byte, error) {
	mapOutRtcp := map[string]interface{}{

		"CORRELATION_ID":  string(ipfix.Data.QOS.OutCallID),
		"RTP_SIP_CALL_ID": string(ipfix.Data.QOS.OutCallID),
		"REPORT_TS":       ipfix.Data.QOS.BeginTimeSec,
		"TL_BYTE":         ipfix.Data.QOS.OutRtcpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        ipfix.Data.QOS.OutRtcpPackets,
		"EXPECTED_PK":     (ipfix.Data.QOS.OutRtcpPackets + ipfix.Data.QOS.OutRtcpLostPackets),
		"PACKET_LOSS":     ipfix.Data.QOS.OutRtcpLostPackets,
		"SEQ":             0,
		"JITTER":          ipfix.Data.QOS.OutRtcpAvgJitter,
		"MAX_JITTER":      ipfix.Data.QOS.OutRtcpMaxJitter,
		"MEAN_JITTER":     ipfix.Data.QOS.OutRtcpAvgJitter,
		"DELTA":           ipfix.Data.QOS.OutRtcpAvgLat,
		"MAX_DELTA":       ipfix.Data.QOS.OutRtcpMaxLat,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         ipfix.Data.QOS.OutMos,
		"MEAN_MOS":        ipfix.Data.QOS.OutMos,
		"MOS":             ipfix.Data.QOS.OutMos,
		"RFACTOR":         ipfix.Data.QOS.OutrVal,
		"MIN_RFACTOR":     ipfix.Data.QOS.OutrVal,
		"MEAN_RFACTOR":    ipfix.Data.QOS.OutrVal,
		"SRC_IP":          ipfix.Data.QOS.CalleeOutSrcIP,
		"SRC_PORT":        ipfix.Data.QOS.CalleeOutSrcPort,
		"DST_IP":          ipfix.Data.QOS.CalleeOutDstIP,
		"DST_PORT":        ipfix.Data.QOS.CalleeOutSrcPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        ipfix.Data.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      ipfix.Data.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(ipfix.Data.QOS.OutRealm),
		"PARTY":           1,
		"TYPE":            "PERIODIC",
	}

	return json.Marshal(mapOutRtcp)

}
