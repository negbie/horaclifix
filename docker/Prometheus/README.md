[![Build Status](https://travis-ci.org/vegasbrianc/prometheus.svg?branch=version-2)](https://travis-ci.org/vegasbrianc/prometheus)

# A Prometheus & Grafana docker-compose stack

Here's a quick start to stand-up a [Prometheus](http://prometheus.io/) stack containing Prometheus, statsd & node exporter, Grafana.

## Pre-requisites
* Latest docker version and [docker-compose](https://docs.docker.com/compose/install/)

## Installation & Configuration
Clone the project locally to your Docker host and run from the horaclifix/docker/Prometheus directory:

	$ ./replaceIP.sh 1.1.1.1 YOUR_OWN_IP
    $ docker-compose up -d


Since we don't want to listen on 0.0.0.0 use the replaceIP script which will replace the preconfigured IP 1.1.1.1 with your own IP.

The Grafana Dashboard is now accessible via: `http://<YOUR_OWN_IP>:3000`

* username:admin
* password:foobar (Password is stored in the `config.monitoring` env file)

## Post Configurations
Now we need to create the Prometheus Datasource in order to connect Grafana to Prometheus 
* Click the `Grafana` Menu at the top left corner
* Click `Data Sources`
* Click the green button `Add Data Source`.

<img src="https://github.com/vegasbrianc/prometheus/blob/version-2/images/Add_Data_Source.png" width="400" heighth="400">

## Alerting
Alerting has been added to the stack with Slack integration. 2 Alerts have been added and are managed 

Alerts              - `prometheus/alert.rules`
Slack configuration - `alertmanager/config.yml`

The Slack configuration requires to build a custom integration.
* Open your slack team in your browser `https://<your-slack-team>.slack.com/apps`
* Click build in the upper right corner
* Make a Custom integration
* Choose Incoming Web Hooks
* Select which channel
* Click on Add Incoming WebHooks integration
* Copy the Webhook URL into the `alertmanager/config.yml` URL section
* Fill in Slack username and channel

View Prometheus alerts `http://<Host IP Address>:9090/alerts`
View Alert Manager `http://<Host IP Address>:9093`

### Test Alerts
A quick test for your alerts is to stop a service. Stop the node_exporter container and you should notice shortly the alert arrive in Slack. Also check the alerts in both the Alert Manager and Prometheus Alerts just to understand how they flow through the system.

High load test alert - `docker run --rm -it busybox sh -c "while true; do :; done"`

Let this run for a few minutes and you will notice the load alert appear.

## Install Dashboard
I created a Dashboard template which is available on [Grafana Docker Dashboard](https://grafana.net/dashboards/179). Simply download the dashboard and select from the Grafana menu -> Dashboards -> Import

This dashboard is intended to help you get started with monitoring. If you have any changes you would like to see in the Dashboard let me know so I can update Grafana site as well.

Here's the Dashboard Template

![Grafana Dashboard](https://github.com/vegasbrianc/prometheus/blob/version-2/images/Dashboard.png)

Grafana Dashboard - `dashboards/Grana_Dashboad.json`
Alerting Dashboard - `dashboards/System_Monitoring.json`

## Troubleshooting
It appears some people have reported no data appearing in Grafana. If this is happening to you be sure to check the time range being queried within Grafana to ensure it is using Today's date with current time.
