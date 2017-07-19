package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

// SendHep will write the final HEP message into the buffer and send it to wire
// The first 4 bytes are the string "HEP3". The next 2 bytes are the length of the
// whole message (len("HEP3") + length of all the chucks we have). The next bytes
// are all the chuncks created by NewHEPChuncks()
// Bytes: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31......
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//        | "HEP3"|len|chuncks(0x0001|0x0002|0x0003|0x0004|0x0007|0x0008|0x0009|0x000a|0x000b|......)
//        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
func (conn *Connections) SendHep(i *IPFIX, s string) {
	chuncks := i.NewHEPChuncks(s)
	hepMsg := make([]byte, len(chuncks)+6)
	copy(hepMsg[6:], chuncks)
	binary.BigEndian.PutUint32(hepMsg[:4], uint32(0x48455033))   // ASCII "HEP3"
	binary.BigEndian.PutUint16(hepMsg[4:6], uint16(len(hepMsg))) // Total length
	conn.Homer.Write(hepMsg)
}

// MakeChunck will construct the respective HEP chunck
func (i *IPFIX) MakeChunck(chunckVen uint16, chunckType uint16, payloadType string) []byte {
	var chunck []byte
	switch chunckType {
	// Chunk IP protocol family (0x02=IPv4)
	case 0x0001:
		chunck = make([]byte, 6+1)
		chunck[6] = 0x02

	// Chunk IP protocol ID (0x11=UDP)
	case 0x0002:
		chunck = make([]byte, 6+1)
		chunck[6] = 0x11

	// Chunk IPv4 source address
	case 0x0003:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.SrcIP)
		case "QOS":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.CallerIncSrcIP)
		case "logQOS":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.CallerIncSrcIP)
		}

	// Chunk IPv4 destination address
	case 0x0004:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.DstIP)
		case "QOS":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.CallerIncDstIP)
		case "logQOS":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.CallerIncDstIP)
		}

	// Chunk IPv6 source address
	// case 0x0005:

	// Chunk IPv6 destination address
	// case 0x0006:

	// Chunk protocol source port
	case 0x0007:
		chunck = make([]byte, 6+2)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.SIP.SrcPort)
		case "QOS":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.QOS.CallerIncSrcPort)
		case "logQOS":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.QOS.CallerIncSrcPort)
		}

	// Chunk destination source port
	case 0x0008:
		chunck = make([]byte, 6+2)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.SIP.DstPort)
		case "QOS":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.QOS.CallerIncDstPort)
		case "logQOS":
			binary.BigEndian.PutUint16(chunck[6:], i.Data.QOS.CallerIncDstPort)
		}

	// Chunk unix timestamp, seconds
	case 0x0009:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.TimeSec)
		default:
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.EndTimeSec)
		}

	// Chunk unix timestamp, microseconds offset
	case 0x000a:
		chunck = make([]byte, 6+4)
		switch payloadType {
		case "SIP":
			binary.BigEndian.PutUint32(chunck[6:], i.Data.SIP.TimeMic)
		default:
			binary.BigEndian.PutUint32(chunck[6:], i.Data.QOS.EndinTimeMic)
		}

	// Chunk protocol type (SIP/H323/RTP/MGCP/M2UA)
	case 0x000b:
		chunck = make([]byte, 6+1)
		switch payloadType {
		case "SIP":
			chunck[6] = 1
		case "QOS":
			chunck[6] = 33
		case "logQOS":
			chunck[6] = 100
		}

	// Chunk capture agent ID
	case 0x000c:
		chunck = make([]byte, 6+4)
		binary.BigEndian.PutUint32(chunck[6:], 0x00000BEE)

		// case 0x000d:
		// Chunk keep alive timer

	// Chunk authenticate key (plain text / TLS connection)
	case 0x000e:
		chunck = make([]byte, len(*hepPw)+6)
		copy(chunck[6:], *hepPw)

	// Chunk captured packet payload
	case 0x000f:
		switch payloadType {
		case "SIP":
			chunck = make([]byte, len(i.Data.SIP.SipMsg)+6)
			copy(chunck[6:], i.Data.SIP.SipMsg)
		case "QOS":
			payload, _ := json.Marshal(i.PrepQOS())
			chunck = make([]byte, len(payload)+6)
			copy(chunck[6:], payload)
		case "logQOS":
			payload, _ := json.Marshal(i.PrepLogQOS())
			chunck = make([]byte, len(payload)+6)
			copy(chunck[6:], payload)
		}

	// case 0x0010:
	// Chunk captured compressed payload (gzip/inflate)

	// Chunk internal correlation id
	case 0x0011:
		if len(i.Data.QOS.IncCallID) > 0 {
			chunck = make([]byte, len(i.Data.QOS.IncCallID)+6)
			copy(chunck[6:], i.Data.QOS.IncCallID)
		} else {
			chunck = make([]byte, len(i.Data.QOS.OutCallID)+6)
			copy(chunck[6:], i.Data.QOS.OutCallID)
		}

		/*
		   // TODO rewrite cast to uint16
		   // Chunk MOS
		   	case 0x0020:
		   		chunck = make([]byte, 6+2)
		   		switch payloadType {
		   		case "incRTP":
		   			binary.BigEndian.PutUint16(chunck[6:], uint16(i.Data.QOS.IncMos))
		   		case "outRTP":
		   			binary.BigEndian.PutUint16(chunck[6:], uint16(i.Data.QOS.OutMos))
		   		case "incRTCP":
		   			binary.BigEndian.PutUint16(chunck[6:], uint16(i.Data.QOS.IncMos))
		   		case "outRTCP":
		   			binary.BigEndian.PutUint16(chunck[6:], uint16(i.Data.QOS.OutMos))
		   		}
		*/
	}
	binary.BigEndian.PutUint16(chunck[:2], chunckVen)
	binary.BigEndian.PutUint16(chunck[2:4], chunckType)
	binary.BigEndian.PutUint16(chunck[4:6], uint16(len(chunck)))
	return chunck
}

