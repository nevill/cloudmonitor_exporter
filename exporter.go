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
	client *cms.Client

	netTxRate        *prometheus.Desc
	netTxRatePercent *prometheus.Desc
	snatConnections  *prometheus.Desc
}

// NewExporter instantiate an CloudmonitorExport
func NewExporter(c *cms.Client) *CloudmonitorExporter {
	return &CloudmonitorExporter{
		client: c,
		netTxRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "net_tx_rate", "bytes"),
			"Outbound bandwith of gateway in bits/s",
			[]string{
				"id", // instance id
			},
			nil,
		),
		netTxRatePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "net_tx_rate", "percent"),
			"Outbound bandwith of gateway used in percentage",
			[]string{
				"id", // instance id
			},
			nil,
		),
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
	ch <- e.netTxRate
	ch <- e.netTxRatePercent
	ch <- e.snatConnections
}

// Collect fetches the metrics from Aliyun cms
// It implements prometheus.Collector.
func (e *CloudmonitorExporter) Collect(ch chan<- prometheus.Metric) {
	natGateway := NewNatGateway(e.client)

	for _, point := range natGateway.retrieveNetTxRate() {
		ch <- prometheus.MustNewConstMetric(
			e.netTxRate,
			prometheus.GaugeValue,
			float64(point.Value),
			point.InstanceId,
		)
	}

	for _, point := range natGateway.retrieveNetTxRatePercent() {
		ch <- prometheus.MustNewConstMetric(
			e.netTxRatePercent,
			prometheus.GaugeValue,
			float64(point.Value),
			point.InstanceId,
		)
	}

	for _, point := range natGateway.retrieveSnatConn() {
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

// GetResponseFunc returns a function to retrieve queryMetricLast
type GetResponseFunc func(client *cms.Client, request *cms.QueryMetricLastRequest) string

// Project represents the dashborad from which metrics collected
type Project struct {
	client      *cms.Client
	getResponse GetResponseFunc
	Name        string
}

func defaultGetResponseFunc(client *cms.Client, request *cms.QueryMetricLastRequest) string {
	response, err := client.QueryMetricLast(request)
	if err != nil {
		panic(err)
	}
	return response.Datapoints
}

func (p *Project) retrieve(metric string) []datapoint {
	request := cms.CreateQueryMetricLastRequest()
	request.Project = p.Name
	request.Metric = metric

	requestsStats.Inc()
	source := p.getResponse(p.client, request)

	datapoints := make([]datapoint, 10)
	if err := json.Unmarshal([]byte(source), &datapoints); err != nil {
		panic(err)
	}

	return datapoints
}
