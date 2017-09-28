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
	h := &HandShake{}
	r := reader{r: bytes.NewReader(hs)}
	r.binRead(h)

	// increment setID to 257
	h.IpfixSetHeader.ID++

	// disable timeout
	h.Timeout = 0

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, h)
	checkErr(err)

	// Convert []byte into []int8
	b := buf.Bytes()
	bi := make([]int8, len(b))
	for i, v := range b {
		bi[i] = int8(v)
	}

	log.Printf("Send handshake message %v to %s at %v\n", bi, *name, c.RemoteAddr())
	err = binary.Write(c, binary.BigEndian, bi)
	checkErr(err)
}
