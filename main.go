package main

import (
	"log"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	snatMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "snat",
			Subsystem: "gateway",
			Name:      "max_connections",
			Help:      "Max number of snat connections per minute",
		},
		[]string{
			"id", // instance id
		},
	)
)

func newCmsClient() *cms.Client {
	config := getConfigFromEnv()
	cmsClient, err := cms.NewClientWithAccessKey(
		config.Region,
		config.AccessKeyId,
		config.AccessKeySecret,
	)

	if err != nil {
		panic(err)
	}

	return cmsClient
}

func start() {
	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	// register metrics to Prometheus
	prometheus.MustRegister(snatMetrics)
}

func main() {
	client := newCmsClient()

	for _, point := range retrieveSnatConn(client) {
		snatMetrics.WithLabelValues(point.InstanceId).Set(float64(point.Maximum))
	}

	start()
}
