/*


chunk type ID payload type chunk type description
-------------------------------------------------
0x0001 uint8 IP protocol family
0x0002 uint8 IP protocol ID
0x0003 inet4-addr IPv4 source address
0x0004 inet4-addr IPv4 destination address
0x0005 inet6-addr IPv6 source address
0x0006 inet6-addr IPv6 destination address
0x0007 uint16 protocol source port (UDP, TCP, SCTP)
0x0008 uint16 protocol destination port (UDP, TCP, SCTP)
0x0009 uint32 timestamp, seconds since 01/01/1970 (epoch)
0x000a uint32 timestamp microseconds offset (added to timestamp)
0x000b uint8 protocol type (SIP/H323/RTP/MGCP/M2UA)
0x000c uint32 capture agent ID (202, 1201, 2033...)
0x000d uint16 keep alive timer (sec)
0x000e octet-string authenticate key (plain text / TLS connection)
0x000f octet-string captured packet payload
0x0010 octet-string captured compressed payload (gzip/inflate)
0x0011 octet-string Internal correlation id
0x0012 uint8 Vland ID
0x0013 octet-string Group ID



chunk protocol ID assigned vendor
---------------------------------
0x00 reserved
0x01 SIP
0x02 XMPP
0x03 SDP
0x04 RTP
0x05 RTCP
0x06 MGCP
0x07 MEGACO (H.248)
0x08 M2UA (SS7/SIGTRAN)
0x09 M3UA (SS7/SIGTRAN)
0x0a IAX
0x0b H3222
0x0c H321
0x0d M2PA


chunk vendor ID assigned vendor
-------------------------------
0x0000 No specific vendor, generic chunk types
0x0001 FreeSWITCH
0x0002 Kamailio
0x0003 OpenSIPS
0x0004 Asterisk
0x0005 Homer
0x0006 SipXecs


*/