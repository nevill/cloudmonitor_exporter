package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

var cacheName = map[string]map[string]string{
	"rds": map[string]string{},
	"slb": map[string]string{},
}

// ResultResponse Ali cloud interface response field
type ResultResponse struct {
	TotalRecordCount int                                 `json:"TotalRecordCount"`
	LoadBalancers    map[string][]map[string]interface{} `json:"LoadBalancers"`
	Items            map[string][]map[string]interface{} `json:"Items"`
}

// CacheDescriptionSLB Call Ali interface cache instance SLB information
func CacheDescriptionSLB() {
	var (
		result ResultResponse
		client = newSLBClient()
	)
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		log.Println("Cache SLB response error from Aliyun: ", err)
	}
	contentString := response.GetHttpContentString()
	if err := json.Unmarshal([]byte(contentString), &result); err == nil {
		Balancer := result.LoadBalancers["LoadBalancer"]
		for _, v := range Balancer {
			LoadBalancerIDStr := fmt.Sprintf("%v", v["LoadBalancerId"])
			LoadBalancerNameStr := fmt.Sprintf("%v", v["LoadBalancerName"])
			cacheName["slb"][LoadBalancerIDStr] = LoadBalancerNameStr
		}
	}
}

// CacheDescriptionRDS Call Ali interface cache instance RDS information
func CacheDescriptionRDS() {
	var (
		result ResultResponse
		num    = 1
		size   = 100
		client = newRDSClient()
	)
	for PageTurning := true; PageTurning != false; num++ {
		request := rds.CreateDescribeDBInstancesRequest()
		request.Scheme = "https"
		request.PageSize = requests.NewInteger(100)
		request.PageNumber = requests.NewInteger(num)

		response, err := client.DescribeDBInstances(request)
		if err != nil {
			log.Println("Cache RDS response error from Aliyun: ", err)
		}
		contentString := response.GetHttpContentString()

		if err := json.Unmarshal([]byte(contentString), &result); err == nil {
			totalCount := result.TotalRecordCount
			DBInstances := result.Items["DBInstance"]
			for _, v := range DBInstances {
				DBInstanceIDStr := fmt.Sprintf("%v", v["DBInstanceId"])
				DBInstanceDescriptionStr := fmt.Sprintf("%v", v["DBInstanceDescription"])
				cacheName["rds"][DBInstanceIDStr] = DBInstanceDescriptionStr
			}
			if totalCount > size {
				size = size + 100
			} else {
				PageTurning = false
			}
		}
	}
}