// NewHEPChuncks will fill a buffer with all the chuncks
// we need and give this buffer to SendHepMsg
func (i *IPFIX) NewHEPChuncks(s string) []byte {
	buf := new(bytes.Buffer)

	buf.Write(i.MakeChunck(0x0000, 0x0001, s))
	buf.Write(i.MakeChunck(0x0000, 0x0002, s))
	buf.Write(i.MakeChunck(0x0000, 0x0003, s))
	buf.Write(i.MakeChunck(0x0000, 0x0004, s))
	buf.Write(i.MakeChunck(0x0000, 0x0007, s))
	buf.Write(i.MakeChunck(0x0000, 0x0008, s))
	buf.Write(i.MakeChunck(0x0000, 0x0009, s))
	buf.Write(i.MakeChunck(0x0000, 0x000a, s))
	buf.Write(i.MakeChunck(0x0000, 0x000b, s))
	buf.Write(i.MakeChunck(0x0000, 0x000c, s))
	buf.Write(i.MakeChunck(0x0000, 0x000e, s))
	buf.Write(i.MakeChunck(0x0000, 0x000f, s))
	if s == "QOS" || s == "logQOS" {
		buf.Write(i.MakeChunck(0x0000, 0x0011, s))
		//buf.Write(i.MakeChunck(0x0000, 0x0020, s))
	}
	return buf.Bytes()
}

