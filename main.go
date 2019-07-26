package main

import (
	"log"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config   = getConfigFromEnv()
	exporter = NewExporter(newCmsClient())
)

func newClient() *sdk.Client {
	client, err := sdk.NewClientWithAccessKey(
		config.Region,
		config.AccessKeyId,
		config.AccessKeySecret,
	)
	if err != nil {
		panic(err)
	}

	return client
}

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

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func init() {
	// register metrics to Prometheus
	prometheus.MustRegister(exporter)

	// Cache instance SLB name and ID
	timedTask(newClient())
}

func main() {
	start()
}
