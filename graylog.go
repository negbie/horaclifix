package main

import (
	"encoding/json"
	"log"
	"net"
)

func (s *ByteString) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(string(*s))
	return bytes, err
}

/*func (i *ByteIP) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.IP(*i).String())
}*/

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

	err = json.NewEncoder(gconn).Encode(&i.Data.SIP)
	if err != nil {
		log.Println("LogSIP json.NewEncoder failed:", err)
	}
}

func LogQOS(i *IPFIX) {
	gconn, err := net.Dial("udp", *gaddr)
	checkErr(err)
	defer gconn.Close()
	//qLog, _ := json.Marshal(&i.Data.QOS)
	//gconn.Write(qLog)

	err = json.NewEncoder(gconn).Encode(&i.Data.QOS)
	if err != nil {
		log.Println("LogQOS json.NewEncoder failed:", err)
	}
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
		log.Println("Json output:")
		log.Printf("%s\n\n\n", sLog)
	}
}
*/
