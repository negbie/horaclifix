package main

import (
	"bytes"
	"strconv"
)

func (i *IPFIX) mapLogSIP() *map[string]interface{} {
	m := map[string]interface{}{
		"version":       "1.1",
		"host":          *name,
		"short_message": i.SIP.SipMsg.Msg,
		"level":         5,
		"_id":           i.SIP.SipMsg.CallID,
		"_from":         i.SIP.SipMsg.FromUser,
		"_to":           i.SIP.SipMsg.ToUser,
		"_method":       i.SIP.SipMsg.StartLine.Method,
		"_statusCode":   i.SIP.SipMsg.StartLine.Resp,
		"_ua":           i.SIP.SipMsg.UserAgent,
		"_srcIP":        i.SIP.SrcIPString,
		"_dstIP":        i.SIP.DstIPString,
		"_srcPort":      i.SIP.SrcPort,
		"_dstPort":      i.SIP.DstPort,
		"_ipLen":        i.SIP.IPlen,
		"_udpLen":       i.SIP.UDPlen,
		"_intVlan":      i.SIP.IntVlan,
		"_vl":           i.SIP.VL,
		"_tos":          i.SIP.TOS,
		"_tlen":         i.SIP.TLen,
		"_tid":          i.SIP.TID,
		"_tflags":       i.SIP.TFlags,
		"_ttl":          i.SIP.TTL,
		"_tproto":       i.SIP.TProto,
	}
	return &m
}

func (i *IPFIX) mapLogQOS() *map[string]interface{} {
	m := map[string]interface{}{
		"version":             "1.1",
		"host":                *name,
		"short_message":       "LogQOS",
		"level":               5,
		"_incRtpBytes":        i.QOS.IncRtpBytes,
		"_incRtpPackets":      i.QOS.IncRtpPackets,
		"_incRtpLostPackets":  i.QOS.IncRtpLostPackets,
		"_incRtpAvgJitter":    i.QOS.IncRtpAvgJitter,
		"_incRtpMaxJitter":    i.QOS.IncRtpMaxJitter,
		"_incRtcpBytes":       i.QOS.IncRtcpBytes,
		"_incRtcpPackets":     i.QOS.IncRtcpPackets,
		"_incRtcpLostPackets": i.QOS.IncRtcpLostPackets,
		"_incRtcpAvgJitter":   i.QOS.IncRtcpAvgJitter,
		"_incRtcpMaxJitter":   i.QOS.IncRtcpMaxJitter,
		"_incRtcpAvgLat":      i.QOS.IncRtcpAvgLat,
		"_incRtcpMaxLat":      i.QOS.IncRtcpMaxLat,
		"_incrVal":            float64(i.QOS.IncrVal) / 100,
		"_incMos":             float64(i.QOS.IncMos) / 100,
		"_outRtpBytes":        i.QOS.OutRtpBytes,
		"_outRtpPackets":      i.QOS.OutRtpPackets,
		"_outRtpLostPackets":  i.QOS.OutRtpLostPackets,
		"_outRtpAvgJitter":    i.QOS.OutRtpAvgJitter,
		"_outRtpMaxJitter":    i.QOS.OutRtpMaxJitter,
		"_outRtcpBytes":       i.QOS.OutRtcpBytes,
		"_outRtcpPackets":     i.QOS.OutRtcpPackets,
		"_outRtcpLostPackets": i.QOS.OutRtcpLostPackets,
		"_outRtcpAvgJitter":   i.QOS.OutRtcpAvgJitter,
		"_outRtcpMaxJitter":   i.QOS.OutRtcpMaxJitter,
		"_outRtcpAvgLat":      i.QOS.OutRtcpAvgLat,
		"_outRtcpMaxLat":      i.QOS.OutRtcpMaxLat,
		"_outrVal":            float64(i.QOS.OutrVal) / 100,
		"_outMos":             float64(i.QOS.OutMos) / 100,
		"_type":               i.QOS.Type,
		"_callerIncSrcIP":     stringIPv4(i.QOS.CallerIncSrcIP),
		"_callerIncDstIP":     stringIPv4(i.QOS.CallerIncDstIP),
		"_callerIncSrcPort":   i.QOS.CallerIncSrcPort,
		"_callerIncDstPort":   i.QOS.CallerIncDstPort,
		"_calleeIncSrcIP":     stringIPv4(i.QOS.CalleeIncSrcIP),
		"_calleeIncDstIP":     stringIPv4(i.QOS.CalleeIncDstIP),
		"_calleeIncSrcPort":   i.QOS.CalleeIncSrcPort,
		"_calleeIncDstPort":   i.QOS.CalleeIncDstPort,
		"_callerOutSrcIP":     stringIPv4(i.QOS.CallerOutSrcIP),
		"_callerOutDstIP":     stringIPv4(i.QOS.CallerOutDstIP),
		"_callerOutSrcPort":   i.QOS.CallerOutSrcPort,
		"_callerOutDstPort":   i.QOS.CallerOutDstPort,
		"_calleeOutSrcIP":     stringIPv4(i.QOS.CalleeOutSrcIP),
		"_calleeOutDstIP":     stringIPv4(i.QOS.CalleeOutDstIP),
		"_calleeOutSrcPort":   i.QOS.CalleeOutSrcPort,
		"_calleeOutDstPort":   i.QOS.CalleeOutDstPort,
		"_callerIntSlot":      i.QOS.CallerIntSlot,
		"_callerIntPort":      i.QOS.CallerIntPort,
		"_callerIntVlan":      i.QOS.CallerIntVlan,
		"_calleeIntSlot":      i.QOS.CalleeIntSlot,
		"_calleeIntPort":      i.QOS.CalleeIntPort,
		"_calleeIntVlan":      i.QOS.CalleeIntVlan,
		"_beginTimeSec":       i.QOS.BeginTimeSec,
		"_beginTimeMic":       i.QOS.BeginTimeMic,
		"_endTimeSec":         i.QOS.EndTimeSec,
		"_endTimeMic":         i.QOS.EndTimeMic,
		"_duration":           int(i.QOS.EndTimeSec - i.QOS.BeginTimeSec),
		"_seperator":          i.QOS.Seperator,
		"_incRealm":           string(i.QOS.IncRealm),
		"_outRealm":           string(i.QOS.OutRealm),
		"_idLen":              i.QOS.IncCallIDLen,
		"_id":                 string(i.QOS.IncCallID),
		"_outCallIDLen":       i.QOS.OutCallIDLen,
		"_outCallID":          string(i.QOS.OutCallID),
	}
	return &m
}

