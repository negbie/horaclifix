# HORACLIFIX Docker
Docker Container running GO [Horaclifix](https://github.com/negbie/horaclifix)

#### About 
Horaclifix decodes IPFIX messages from Oracle SBC's and converts them into HEP Messages and Stats.
It supports HEP, StatsD, InfluxDB and Graylog output, configured with the command argument
inside the docker-compose.yml

##### Usage
Change inside the docker-compose.yml the command argument to the command line argument's you want for horaclifix.
This example command will send messages to Homer, Graylog and InfluxDB at 192.168.2.22 with the name acme9000.

```command: "./horaclifix -H 192.168.2.22:9060 -g 192.168.2.22:4488 -I 192.168.2.22:8086 -n acme9000 -v"```

Now you can simply run

```docker-compose up -d```