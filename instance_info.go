package main

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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
	req := slb.CreateDescribeLoadBalancersRequest()
	req.Scheme = "https"
	response, err := client.DescribeLoadBalancers(req)
	if err != nil {
		log.Println("Get SLB response error: ", err)
	} else {
		for _, lb := range response.LoadBalancers.LoadBalancer {
			instance_info.WithLabelValues("slb", lb.LoadBalancerId, lb.LoadBalancerName).Set(1)
		}
	}
}

func collectRDSInfo() {
	// 100 is the maximum number of items can be returned in One request
	const pageSize int = 100

	client := newRDSClient()
	req := rds.CreateDescribeDBInstancesRequest()
	req.Scheme = "https"
	req.PageSize = requests.NewInteger(pageSize)

	for hasNextPage, pageNum := true, 1; hasNextPage != false; pageNum++ {
		req.PageNumber = requests.NewInteger(pageNum)
		response, err := client.DescribeDBInstances(req)
		if response.PageRecordCount < pageSize {
			hasNextPage = false
		}

		if err != nil {
			log.Println("Get RDS response error: ", err)
		} else {
			for _, db := range response.Items.DBInstance {
				//TODO can define which field will be set as name in Config file
				instance_info.WithLabelValues("rds", db.DBInstanceId, db.DBInstanceDescription).Set(1)
			}
		}
	}
}
