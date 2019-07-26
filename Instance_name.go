package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// ResultResponse Ali cloud interface response field
type ResultResponse struct {
	PageNumber    int                                 `json:"PageNumber"`
	TotalCount    int                                 `json:"TotalCount"`
	PageSize      int                                 `json:"PageSize"`
	LoadBalancers map[string][]map[string]interface{} `json:"LoadBalancers"`
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
func CacheDescriptionSLB(client *sdk.Client) {
	var result ResultResponse
	request := requests.NewCommonRequest()
	request.Domain = "slb.aliyuncs.com"
	request.Version = "2014-05-15"
	request.ApiName = "DescribeLoadBalancers"
	// request.QueryParams["ServerIntranetAddress"] = "10.32.12.19"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cache response error from Aliyun: %s\n", err)
		// panic(err)
	}
	res := response.GetHttpContentString()
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(res), &result); err == nil {
		Balancer := result.LoadBalancers["LoadBalancer"]
		for _, v := range Balancer {
			LoadBalancerIDStr := fmt.Sprintf("%v", v["LoadBalancerId"])
			data[LoadBalancerIDStr] = v["LoadBalancerName"]
		}
	}
	pureRes, err := json.Marshal(data)
	if err == nil {
		WriteCache("/tmp/slb.cache", pureRes)
	}
}

// timedTask 循环定时任务
func timedTask(client *sdk.Client) {
	go func() {
		for {
			CacheDescriptionSLB(client)
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 10, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}
