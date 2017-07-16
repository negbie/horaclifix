package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

/*
func (s *ByteString) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(string(*s))
	return bytes, err
}

func (i *ByteIP) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.IP(*i).String())
}

func (s *ByteString) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, &x)
	*s = ByteString(x)
	return err
}

func (i *IPFIX) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(*i)
	return bytes, err
}
*/

// SendLog will encode the QOS & SIP maps to json
// and send them over UDP to Graylog
func (conn Connections) SendLog(i *IPFIX, s string) {
	var gLog []byte
	var err error
	switch s {
	case "SIP":

		if gLog, err = json.Marshal(i.PrepLogSIP()); err != nil {
			log.Println("SIP json.NewEncoder failed:", err, gLog)
		}

		/*		if err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogSIP()); err != nil {
				log.Println("SIP json.NewEncoder failed:", err, gLog)
			}*/
	case "QOS":

		if gLog, err = json.Marshal(i.PrepLogQOS()); err != nil {
			log.Println("SIP json.NewEncoder failed:", err, gLog)
		}

		/*		if err := json.NewEncoder(conn.Graylog).Encode(i.PrepLogQOS()); err != nil {
				log.Println("SIP json.NewEncoder failed:", err, gLog)
			}*/
	}
	if _, err := conn.Graylog.Write(gLog); err != nil {
		checkErr(err)
	}
}

// PrepLogSIP retruns a map with SIP stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) PrepLogSIP() *map[string]interface{} {
	mapSIP := map[string]interface{}{
		"intVlan": i.Data.SIP.IntVlan,
		"id":      string(i.Data.SIP.CallID),
		"ipLen":   i.Data.SIP.IPlen,
		"udpLen":  i.Data.SIP.UDPlen,
		"vl":      i.Data.SIP.VL,
		"tos":     i.Data.SIP.TOS,
		"tlen":    i.Data.SIP.TLen,
		"tid":     i.Data.SIP.TID,
		"tflags":  i.Data.SIP.TFlags,
		"ttl":     i.Data.SIP.TTL,
		"tproto":  i.Data.SIP.TProto,
		"srcIp":   stringIPv4(i.Data.SIP.SrcIP),
		"dstIp":   stringIPv4(i.Data.SIP.DstIP),
		"srcPort": i.Data.SIP.SrcPort,
		"dstPort": i.Data.SIP.DstPort,
		"sipMsg":  string(i.Data.SIP.SipMsg),
		"sbcName": *name,
	}
	return &mapSIP
}

