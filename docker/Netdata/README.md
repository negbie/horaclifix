# Quickstart for Netdata and statsd with docker-compose 

## Pre-requisites
* Latest docker version and [docker-compose](https://docs.docker.com/compose/install/)

## Installation & Configuration
Clone the project locally to your Docker host and run from the horaclifix/docker/Netdata directory:

	$ ./replaceIP.sh 1.1.1.1 YOUR_OWN_IP
    $ docker-compose up -d


## Metrics
You can define your own metrics in sbc_qos.conf.