// PrepQOS
func (i *IPFIX) PrepQOS() *map[string]interface{} {
	mapQOS := map[string]interface{}{
		"INC_ID":              string(i.Data.QOS.IncCallID),
		"INC_RTP_BYTE":        i.Data.QOS.IncRtpBytes,
		"INC_RTP_PK":          i.Data.QOS.IncRtpPackets,
		"INC_RTP_PK_LOSS":     i.Data.QOS.IncRtpLostPackets,
		"INC_RTP_AVG_JITTER":  i.Data.QOS.IncRtpAvgJitter,
		"INC_RTP_MAX_JITTER":  i.Data.QOS.IncRtpMaxJitter,
		"INC_RTCP_BYTE":       i.Data.QOS.IncRtcpBytes,
		"INC_RTCP_PK":         i.Data.QOS.IncRtcpPackets,
		"INC_RTCP_PK_LOSS":    i.Data.QOS.IncRtcpLostPackets,
		"INC_RTCP_AVG_JITTER": i.Data.QOS.IncRtcpAvgJitter,
		"INC_RTCP_MAX_JITTER": i.Data.QOS.IncRtcpMaxJitter,
		"INC_RTCP_AVG_LAT":    i.Data.QOS.IncRtcpAvgLat,
		"INC_RTCP_MAX_LAT":    i.Data.QOS.IncRtcpMaxLat,
		"INC_MOS":             i.Data.QOS.IncMos,
		"INC_RVAL":            i.Data.QOS.IncrVal,
		"CALLER_VLAN":         i.Data.QOS.CallerIntVlan,
		"CALLER_SRC_IP":       stringIPv4(i.Data.QOS.CallerIncSrcIP),
		"CALLER_SRC_PORT":     i.Data.QOS.CallerIncSrcPort,
		"CALLER_DST_IP":       stringIPv4(i.Data.QOS.CallerOutDstIP),
		"CALLER_DST_PORT":     i.Data.QOS.CallerOutDstPort,
		"INC_REALM":           string(i.Data.QOS.IncRealm),

		"OUT_ID":              string(i.Data.QOS.OutCallID),
		"OUT_RTP_BYTE":        i.Data.QOS.OutRtpBytes,
		"OUT_RTP_PK":          i.Data.QOS.OutRtpPackets,
		"OUT_RTP_PK_LOSS":     i.Data.QOS.OutRtpLostPackets,
		"OUT_RTP_AVG_JITTER":  i.Data.QOS.OutRtpAvgJitter,
		"OUT_RTP_MAX_JITTER":  i.Data.QOS.OutRtpMaxJitter,
		"OUT_RTCP_BYTE":       i.Data.QOS.OutRtcpBytes,
		"OUT_RTCP_PK":         i.Data.QOS.OutRtcpPackets,
		"OUT_RTCP_PK_LOSS":    i.Data.QOS.OutRtcpLostPackets,
		"OUT_RTCP_AVG_JITTER": i.Data.QOS.OutRtcpAvgJitter,
		"OUT_RTCP_MAX_JITTER": i.Data.QOS.OutRtcpMaxJitter,
		"OUT_RTCP_AVG_LAT":    i.Data.QOS.OutRtcpAvgLat,
		"OUT_RTCP_MAX_LAT":    i.Data.QOS.OutRtcpMaxLat,
		"OUT_MOS":             i.Data.QOS.OutMos,
		"OUT_RVAL":            i.Data.QOS.OutrVal,
		"CALLEE_VLAN":         i.Data.QOS.CalleeIntVlan,
		"CALLEE_SRC_IP":       stringIPv4(i.Data.QOS.CalleeOutSrcIP),
		"CALLEE_SRC_PORT":     i.Data.QOS.CalleeOutSrcPort,
		"CALLEE_DST_IP":       stringIPv4(i.Data.QOS.CalleeIncDstIP),
		"CALLEE_DST_PORT":     i.Data.QOS.CalleeIncDstPort,
		"OUT_REALM":           string(i.Data.QOS.OutRealm),
		"MEDIA_TYPE":          i.Data.QOS.Type,
	}
	return &mapQOS
}

