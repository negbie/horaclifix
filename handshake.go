package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

//SendHandshake writes the binary Handshake representation into the buffer
func SendHandshake(c *net.TCPConn, hs []byte) {

	var ipfix IPFIX
	r := bytes.NewReader(hs)
	err := binary.Read(r, binary.BigEndian, &ipfix.Header)
	err = binary.Read(r, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Read(r, binary.BigEndian, &ipfix.Data.HandShake)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}
	ipfix.SetHeader.ID++
	ipfix.Data.HandShake.Timeout = 0

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &ipfix.Header)
	err = binary.Write(buf, binary.BigEndian, &ipfix.SetHeader)
	err = binary.Write(buf, binary.BigEndian, &ipfix.Data.HandShake)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}

	b := buf.Bytes()
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}

	err = binary.Write(c, binary.BigEndian, bi)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
}
