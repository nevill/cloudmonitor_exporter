package main

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "cloudmonitor"
)

// CloudmonitorExporter collects metrics from Aliyun via cms API
type CloudmonitorExporter struct {
	client          *cms.Client
	snatConnections *prometheus.Desc
}

// NewExporter instantiate an CloudmonitorExport
func NewExporter(c *cms.Client) *CloudmonitorExporter {
	return &CloudmonitorExporter{
		client: c,
		snatConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "snat", "connections"),
			"Max number of snat connections per minute",
			[]string{
				"id", // instance id
			},
			nil,
		),
	}
}

// Describe describes all the metrics exported by the cloudmonitor exporter.
// It implements prometheus.Collector.
func (e *CloudmonitorExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.snatConnections
}

// Collect fetches the metrics from Aliyun cms
// It implements prometheus.Collector.
func (e *CloudmonitorExporter) Collect(ch chan<- prometheus.Metric) {
	project := NewNatGateway(e.client)

	for _, point := range project.retrieveSnatConn() {
		ch <- prometheus.MustNewConstMetric(
			e.snatConnections,
			prometheus.GaugeValue,
			float64(point.Maximum),
			point.InstanceId,
		)
	}

}

// datapoint represents the member of Datapoints field from QueryMetricLastResponse
type datapoint struct {
	Average    float64 `json:"Average"`
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Value      float64 `json:"Value"`
	InstanceId string  `json:"instanceId"`
	Timestamp  int64   `json:"timestamp"`
	UserId     string  `json:"userId"`
}

// Project represents the dashborad from which metrics collected
type Project struct {
	client *cms.Client
	Name   string
}

func (p *Project) retrieve(name string) []datapoint {
	request := cms.CreateQueryMetricLastRequest()
	request.Project = p.Name
	request.Metric = name
	response, err := p.client.QueryMetricLast(request)

	if err != nil {
		panic(err)
	}

	source := response.Datapoints
	datapoints := make([]datapoint, 10)
	if err := json.Unmarshal([]byte(source), &datapoints); err != nil {
		panic(err)
	}

	return datapoints
}
