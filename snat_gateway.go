package main

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

const (
	dashBoard = "acs_nat_gateway"
)

// MaxinumDatapoint represents the member of Datapoints field from QueryMetricLastResponse
type MaxinumDatapoint struct {
	Maximum    int    `json:"Maximum"`
	InstanceId string `json:"instanceId"`
	Timestamp  int64  `json:"timestamp"`
	UserID     string `json:"userId"`
}

func retrieveSnatConn(client *cms.Client) []MaxinumDatapoint {
	request := cms.CreateQueryMetricLastRequest()
	request.Project = dashBoard
	request.Metric = "SnatConnection"
	response, err := client.QueryMetricLast(request)

	if err != nil {
		panic(err)
	}

	source := response.Datapoints
	datapoints := make([]MaxinumDatapoint, 3)
	if err := json.Unmarshal([]byte(source), &datapoints); err != nil {
		panic(err)
	}

	return datapoints
}