func (i *IPFIX) mapMetricQOS() map[string]interface{} {
	m := map[string]interface{}{
		"inc.rtp.mos":          float64(i.QOS.IncMos) / 100,
		"out.rtp.mos":          float64(i.QOS.OutMos) / 100,
		"inc.rtp.rval":         float64(i.QOS.IncrVal) / 100,
		"out.rtp.rval":         float64(i.QOS.OutrVal) / 100,
		"inc.rtp.packets":      float64(i.QOS.IncRtpPackets),
		"out.rtp.packets":      float64(i.QOS.OutRtpPackets),
		"inc.rtcp.packets":     float64(i.QOS.IncRtcpPackets),
		"out.rtcp.packets":     float64(i.QOS.OutRtcpPackets),
		"inc.rtp.lostPackets":  float64(i.QOS.IncRtpLostPackets),
		"out.rtp.lostPackets":  float64(i.QOS.OutRtpLostPackets),
		"inc.rtcp.lostPackets": float64(i.QOS.IncRtcpLostPackets),
		"out.rtcp.lostPackets": float64(i.QOS.OutRtcpLostPackets),
		"inc.rtp.avgJitter":    float64(i.QOS.IncRtpAvgJitter),
		"out.rtp.avgJitter":    float64(i.QOS.OutRtpAvgJitter),
		"inc.rtp.maxJitter":    float64(i.QOS.IncRtpMaxJitter),
		"out.rtp.maxJitter":    float64(i.QOS.OutRtpMaxJitter),
		"inc.rtcp.avgJitter":   float64(i.QOS.IncRtcpAvgJitter),
		"out.rtcp.avgJitter":   float64(i.QOS.OutRtcpAvgJitter),
		"inc.rtcp.maxJitter":   float64(i.QOS.IncRtcpMaxJitter),
		"out.rtcp.maxJitter":   float64(i.QOS.OutRtcpMaxJitter),
		"inc.rtcp.avgLat":      float64(i.QOS.IncRtcpAvgLat),
		"out.rtcp.avgLat":      float64(i.QOS.OutRtcpAvgLat),
		"inc.rtcp.maxLat":      float64(i.QOS.IncRtcpMaxLat),
		"out.rtcp.maxLat":      float64(i.QOS.OutRtcpMaxLat),
	}
	return m
}

