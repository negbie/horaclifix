package main

func (conn *Connections) SendSyslog(i *IPFIX, s string) {
	payload := i.mapQOS()
	_, err := conn.Syslog.Write(payload)
	checkErr(err)
}
