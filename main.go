package main

import (
	"flag"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var config = struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	ListenAddress   string
}{}

func main() {
	flag.StringVar(&(config.AccessKeyId), "id", os.Getenv("ACCESS_KEY_ID"), "Access key ID")
	flag.StringVar(&(config.AccessKeySecret), "secret", os.Getenv("ACCESS_KEY_SECRET"), "Access key secret")
	flag.StringVar(&(config.Region), "region", "cn-hangzhou", "The region")
	flag.StringVar(&(config.ListenAddress), "listenaddress", ":8080", "The address it will listen on")
	flag.Parse()

	startAt := time.Now()
	dangledPeriod := func() float64 {
		return 60 + 2 + math.Sin(float64(time.Since(startAt) / (10 * time.Minute)))
	}

	go func() {
		for {
			collectSLBInfo()
			collectRDSInfo()
			time.Sleep(time.Duration(dangledPeriod()) * time.Minute)
		}
	}()

	exporter := NewExporter(newCmsClient())
	prometheus.MustRegister(instance_info, exporter)

	// serve on /metrics
	listenAddress := config.ListenAddress
	log.Println("Running on ", listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
