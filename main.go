package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	/*
		"net"
		_ "net/http/pprof"
	*/
	"os"
)

var (
	addr    = flag.String("l", ":4739", "IPFIX listen address")
	haddr   = flag.String("H", "", "Homer server address")
	saddr   = flag.String("s", "", "StatsD server address")
	debug   = flag.Bool("d", false, "Debug output to stdout")
	verbose = flag.Bool("v", false, "Verbose output to stdout")
	gaddr   = flag.String("g", "", "Graylog server address")
	name    = flag.String("n", "sbc", "SBC name in graylog")
	hepPw   = flag.String("P", "myhep", "HEP capture password")
)

const headerLen int = 20

type Connections struct {
	Graylog net.Conn
	Homer   net.Conn
	StatsD  net.Conn
}

func newUDPConnections() *Connections {
	conn := new(Connections)
	if *gaddr != "" {
		gconn, err := net.Dial("udp", *gaddr)
		checkCritErr(err)
		conn.Graylog = gconn
	}
	if *haddr != "" {
		hconn, err := net.Dial("udp", *haddr)
		checkCritErr(err)
		conn.Homer = hconn
	}
	if *saddr != "" {
		sconn, err := net.Dial("udp", *saddr)
		checkCritErr(err)
		conn.StatsD = sconn
	}
	return conn
}

func checkErr(err error) {
	if err != nil {
		log.Println("Warning:", err)
	}
}

func checkCritErr(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}

/*
// Template for a sync.Pool buffer
var buffers = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 65536)
	},
}
packet := buffers.Get().([]byte)
buffers.Put(packet)
*/

// Start handles the incoming packets
func start(conn *net.TCPConn) {
	log.Printf("Handling new connection for %v|%v\n", *addr, *name)
	uConn := newUDPConnections()

	// Close connection when this function ends
	defer func() {
		log.Printf("Closing connection for %v|%v\n", *addr, *name)
		uConn.Graylog.Close()
		uConn.Homer.Close()
		uConn.StatsD.Close()
		conn.Close()
	}()

	// Create a buffer for incoming packets
	r := bufio.NewReader(conn)
	header := make([]byte, headerLen)
	var (
		version int
		dataLen int
		setID   int
		dataSet []byte
	)

	for {
		// Read the next 20 bytes which represent
		// the ipfix header + setheader at once
		if _, err := io.ReadFull(r, header); err == nil {
			// Extract the ipfix version
			version = int(uint16(header[0])<<8 + uint16(header[1]))
			// Extract the length of the whole packet
			dataLen = int(uint16(header[2])<<8 + uint16(header[3]))
			// Extract the setID
			setID = int(uint16(header[16])<<8 + uint16(header[17]))

			if setID > 280 || setID < 256 || version != 10 {
				break
			}
			// Create a buffer which holds exactly one
			// dataSet of the length of (dataLen-headerLen)
			// or in words the packet length minus the header length
			dataSet = make([]byte, dataLen-headerLen)

		} else if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s", err)
			}
			break
		}
		// Read the next (dataLen-headerLen) bytes into
		// the dataSet buffer at once
		if _, err := io.ReadFull(r, dataSet); err == nil {
			// Now merge the header with the dataSet.
			data := append(header, dataSet...)

			if *verbose {
				fmt.Println("########################################################################")
				fmt.Printf("Headerversion: %d, Headerlength: %d, SetID: %d\n", version, dataLen, setID)
				fmt.Println("Header in raw:", header)
			}
			if *debug {
				fmt.Println("Hexdump output:")
				fmt.Printf("%s\n", hex.Dump(data))
			}
			switch setID {
			case 256:
				SendHandshake(conn, data)
			case 258:
				msg := NewRecSipUDP(data)
				if *haddr != "" {
					uConn.SendHep(msg, "SIP")
				}

				if *gaddr != "" {
					uConn.SendLog(msg, "SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 259:
				msg := NewSendSipUDP(data)
				if *haddr != "" {
					uConn.SendHep(msg, "SIP")
				}

				if *gaddr != "" {
					uConn.SendLog(msg, "SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 260:
				msg := NewRecSipTCP(data)
				if *haddr != "" {
					uConn.SendHep(msg, "SIP")
				}

				if *gaddr != "" {
					uConn.SendLog(msg, "SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 261:
				msg := NewSendSipTCP(data)
				if *haddr != "" {
					uConn.SendHep(msg, "SIP")
				}

				if *gaddr != "" {
					uConn.SendLog(msg, "SIP")
				}
				if *debug {
					fmt.Println("SIP output:")
					fmt.Printf("%s\n", msg.Data.SIP.SipMsg)
				}
			case 268:
				msg := NewQosStats(data)
				if *haddr != "" {
					// Send only QOS stats with meaningful values
					if msg.Data.QOS.IncMos > 0 && msg.Data.QOS.OutMos > 0 {
						uConn.SendHep(msg, "QOS")
						//uConn.SendHep(msg, "logQOS")
					}
				}

				if *saddr != "" {
					// Send only QOS stats with meaningful values
					if msg.Data.QOS.IncMos > 0 && msg.Data.QOS.OutMos > 0 {
						uConn.SendStatsD(msg, "QOS")
					}
				}
				if *gaddr != "" {
					uConn.SendLog(msg, "QOS")
				}
			default:
				log.Printf("Unhandled setID %v\n", setID)
			}

		} else if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s", err)
			}
			break
		}
	}
}

func main() {
	//go http.ListenAndServe(":8080", http.DefaultServeMux)
	flag.Parse()
	f, err := os.OpenFile("horaclifix.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	checkCritErr(err)
	defer f.Close()
	log.SetOutput(f)

	log.Printf("Start horaclifix interfaces IPFIX:%v Homer:%v Graylog:%v\n", *addr, *haddr, *gaddr)

	laddr, err := net.ResolveTCPAddr("tcp", *addr)
	checkCritErr(err)

	listener, err := net.ListenTCP("tcp", laddr)
	checkCritErr(err)

	for {
		conn, err := listener.AcceptTCP()
		checkCritErr(err)
		go start(conn)
	}

}
