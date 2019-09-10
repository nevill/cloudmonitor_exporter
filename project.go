package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// datapoint represents the member of Datapoints field from QueryMetricLastResponse
type datapoint struct {
	Average    float64 `json:"Average"`
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Value      float64 `json:"Value"`
	InstanceId string  `json:"instanceId"`
	Timestamp  int64   `json:"timestamp"`
	UserId     string  `json:"userId"`
	Port       string  `json:port`
	Vip        string  `json:vip`
}

// GetResponseFunc returns a function to retrieve queryMetricLast
type GetResponseFunc func(client *cms.Client, request *cms.DescribeMetricLastRequest) (string, error)

// Project represents the dashborad from which metrics collected
type Project struct {
	client      *cms.Client
	getResponse GetResponseFunc
	Namespace   string
}

func defaultGetResponseFunc(client *cms.Client, request *cms.DescribeMetricLastRequest) (string, error) {
	response, err := client.DescribeMetricLast(request)
	if err != nil {
		return nil, err
	}
	else {
		return response.Datapoints, error
	}
}

func retrieve(metric string, p Project) []datapoint {
	request := cms.CreateDescribeMetricLastRequest()
	request.Namespace = p.Namespace
	request.MetricName = metric

	requestsStats.Inc()

	datapoints := make([]datapoint, 0)

	getResponseFunc := p.getResponse
	if getResponseFunc == nil {
		getResponseFunc = defaultGetResponseFunc
	}

	source, err := getResponseFunc(p.client, request)

	if (err != nil) {
		responseError.Inc()
		log.Println("Encounter response error from Aliyun:", err)
	} else if err := json.Unmarshal([]byte(source), &datapoints); err != nil {
		responseFormatError.Inc()
		log.Println("Cannot decode json reponse:", err)
	}

	return datapoints
}
