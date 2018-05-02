package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"unsafe"
)

type reader struct {
	r   io.Reader
	err error
}

func checkErr(err error) {
	if err != nil {
		log.Printf("[WARN] <%s> %s\n", *name, err)
	}
}

func checkCritErr(err error) {
	if err != nil {
		log.Fatalf("[CRIT] <%s> %s\n", *name, err)
	}
}

func (r *reader) binRead(data interface{}) {
	r.err = binary.Read(r.r, binary.BigEndian, data)
	checkErr(r.err)
}

// stringIPv4 converts a ipv4 unit32 into a string
func stringIPv4(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}

func unsafeBytesToStr(z []byte) string {
	return *(*string)(unsafe.Pointer(&z))
}

func normMaxQOS(val uint32) uint32 {
	if val > 10000000 {
		return 0
	}
	return val
}
