package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const metricsNamespace = "horaclifix"

var (
	packetsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "packets_total",
			Help:      "Number of received packets",
		},
		[]string{"packets"},
	)
)
