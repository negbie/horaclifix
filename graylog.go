package main

import (
	"encoding/json"
	"net"
)

func (s *ByteString) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(string(*s))
	return bytes, err
}

func (s *ByteString) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, &x)
	*s = ByteString(x)
	return err
}

func (ipfix *IPFIX) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(*ipfix)
	return bytes, err
}

func LogSIP(i *IPFIX) {
	gconn, err := net.Dial("udp", *gaddr)
	checkErr(err)
	defer gconn.Close()
	//sLog, _ := json.Marshal(&i.Data.SIP)
	//gconn.Write(sLog)

	json.NewEncoder(gconn).Encode(&i.Data.SIP)
}

func LogQOS(i *IPFIX) {
	gconn, err := net.Dial("udp", *gaddr)
	checkErr(err)
	defer gconn.Close()
	//qLog, _ := json.Marshal(&i.Data.QOS)
	//gconn.Write(qLog)

	json.NewEncoder(gconn).Encode(&i.Data.QOS)
}

/*
func LogSIP(ipfix *IPFIX) {
	gconn, err := net.Dial("udp", *gaddr)
	checkErr(err)
	defer gconn.Close()
	mapSIP := map[string]interface{}{

		"timeSec": ipfix.Data.SIP.TimeSec,
		"timeMic": ipfix.Data.SIP.TimeMic,
		"intVlan": ipfix.Data.SIP.IntVlan,
		"id":      string(ipfix.Data.SIP.CallID),
		"ipLen":   ipfix.Data.SIP.IPlen,
		"udpLen":  ipfix.Data.SIP.UDPlen,
		"vl":      ipfix.Data.SIP.VL,
		"tos":     ipfix.Data.SIP.TOS,
		"tlen":    ipfix.Data.SIP.TLen,
		"tid":     ipfix.Data.SIP.TID,
		"tflags":  ipfix.Data.SIP.TFlags,
		"ttl":     ipfix.Data.SIP.TTL,
		"tproto":  ipfix.Data.SIP.TProto,
		"srcIp":   ipfix.Data.SIP.SrcIP,
		"dstIp":   ipfix.Data.SIP.DstIP,
		"srcPort": ipfix.Data.SIP.SrcPort,
		"dstPort": ipfix.Data.SIP.DstPort,
		"context": ipfix.Data.SIP.Context,
		"sipMsg":  string(ipfix.Data.SIP.SipMsg),
		"sbc":     "sbcSIP",
	}

	sLog, _ := json.Marshal(mapSIP)

	gconn.Write(sLog)

	if *debug {
		fmt.Println("Json output:")
		fmt.Printf("%s\n\n\n", sLog)
	}
}
*/
