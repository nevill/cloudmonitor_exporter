package main

import (
	"flag"
	"log"
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

func start() {
	exporter := NewExporter(newCmsClient())
	prometheus.MustRegister(exporter)

	listenAddress := config.ListenAddress

	log.Println("Running on ", listenAddress)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func main() {
	flag.StringVar(&(config.AccessKeyId), "id", os.Getenv("ACCESS_KEY_ID"), "Access key ID")
	flag.StringVar(&(config.AccessKeySecret), "secret", os.Getenv("ACCESS_KEY_SECRET"), "Access key secret")
	flag.StringVar(&(config.Region), "region", "cn-hangzhou", "The region")
	flag.StringVar(&(config.ListenAddress), "listenaddress", ":8080", "The address it will listen on")
	flag.Parse()

	go timedTask()
	start()
}

// timedTask 循环定时任务
func timedTask() {
	for {
		CacheDescriptionSLB()
		CacheDescriptionRDS()
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 10, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}
}
