package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cacheName = make(map[string]map[string]string) // Cache variable
	config    = getConfigFromEnv()
	exporter  = NewExporter(newCmsClient())
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

	// Refresh the cache once a day
	timedTask()
}

func main() {
	start()
}
