package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "cloudmonitor"
)

// CloudmonitorExporter collects metrics from Aliyun via cms API
type CloudmonitorExporter struct {
	snatConnections *prometheus.Desc
}

// NewExporter instantiate an CloudmonitorExport
func NewExporter() *CloudmonitorExporter {
	return &CloudmonitorExporter{
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
	client := newCmsClient()

	for _, point := range retrieveSnatConn(client) {
		ch <- prometheus.MustNewConstMetric(
			e.snatConnections,
			prometheus.GaugeValue,
			float64(point.Maximum),
			point.InstanceId,
		)
	}

}
