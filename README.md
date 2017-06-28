![horaclifix](https://user-images.githubusercontent.com/20154956/27519118-9a0720a4-59ed-11e7-95ba-f0e9ce529624.png)

# horaclifix
horaclifix sends IPFIX messages from Oracle SBC's into Homer


### Install:

Get it from the releases:
https://github.com/negbie/horaclifix/releases

Or:
```bash
go install github.com/negbie/horaclifix
```


### Usage of ./horaclifix:

```bash
  -H string
        Homer server address (default "127.0.0.1:9060")
  -d    Debug output to stdout
  -g string
        Graylog server address (default "127.0.0.1:4488")
  -l string
        Host ipfix listen address (default ":4739")

        
################################################################

./horaclifix -H 192.168.2.22:9060
./horaclifix -H 192.168.2.22:9060 -d
```
