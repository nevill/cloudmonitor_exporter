package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	slbCacheFile = "/tmp/slb.cache"
	rdsCacheFile = "/tmp/rds.cache"
)

var (
	config   = getConfigFromEnv()
	exporter = NewExporter(newCmsClient())
	slbName  = ReadCache(slbCacheFile) // Read SLB Instance Cache
	rdsName  = ReadCache(rdsCacheFile) // Read RDS Instance Cache
)

func start() {
	listenAddress := config.ListenAddress
	if len(listenAddress) == 0 {
		listenAddress = ":8080"
	}

	log.Println("Running on ", listenAddress)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func init() {
	// register metrics to Prometheus
	prometheus.MustRegister(exporter)

	// Write and read cache RDS name and ID to local file
	timedTask(newSLBClient(), newRDSClient())
}

func main() {
	start()
}
