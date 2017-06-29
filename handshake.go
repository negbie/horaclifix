package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

//SendHandshake writes the binary Handshake representation into the buffer
func SendHandshake(c *net.TCPConn, hs []byte) {

	var ipfix IPFIX
	r := bytes.NewReader(hs)
	binary.Read(r, binary.BigEndian, &ipfix.Header)
	binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake)
	ipfix.SetHeader.ID++
	ipfix.Data.HandShake.Timeout = 0

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &ipfix.Header)
	binary.Write(buf, binary.BigEndian, &ipfix.SetHeader)
	binary.Write(buf, binary.BigEndian, &ipfix.Data.HandShake)

	b := buf.Bytes()
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}

	binary.Write(c, binary.BigEndian, bi)
}
