package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/Graylog2/go-gelf/gelf"
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

func NewGelfLogger() {
	if *graylogAddr != "" {
		gelfWriter, err := gelf.NewWriter(*graylogAddr)
		if err != nil {
			log.Fatalf("gelf.NewWriter: %s", err)
		}
		// log to both stderr and graylog2
		//log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
		//log.Printf("logging to stderr & graylog2@'%s'", *graylogAddr)

		log.SetOutput(io.MultiWriter(gelfWriter))
		// Sine we use json extractor in graylog we dont wont to prefix the log output with timestamps
		log.SetFlags(0)
	}
}

func LogSip(i *IPFIX) {

	s := &i.Data.SIP
	//json.Unmarshal([]byte(`{}`), s)

	sLog, _ := json.Marshal(s)
	log.Printf("%s\n", sLog)
	//fmt.Printf("%s\n", sLog)

}

func LogQos(i *IPFIX) {

	q := &i.Data.QOS
	//json.Unmarshal([]byte(`{}`), q)

	qLog, _ := json.Marshal(q)
	log.Printf("%s\n", qLog)
	//fmt.Printf("%s\n", qLog)

}
