package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

// SendHandshake will write the binary Handshake representation
// into the buffer and send it to wire
func SendHandshake(c *net.TCPConn, hs []byte) {
	var i IPFIX
	r := bytes.NewReader(hs)
	binary.Read(r, binary.BigEndian, &i.Header)
	binary.Read(r, binary.BigEndian, &i.SetHeader)
	binary.Read(r, binary.BigEndian, &i.Data.Hs)

	// increment setID to 257
	i.SetHeader.ID++
	// disable timeout
	i.Data.Hs.Timeout = 0

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &i.Header)
	binary.Write(buf, binary.BigEndian, &i.SetHeader)
	binary.Write(buf, binary.BigEndian, &i.Data.Hs)

	// Convert []byte into []int8
	b := buf.Bytes()
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}

	log.Printf("Send handshake message %v to %s at %v\n", bi, *name, c.RemoteAddr())
	err := binary.Write(c, binary.BigEndian, bi)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
}
