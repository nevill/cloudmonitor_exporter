package main

import (
	"log"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	config   = getConfigFromEnv()
	exporter = NewExporter()
)

func newCmsClient() *cms.Client {
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
	listenAddress := config.ListenAddress
	if len(listenAddress) == 0 {
		listenAddress = ":8080"
	}

	log.Println("Running on ", listenAddress)

	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func init() {
	// register metrics to Prometheus
	prometheus.MustRegister(exporter)
}

func main() {
	start()
}
