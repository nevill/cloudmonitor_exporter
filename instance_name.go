package main

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/prometheus/client_golang/prometheus"
)

var instance_info = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "cloudmonitor_instance_info",
		Help: "A metric with a constant '1' value labeled by service, id, name.",
	},
	[]string{"service", "id", "name"},
)

func collectSLBInfo() {
	client := newSLBClient()
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		log.Println("Get SLB response error: ", err)
	} else {
		for _, lb := range response.LoadBalancers.LoadBalancer {
			instance_info.WithLabelValues("slb", lb.LoadBalancerId, lb.LoadBalancerName).Set(1)
		}
	}
}

func collectRDSInfo() {
	client := newRDSClient()
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		log.Println("Get RDS response error: ", err)
	} else {
		for _, db := range response.Items.DBInstance {
			instance_info.WithLabelValues("rds", db.DBInstanceId, "").Set(1)
		}
	}
}