func (i *IPFIX) mapQOS() []byte {
	var b bytes.Buffer

	b.WriteString("{")
	b.WriteString("\"NAME\":\"")
	b.WriteString(*name)
	b.WriteString("\",\"INC_REALM\":\"")
	b.WriteString(string(i.QOS.IncRealm))
	b.WriteString("\",\"OUT_REALM\":\"")
	b.WriteString(string(i.QOS.OutRealm))
	b.WriteString("\",\"INC_ID\":\"")
	b.WriteString(string(i.QOS.IncCallID))
	b.WriteString("\",\"INC_RTP_BYTE\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtpBytes)))
	b.WriteString(",\"INC_RTP_PK\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtpPackets)))
	b.WriteString(",\"INC_RTP_PK_LOSS\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtpLostPackets)))
	b.WriteString(",\"INC_RTP_AVG_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtpAvgJitter)))
	b.WriteString(",\"INC_RTP_MAX_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtpMaxJitter)))
	b.WriteString(",\"INC_RTCP_BYTE\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpBytes)))
	b.WriteString(",\"INC_RTCP_PK\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpPackets)))
	b.WriteString(",\"INC_RTCP_PK_LOSS\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpLostPackets)))
	b.WriteString(",\"INC_RTCP_AVG_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpAvgJitter)))
	b.WriteString(",\"INC_RTCP_MAX_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpMaxJitter)))
	b.WriteString(",\"INC_RTCP_AVG_LAT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpAvgLat)))
	b.WriteString(",\"INC_RTCP_MAX_LAT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncRtcpMaxLat)))
	b.WriteString(",\"INC_MOS\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncMos)))
	b.WriteString(",\"INC_RVAL\":")
	b.WriteString(strconv.Itoa(int(i.QOS.IncrVal)))
	b.WriteString(",\"CALLER_VLAN\":")
	b.WriteString(strconv.Itoa(int(i.QOS.CallerIntVlan)))
	b.WriteString(",\"CALLER_SRC_IP\":\"")
	b.WriteString(stringIPv4(i.QOS.CallerIncSrcIP))
	b.WriteString("\",\"CALLER_SRC_PORT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.CallerIncSrcPort)))
	b.WriteString(",\"CALLER_DST_IP\":\"")
	b.WriteString(stringIPv4(i.QOS.CallerOutDstIP))
	b.WriteString("\",\"CALLER_DST_PORT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.CallerOutDstPort)))

	b.WriteString(",\"OUT_ID\":\"")
	b.WriteString(string(i.QOS.OutCallID))
	b.WriteString("\",\"OUT_RTP_BYTE\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtpBytes)))
	b.WriteString(",\"OUT_RTP_PK\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtpPackets)))
	b.WriteString(",\"OUT_RTP_PK_LOSS\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtpLostPackets)))
	b.WriteString(",\"OUT_RTP_AVG_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtpAvgJitter)))
	b.WriteString(",\"OUT_RTP_MAX_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtpMaxJitter)))
	b.WriteString(",\"OUT_RTCP_BYTE\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpBytes)))
	b.WriteString(",\"OUT_RTCP_PK\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpPackets)))
	b.WriteString(",\"OUT_RTCP_PK_LOSS\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpLostPackets)))
	b.WriteString(",\"OUT_RTCP_AVG_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpAvgJitter)))
	b.WriteString(",\"OUT_RTCP_MAX_JITTER\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpMaxJitter)))
	b.WriteString(",\"OUT_RTCP_AVG_LAT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpAvgLat)))
	b.WriteString(",\"OUT_RTCP_MAX_LAT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutRtcpMaxLat)))
	b.WriteString(",\"OUT_MOS\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutMos)))
	b.WriteString(",\"OUT_RVAL\":")
	b.WriteString(strconv.Itoa(int(i.QOS.OutrVal)))
	b.WriteString(",\"CALLEE_VLAN\":")
	b.WriteString(strconv.Itoa(int(i.QOS.CalleeIntVlan)))
	b.WriteString(",\"CALLEE_SRC_IP\":\"")
	b.WriteString(stringIPv4(i.QOS.CalleeOutSrcIP))
	b.WriteString("\",\"CALLEE_SRC_PORT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.CalleeOutSrcPort)))
	b.WriteString(",\"CALLEE_DST_IP\":\"")
	b.WriteString(stringIPv4(i.QOS.CalleeIncDstIP))
	b.WriteString("\",\"CALLEE_DST_PORT\":")
	b.WriteString(strconv.Itoa(int(i.QOS.CalleeIncDstPort)))
	b.WriteString(",\"MEDIA_TYPE\":")
	b.WriteString(strconv.Itoa(int(i.QOS.Type)))
	b.WriteString("}")

	return b.Bytes()
}

func (i *IPFIX) mapIncQOS() *map[string]interface{} {
	m := map[string]interface{}{

		"CORRELATION_ID":  string(i.QOS.IncCallID),
		"RTP_SIP_CALL_ID": string(i.QOS.IncCallID),
		"REPORT_TS":       i.QOS.BeginTimeSec,
		"TL_BYTE":         i.QOS.IncRtpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        i.QOS.IncRtpPackets,
		"EXPECTED_PK":     (i.QOS.IncRtpPackets + i.QOS.IncRtpLostPackets),
		"PACKET_LOSS":     i.QOS.IncRtpLostPackets,
		"SEQ":             0,
		"JITTER":          i.QOS.IncRtpAvgJitter,
		"MAX_JITTER":      i.QOS.IncRtpMaxJitter,
		"MEAN_JITTER":     i.QOS.IncRtpAvgJitter,
		"DELTA":           i.QOS.IncRtcpAvgLat,
		"MAX_DELTA":       i.QOS.IncRtcpMaxLat,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         i.QOS.IncMos,
		"MEAN_MOS":        i.QOS.IncMos,
		"MOS":             i.QOS.IncMos,
		"RFACTOR":         i.QOS.IncrVal,
		"MIN_RFACTOR":     i.QOS.IncrVal,
		"MEAN_RFACTOR":    i.QOS.IncrVal,
		"SRC_IP":          stringIPv4(i.QOS.CallerIncSrcIP),
		"SRC_PORT":        i.QOS.CallerIncSrcPort,
		"DST_IP":          stringIPv4(i.QOS.CallerIncDstIP),
		"DST_PORT":        i.QOS.CallerIncDstPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        i.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      i.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(i.QOS.IncRealm),
		"PARTY":           0,
		"TYPE":            "PERIODIC",
	}
	return &m
}

func (i *IPFIX) mapOutQOS() *map[string]interface{} {
	m := map[string]interface{}{

		"CORRELATION_ID":  string(i.QOS.OutCallID),
		"RTP_SIP_CALL_ID": string(i.QOS.OutCallID),
		"REPORT_TS":       i.QOS.BeginTimeSec,
		"TL_BYTE":         i.QOS.OutRtpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        i.QOS.OutRtpPackets,
		"EXPECTED_PK":     (i.QOS.OutRtpPackets + i.QOS.OutRtpLostPackets),
		"PACKET_LOSS":     i.QOS.OutRtpLostPackets,
		"SEQ":             0,
		"JITTER":          i.QOS.OutRtpAvgJitter,
		"MAX_JITTER":      i.QOS.OutRtpMaxJitter,
		"MEAN_JITTER":     i.QOS.OutRtpAvgJitter,
		"DELTA":           i.QOS.OutRtcpAvgLat,
		"MAX_DELTA":       i.QOS.OutRtcpMaxLat,
		"MAX_SKEW":        0.000,
		"MIN_MOS":         i.QOS.OutMos,
		"MEAN_MOS":        i.QOS.OutMos,
		"MOS":             i.QOS.OutMos,
		"RFACTOR":         i.QOS.OutrVal,
		"MIN_RFACTOR":     i.QOS.OutrVal,
		"MEAN_RFACTOR":    i.QOS.OutrVal,
		"SRC_IP":          stringIPv4(i.QOS.CalleeOutSrcIP),
		"SRC_PORT":        i.QOS.CalleeOutSrcPort,
		"DST_IP":          stringIPv4(i.QOS.CalleeOutDstIP),
		"DST_PORT":        i.QOS.CalleeOutDstPort,
		"SRC_MAC":         "00-00-00-00-00-00",
		"DST_MAC":         "00-00-00-00-00-00",
		"OUT_ORDER":       0,
		"SSRC_CHG":        0,
		"CODEC_PT":        i.QOS.Type,
		"CLOCK":           8000,
		"CODEC_NAME":      i.QOS.Type,
		"DIR":             0,
		"REPORT_NAME":     string(i.QOS.OutRealm),
		"PARTY":           1,
		"TYPE":            "PERIODIC",
	}
	return &m
}
