package main

import (
	"fmt"

	"github.com/negbie/siprocket"
)

// mapLogSIP retruns a map with SIP stats which can be
// json encoded and send send as gelf to graylog
func (i *IPFIX) mapLogSIP() *map[string]string {
	sipMSG := siprocket.Parse(i.Data.SIP.SipMsg)
	//siprocket.PrintSipStruct(&sipMSG)
	mLogSIP := map[string]string{
		"version":       "1.1",
		"host":          *name,
		"short_message": string(i.Data.SIP.SipMsg),
		"level":         "5",
		"_id":           string(sipMSG.CallId.Value),
		"_from":         string(sipMSG.From.User),
		"_to":           string(sipMSG.To.User),
		"_method":       string(sipMSG.Req.Method),
		"_statusCode":   string(sipMSG.Req.StatusCode),
		"_ua":           string(sipMSG.Ua.Value),
		"_srcIp":        stringIPv4(i.Data.SIP.SrcIP),
		"_dstIp":        stringIPv4(i.Data.SIP.DstIP),
		"_srcPort":      fmt.Sprint(i.Data.SIP.SrcPort),
		"_dstPort":      fmt.Sprint(i.Data.SIP.DstPort),
		"_ipLen":        fmt.Sprint(i.Data.SIP.IPlen),
		"_udpLen":       fmt.Sprint(i.Data.SIP.UDPlen),
		"_intVlan":      fmt.Sprint(i.Data.SIP.IntVlan),
		"_vl":           fmt.Sprint(i.Data.SIP.VL),
		"_tos":          fmt.Sprint(i.Data.SIP.TOS),
		"_tlen":         fmt.Sprint(i.Data.SIP.TLen),
		"_tid":          fmt.Sprint(i.Data.SIP.TID),
		"_tflags":       fmt.Sprint(i.Data.SIP.TFlags),
		"_ttl":          fmt.Sprint(i.Data.SIP.TTL),
		"_tproto":       fmt.Sprint(i.Data.SIP.TProto),
	}
	return &mLogSIP
}

// mapLogQOS retruns a map with QOS stats which can be
// json encoded and send as gelf to graylog
func (i *IPFIX) mapLogQOS() *map[string]interface{} {
	mLogQOS := map[string]interface{}{
		"version":             "1.1",
		"host":                *name,
		"short_message":       "LogQOS",
		"level":               5,
		"_incRtpBytes":        i.Data.QOS.IncRtpBytes,
		"_incRtpPackets":      i.Data.QOS.IncRtpPackets,
		"_incRtpLostPackets":  i.Data.QOS.IncRtpLostPackets,
		"_incRtpAvgJitter":    i.Data.QOS.IncRtpAvgJitter,
		"_incRtpMaxJitter":    i.Data.QOS.IncRtpMaxJitter,
		"_incRtcpBytes":       i.Data.QOS.IncRtcpBytes,
		"_incRtcpPackets":     i.Data.QOS.IncRtcpPackets,
		"_incRtcpLostPackets": i.Data.QOS.IncRtcpLostPackets,
		"_incRtcpAvgJitter":   i.Data.QOS.IncRtcpAvgJitter,
		"_incRtcpMaxJitter":   i.Data.QOS.IncRtcpMaxJitter,
		"_incRtcpAvgLat":      i.Data.QOS.IncRtcpAvgLat,
		"_incRtcpMaxLat":      i.Data.QOS.IncRtcpMaxLat,
		"_incrVal":            i.Data.QOS.IncrVal,
		"_incMos":             i.Data.QOS.IncMos,
		"_outRtpBytes":        i.Data.QOS.OutRtpBytes,
		"_outRtpPackets":      i.Data.QOS.OutRtpPackets,
		"_outRtpLostPackets":  i.Data.QOS.OutRtpLostPackets,
		"_outRtpAvgJitter":    i.Data.QOS.OutRtpAvgJitter,
		"_outRtpMaxJitter":    i.Data.QOS.OutRtpMaxJitter,
		"_outRtcpBytes":       i.Data.QOS.OutRtcpBytes,
		"_outRtcpPackets":     i.Data.QOS.OutRtcpPackets,
		"_outRtcpLostPackets": i.Data.QOS.OutRtcpLostPackets,
		"_outRtcpAvgJitter":   i.Data.QOS.OutRtcpAvgJitter,
		"_outRtcpMaxJitter":   i.Data.QOS.OutRtcpMaxJitter,
		"_outRtcpAvgLat":      i.Data.QOS.OutRtcpAvgLat,
		"_outRtcpMaxLat":      i.Data.QOS.OutRtcpMaxLat,
		"_outrVal":            i.Data.QOS.OutrVal,
		"_outMos":             i.Data.QOS.OutMos,
		"_type":               i.Data.QOS.Type,
		"_callerIncSrcIP":     stringIPv4(i.Data.QOS.CallerIncSrcIP),
		"_callerIncDstIP":     stringIPv4(i.Data.QOS.CallerIncDstIP),
		"_callerIncSrcPort":   i.Data.QOS.CallerIncSrcPort,
		"_callerIncDstPort":   i.Data.QOS.CallerIncDstPort,
		"_calleeIncSrcIP":     stringIPv4(i.Data.QOS.CalleeIncSrcIP),
		"_calleeIncDstIP":     stringIPv4(i.Data.QOS.CalleeIncDstIP),
		"_calleeIncSrcPort":   i.Data.QOS.CalleeIncSrcPort,
		"_calleeIncDstPort":   i.Data.QOS.CalleeIncDstPort,
		"_callerOutSrcIP":     stringIPv4(i.Data.QOS.CallerOutSrcIP),
		"_callerOutDstIP":     stringIPv4(i.Data.QOS.CallerOutDstIP),
		"_callerOutSrcPort":   i.Data.QOS.CallerOutSrcPort,
		"_callerOutDstPort":   i.Data.QOS.CallerOutDstPort,
		"_calleeOutSrcIP":     stringIPv4(i.Data.QOS.CalleeOutSrcIP),
		"_calleeOutDstIP":     stringIPv4(i.Data.QOS.CalleeOutDstIP),
		"_calleeOutSrcPort":   i.Data.QOS.CalleeOutSrcPort,
		"_calleeOutDstPort":   i.Data.QOS.CalleeOutDstPort,
		"_callerIntSlot":      i.Data.QOS.CallerIntSlot,
		"_callerIntPort":      i.Data.QOS.CallerIntPort,
		"_callerIntVlan":      i.Data.QOS.CallerIntVlan,
		"_calleeIntSlot":      i.Data.QOS.CalleeIntSlot,
		"_calleeIntPort":      i.Data.QOS.CalleeIntPort,
		"_calleeIntVlan":      i.Data.QOS.CalleeIntVlan,
		"_beginTimeSec":       i.Data.QOS.BeginTimeSec,
		"_beginTimeMic":       i.Data.QOS.BeginTimeMic,
		"_endTimeSec":         i.Data.QOS.EndTimeSec,
		"_endinTimeMic":       i.Data.QOS.EndinTimeMic,
		"_duration":           (i.Data.QOS.EndTimeSec - i.Data.QOS.BeginTimeSec),
		"_seperator":          i.Data.QOS.Seperator,
		"_incRealm":           string(i.Data.QOS.IncRealm),
		"_outRealm":           string(i.Data.QOS.OutRealm),
		"_idLen":              i.Data.QOS.IncCallIDLen,
		"_id":                 string(i.Data.QOS.IncCallID),
		"_outCallIDLen":       i.Data.QOS.OutCallIDLen,
		"_outCallID":          string(i.Data.QOS.OutCallID),
	}
	return &mLogQOS
}

