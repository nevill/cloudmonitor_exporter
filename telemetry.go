package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsStats = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cloudmonitor_exporter_requests_total",
			Help: "The total number of cloudmonitor requests sent to Aliyun.",
		},
	)
)

func init() {
	prometheus.MustRegister(requestsStats)
}
