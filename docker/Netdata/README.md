# Quickstart for Netdata and statsd with docker-compose 

## Pre-requisites
* Latest docker version and [docker-compose](https://docs.docker.com/compose/install/)

## Installation & Configuration
Clone the project locally to your Docker host and run from the horaclifix/docker/Netdata directory:

	$ ./replaceIP.sh 1.1.1.1 YOUR_OWN_IP
    $ docker-compose up -d


Since we don't want to listen on 0.0.0.0 use the replaceIP script which will replace the preconfigured IP 1.1.1.1 with your own IP.

## Metrics
You can define your own metrics in sbc_qos.conf.