// PrepLogQOS retruns a map with QOS stats which can be
// json encoded and send into homer, or graylog
func (i *IPFIX) PrepLogQOS() *map[string]interface{} {
	mapQOS := map[string]interface{}{

		"incRtpBytes":       i.Data.QOS.IncRtpBytes,
		"incRtpPackets":     i.Data.QOS.IncRtpPackets,
		"incRtpLostPackets": i.Data.QOS.IncRtpLostPackets,
		"incRtpAvgJitter":   i.Data.QOS.IncRtpAvgJitter,
		"incRtpMaxJitter":   i.Data.QOS.IncRtpMaxJitter,

		"incRtcpBytes":       i.Data.QOS.IncRtcpBytes,
		"incRtcpPackets":     i.Data.QOS.IncRtcpPackets,
		"incRtcpLostPackets": i.Data.QOS.IncRtcpLostPackets,
		"incRtcpAvgJitter":   i.Data.QOS.IncRtcpAvgJitter,
		"incRtcpMaxJitter":   i.Data.QOS.IncRtcpMaxJitter,
		"incRtcpAvgLat":      i.Data.QOS.IncRtcpAvgLat,
		"incRtcpMaxLat":      i.Data.QOS.IncRtcpMaxLat,

		"incrVal": i.Data.QOS.IncrVal,
		"incMos":  i.Data.QOS.IncMos,

		"outRtpBytes":       i.Data.QOS.OutRtpBytes,
		"outRtpPackets":     i.Data.QOS.OutRtpPackets,
		"outRtpLostPackets": i.Data.QOS.OutRtpLostPackets,
		"outRtpAvgJitter":   i.Data.QOS.OutRtpAvgJitter,
		"outRtpMaxJitter":   i.Data.QOS.OutRtpMaxJitter,

		"outRtcpBytes":       i.Data.QOS.OutRtcpBytes,
		"outRtcpPackets":     i.Data.QOS.OutRtcpPackets,
		"outRtcpLostPackets": i.Data.QOS.OutRtcpLostPackets,
		"outRtcpAvgJitter":   i.Data.QOS.OutRtcpAvgJitter,
		"outRtcpMaxJitter":   i.Data.QOS.OutRtcpMaxJitter,
		"outRtcpAvgLat":      i.Data.QOS.OutRtcpAvgLat,
		"outRtcpMaxLat":      i.Data.QOS.OutRtcpMaxLat,

		"outrVal": i.Data.QOS.OutrVal,
		"outMos":  i.Data.QOS.OutMos,

		"type": i.Data.QOS.Type,

		"callerIncSrcIP":   stringIPv4(i.Data.QOS.CallerIncSrcIP),
		"callerIncDstIP":   stringIPv4(i.Data.QOS.CallerIncDstIP),
		"callerIncSrcPort": i.Data.QOS.CallerIncSrcPort,
		"callerIncDstPort": i.Data.QOS.CallerIncDstPort,

		"calleeIncSrcIP":   stringIPv4(i.Data.QOS.CalleeIncSrcIP),
		"calleeIncDstIP":   stringIPv4(i.Data.QOS.CalleeIncDstIP),
		"calleeIncSrcPort": i.Data.QOS.CalleeIncSrcPort,
		"calleeIncDstPort": i.Data.QOS.CalleeIncDstPort,

		"callerOutSrcIP":   stringIPv4(i.Data.QOS.CallerOutSrcIP),
		"callerOutDstIP":   stringIPv4(i.Data.QOS.CallerOutDstIP),
		"callerOutSrcPort": i.Data.QOS.CallerOutSrcPort,
		"callerOutDstPort": i.Data.QOS.CallerOutDstPort,

		"calleeOutSrcIP":   stringIPv4(i.Data.QOS.CalleeOutSrcIP),
		"calleeOutDstIP":   stringIPv4(i.Data.QOS.CalleeOutDstIP),
		"calleeOutSrcPort": i.Data.QOS.CalleeOutSrcPort,
		"calleeOutDstPort": i.Data.QOS.CalleeOutDstPort,

		"callerIntSlot": i.Data.QOS.CallerIntSlot,
		"callerIntPort": i.Data.QOS.CallerIntPort,
		"callerIntVlan": i.Data.QOS.CallerIntVlan,

		"calleeIntSlot": i.Data.QOS.CalleeIntSlot,
		"calleeIntPort": i.Data.QOS.CalleeIntPort,
		"calleeIntVlan": i.Data.QOS.CalleeIntVlan,

		"beginTimeSec": i.Data.QOS.BeginTimeSec,
		"beginTimeMic": i.Data.QOS.BeginTimeMic,

		"endTimeSec":   i.Data.QOS.EndTimeSec,
		"endinTimeMic": i.Data.QOS.EndinTimeMic,

		"duration": (i.Data.QOS.EndTimeSec - i.Data.QOS.BeginTimeSec),

		"seperator": i.Data.QOS.Seperator,

		"incRealmLen": i.Data.QOS.IncRealmLen,
		"incRealm":    string(i.Data.QOS.IncRealm),
		"incRealmEnd": i.Data.QOS.IncRealmEnd,

		"outRealmLen": i.Data.QOS.OutRealmLen,
		"outRealm":    string(i.Data.QOS.OutRealm),
		"outRealmEnd": i.Data.QOS.OutRealmEnd,

		"incCallIDLen": i.Data.QOS.IncCallIDLen,
		"incCallID":    string(i.Data.QOS.IncCallID),
		"incCallIDEnd": i.Data.QOS.IncCallIDEnd,

		"outCallIDLen": i.Data.QOS.OutCallIDLen,
		"outCallID":    string(i.Data.QOS.OutCallID),
		"sbcName":      *name,
	}
	return &mapQOS
}

// stringIPv4 converts a ipv4 unit32 into a string
func stringIPv4(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}
