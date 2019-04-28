This project is deprecated.

# horaclifix
horaclifix sends IPFIX messages from Oracle SBC's into Homer


### Install:

Get it from the releases:
* https://github.com/negbie/horaclifix/releases
* chmod +x horaclifix
* ./horaclifix -H 192.168.2.22:9060&

Or if you have go installed:
* go install github.com/negbie/horaclifix
* look inside your go bin folder


### Stats:
![netdata_qos](https://user-images.githubusercontent.com/20154956/28118829-01909016-6713-11e7-9b54-80e626af7222.jpeg)

### Usage of ./horaclifix:

```bash
  -H string
    	Homer capture server address
  -I string
    	InfluxDB HTTP server address
  -N int
    	HEP capture node ID (default 2004)
  -P string
    	HEP capture password (default "myhep")
  -di string
    	Discard SIP method
  -g string
    	Graylog gelf TCP server address
  -l string
    	IPFIX TCP listen address (default ":4739")
  -m string
    	MySQL TCP server address
  -mp string
    	MySQL password
  -mu string
    	MySQL user
  -n string
    	SBC name (default "sbc")
  -nt string
    	Network types are [udp, tcp, tls] (default "udp")
  -p string
    	Prometheus address
  -protobuf
    	Use Protobuf on wire
  -s string
    	StatsD UDP server address

  -v    Verbose output to stdout
  -d    Debug output to stdout
  -V    Show version

  
################################################################
./horaclifix -h
./horaclifix -d -v
./horaclifix -H 192.168.2.22:9060 &
./horaclifix -H 192.168.2.22:9060 -nt tls &
./horaclifix -H 192.168.2.22:9060 -n labsbc -p 192.168.2.22:9096 &

```