// mapAllQOS retruns a map with QOS33 stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) mapAllQOS() *map[string]interface{} {
	mAllQOS := map[string]interface{}{
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
	return &mAllQOS
}

// mapIncQOS retruns a map with incomming RTP QOS stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) mapIncQOS() *map[string]interface{} {
	mIncQOS := map[string]interface{}{

		"CORRELATION_ID":  string(i.Data.QOS.IncCallID),
		"RTP_SIP_CALL_ID": string(i.Data.QOS.IncCallID),
		"REPORT_TS":       i.Data.QOS.BeginTimeSec,
		"TL_BYTE":         i.Data.QOS.IncRtpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        i.Data.QOS.IncRtpPackets,
		"EXPECTED_PK":     (i.Data.QOS.IncRtpPackets + i.Data.QOS.IncRtpLostPackets),
		"PACKET_LOSS":     i.Data.QOS.IncRtpLostPackets,
		"SEQ":             0,
		"JITTER":          i.Data.QOS.IncRtpAvgJitter,
		"MAX_JITTER":      i.Data.QOS.IncRtpMaxJitter,
		"MEAN_JITTER":     i.Data.QOS.IncRtpAvgJitter,
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
	}
	return &mIncQOS
}

// mapOutQOS retruns a map with outgoing RTP QOS stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) mapOutQOS() *map[string]interface{} {
	mOutQOS := map[string]interface{}{

		"CORRELATION_ID":  string(i.Data.QOS.OutCallID),
		"RTP_SIP_CALL_ID": string(i.Data.QOS.OutCallID),
		"REPORT_TS":       i.Data.QOS.BeginTimeSec,
		"TL_BYTE":         i.Data.QOS.OutRtpBytes,
		"SKEW":            0.000,
		"TOTAL_PK":        i.Data.QOS.OutRtpPackets,
		"EXPECTED_PK":     (i.Data.QOS.OutRtpPackets + i.Data.QOS.OutRtpLostPackets),
		"PACKET_LOSS":     i.Data.QOS.OutRtpLostPackets,
		"SEQ":             0,
		"JITTER":          i.Data.QOS.OutRtpAvgJitter,
		"MAX_JITTER":      i.Data.QOS.OutRtpMaxJitter,
		"MEAN_JITTER":     i.Data.QOS.OutRtpAvgJitter,
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
	}
	return &mOutQOS
}

/*
// mapIncQOS retruns a map with incomming RTCP QOS stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) mapIncQOS() *map[string]interface{} {
	mIncQOS := map[string]interface{}{

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
	}
	return &mIncQOS
}

// mapOutQOS retruns a map with outgoing RTCP QOS stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) mapOutQOS() *map[string]interface{} {
	mOutQOS := map[string]interface{}{
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
	}
	return &mOutQOS
}
*/
