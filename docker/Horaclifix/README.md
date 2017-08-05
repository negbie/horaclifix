# HORACLIFIX Docker
Docker Container running GO [Horaclifix](https://github.com/negbie/horaclifix)

<img src="https://user-images.githubusercontent.com/20154956/28133509-9d9870fa-6740-11e7-9616-3fd0e7e1fa9c.png" />

#### About 
Horaclifix decodes IPFIX messages from Oracle SBC's and converts them into HEP Messages and Stats

#### Usage
Horaclifix supports HEP, StatsD and Graylog output, configured by ENV variables

##### Example
```docker run --env HEP_HOST=127.0.0.1 --env HEP_PORT=9060 -p 4739:4739 qxip/horaclifix-go```


##### ENV VARIABLES
###### HEP
```
HEP_HOST 
HEP_PORT 
```
###### GRAYLOG
```
GRAYLOG_URL
```
###### STATSD
```
STATSD_URL 
```



