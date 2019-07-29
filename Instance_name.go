package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

// ResultResponse Ali cloud interface response field
type ResultResponse struct {
	TotalRecordCount int                                 `json:"TotalRecordCount"`
	LoadBalancers    map[string][]map[string]interface{} `json:"LoadBalancers"`
	Items            map[string][]map[string]interface{} `json:"Items"`
}

// WriteCache Write Cache instance information
func WriteCache(outputFile string, strContent []byte) {
	fd, _ := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	fd.Write(strContent)
	fd.Close()
}

// ReadCache Read Cache instance information
func ReadCache(inputFile string) map[string]string {
	buf, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		// panic(err.Error())
	}
	var mapResult map[string]string
	if err := json.Unmarshal([]byte(buf), &mapResult); err != nil {
		fmt.Fprintf(os.Stderr, "File Format Error: %s\n", err)
		// panic(err)
	}

	return mapResult
}

// CacheDescriptionSLB Call Ali interface cache instance SLB information
func CacheDescriptionSLB(client *slb.Client) {
	request := slb.CreateDescribeLoadBalancersRequest()
	request.Scheme = "https"
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cache SLB response error from Aliyun: %s\n", err.Error())
	}
	contentString := response.GetHttpContentString()
	data := make(map[string]interface{})
	var result ResultResponse
	if err := json.Unmarshal([]byte(contentString), &result); err == nil {
		Balancer := result.LoadBalancers["LoadBalancer"]
		for _, v := range Balancer {
			LoadBalancerIDStr := fmt.Sprintf("%v", v["LoadBalancerId"])
			data[LoadBalancerIDStr] = v["LoadBalancerName"]
		}
	}
	pureRes, err := json.Marshal(data)
	if err == nil {
		WriteCache(slbCacheFile, pureRes)
	}
}

// CacheDescriptionRDS Call Ali interface cache instance RDS information
func CacheDescriptionRDS(client *rds.Client) {
	var (
		num  = 0
		size = 100
		data = make(map[string]interface{})
	)
PageTurning:
	num++
	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	request.PageSize = requests.NewInteger(100)
	request.PageNumber = requests.NewInteger(num)

	response, err := client.DescribeDBInstances(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cache RDS response error from Aliyun: %s\n", err.Error())
	}
	contentString := response.GetHttpContentString()

	var result ResultResponse
	if err := json.Unmarshal([]byte(contentString), &result); err == nil {
		totalCount := result.TotalRecordCount
		DBInstances := result.Items["DBInstance"]
		for _, v := range DBInstances {
			DBInstanceID := fmt.Sprintf("%v", v["DBInstanceId"])
			data[DBInstanceID] = v["DBInstanceDescription"]
		}
		if totalCount > size {
			size = size + 100
			goto PageTurning
		}
	}
	pureRes, err := json.Marshal(data)
	if err == nil {
		WriteCache(rdsCacheFile, pureRes)
	}
}

// timedTask 循环定时任务
func timedTask(slb *slb.Client, rds *rds.Client) {
	go func() {
		for {
			CacheDescriptionSLB(slb)
			CacheDescriptionRDS(rds)
			slbName = ReadCache(slbCacheFile)
			rdsName = ReadCache(rdsCacheFile)
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 10, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}