/*
func (i *IPFIX) PrepIncRTCP() *map[string]interface{} {
	mapIncRtcp := map[string]interface{}{

		"CORRELATION_ID":  string(i.Data.QOS.IncCallID),
		"RTP_SIP_CALL_ID": string(i.Data.QOS.IncCallID),
		"REPORT_TS":       i.Data.QOS.BeginTimeSec,
		"TL_BYTE":         i.Data.QOS.IncRtcpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        i.Data.QOS.IncRtcpPackets,
		"EXPECTED_PK":     (i.Data.QOS.IncRtcpPackets + i.Data.QOS.IncRtcpLostPackets),
		"PACKET_LOSS":     i.Data.QOS.IncRtcpLostPackets,
		"SEQ":             0,
		"JITTER":          i.Data.QOS.IncRtcpAvgJitter,
		"MAX_JITTER":      i.Data.QOS.IncRtcpMaxJitter,
		"MEAN_JITTER":     i.Data.QOS.IncRtcpAvgJitter,
		"DELTA":           i.Data.QOS.IncRtcpAvgLat,
		"MAX_DELTA":       i.Data.QOS.IncRtcpMaxLat,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         i.Data.QOS.IncMos,
		"MEAN_MOS":        i.Data.QOS.IncMos,
		"MOS":             i.Data.QOS.IncMos,
		"RFACTOR":         i.Data.QOS.IncrVal,
		"MIN_RFACTOR":     i.Data.QOS.IncrVal,
		"MEAN_RFACTOR":    i.Data.QOS.IncrVal,
		"SRC_IP":          stringIPv4(i.Data.QOS.CallerIncSrcIP),
		"SRC_PORT":        i.Data.QOS.CallerIncSrcPort,
		"DST_IP":          stringIPv4(i.Data.QOS.CallerIncDstIP),
		"DST_PORT":        i.Data.QOS.CallerIncSrcPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        i.Data.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      i.Data.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(i.Data.QOS.IncRealm),
		"PARTY":           0,
		"TYPE":            "PERIODIC",
	} return &mapIncRtcp
}

func (i *IPFIX) PrepOutRTCP() *map[string]interface{} {
	mapOutRtcp := map[string]interface{}{

		"CORRELATION_ID":  string(i.Data.QOS.OutCallID),
		"RTP_SIP_CALL_ID": string(i.Data.QOS.OutCallID),
		"REPORT_TS":       i.Data.QOS.BeginTimeSec,
		"TL_BYTE":         i.Data.QOS.OutRtcpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        i.Data.QOS.OutRtcpPackets,
		"EXPECTED_PK":     (i.Data.QOS.OutRtcpPackets + i.Data.QOS.OutRtcpLostPackets),
		"PACKET_LOSS":     i.Data.QOS.OutRtcpLostPackets,
		"SEQ":             0,
		"JITTER":          i.Data.QOS.OutRtcpAvgJitter,
		"MAX_JITTER":      i.Data.QOS.OutRtcpMaxJitter,
		"MEAN_JITTER":     i.Data.QOS.OutRtcpAvgJitter,
		"DELTA":           i.Data.QOS.OutRtcpAvgLat,
		"MAX_DELTA":       i.Data.QOS.OutRtcpMaxLat,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         i.Data.QOS.OutMos,
		"MEAN_MOS":        i.Data.QOS.OutMos,
		"MOS":             i.Data.QOS.OutMos,
		"RFACTOR":         i.Data.QOS.OutrVal,
		"MIN_RFACTOR":     i.Data.QOS.OutrVal,
		"MEAN_RFACTOR":    i.Data.QOS.OutrVal,
		"SRC_IP":          stringIPv4(i.Data.QOS.CalleeOutSrcIP),
		"SRC_PORT":        i.Data.QOS.CalleeOutSrcPort,
		"DST_IP":          stringIPv4(i.Data.QOS.CalleeOutDstIP),
		"DST_PORT":        i.Data.QOS.CalleeOutSrcPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        i.Data.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      i.Data.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(i.Data.QOS.OutRealm),
		"PARTY":           1,
		"TYPE":            "PERIODIC",
	} return &mapOutRtcp
}
*/
