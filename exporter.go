package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "cloudmonitor"
)

// CloudmonitorExporter collects metrics from Aliyun via cms API
type CloudmonitorExporter struct {
	client *cms.Client

	// nat gateway
	netTxRate        *prometheus.Desc
	netTxRatePercent *prometheus.Desc
	snatConnections  *prometheus.Desc

	// slb dashbaord
	activeConnection *prometheus.Desc
	packetRX         *prometheus.Desc
	packetTX         *prometheus.Desc
	trafficRX        *prometheus.Desc
	trafficTX        *prometheus.Desc
	newConnection    *prometheus.Desc
	maxConnection    *prometheus.Desc
	dropConnection   *prometheus.Desc
	dropPacketRX     *prometheus.Desc
	dropPacketTX     *prometheus.Desc
	dropTrafficRX    *prometheus.Desc
	dropTrafficTX    *prometheus.Desc
	qps              *prometheus.Desc
	rt               *prometheus.Desc
	statusCode5xx    *prometheus.Desc
	upstreamCode4xx  *prometheus.Desc
	upstreamCode5xx  *prometheus.Desc
	upstreamRt       *prometheus.Desc

	// rds dashbaord
	cpuUsage        *prometheus.Desc
	connectionUsage *prometheus.Desc
	activeSessions  *prometheus.Desc
}

// NewExporter instantiate an CloudmonitorExport
func NewExporter(c *cms.Client) *CloudmonitorExporter {
	return &CloudmonitorExporter{
		client: c,

		// nat gateway
		netTxRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "net_tx_rate", "bytes"),
			"Outbound bandwith of gateway in bits/s",
			[]string{
				"id",
			},
			nil,
		),
		netTxRatePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "net_tx_rate", "percent"),
			"Outbound bandwith of gateway used in percentage",
			[]string{
				"id",
			},
			nil,
		),
		snatConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "snat", "connections"),
			"Max number of snat connections per minute",
			[]string{
				"id",
			},
			nil,
		),

		// slb dashboard
		activeConnection: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slb", "active_connection"),
			"Number of active connections per minute",
			[]string{
				"id",
				"port",
				"vip",
			},
			nil,
		),

		packetRX: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slb", "packet_rx_average"),
			"Average packets received per second",
			[]string{
				"id",
				"port",
				"vip",
			},
			nil,
		),

		packetTX: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slb", "packet_tx_average"),
			"Average packets sent per second",
			[]string{
				"id",
				"port",
				"vip",
			},
			nil,
		),

		trafficRX: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slb", "traffic_rx_average"),
			"Average traffic received per second",
			[]string{
				"id",
				"port",
				"vip",
			},
			nil,
		),

		trafficTX: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slb", "traffic_tx_average"),
			"Average traffic sent per second",
			[]string{
				"id",
				"port",
				"vip",
			},
			nil,
		),

		newConnection: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slb", "new_connection_average"),
			"Average number of new connections created per second",
			[]string{
				"id",
				"port",
				"vip",
			},
			nil,
		),

		// rds dashbaord
		cpuUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "rds", "cpu_usage_average"),
			"CPU usage per minute",
			[]string{
				"id",
			},
			nil,
		),

		connectionUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "rds", "connection_usage"),
			"Connection usage per minute",
			[]string{
				"id",
			},
			nil,
		),

		activeSessions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "rds", "active_sessions"),
			"Active Sessions per minute",
			[]string{
				"id",
			},
			nil,
		),
	}
}

// Describe describes all the metrics exported by the cloudmonitor exporter.
// It implements prometheus.Collector.
func (e *CloudmonitorExporter) Describe(ch chan<- *prometheus.Desc) {
	// nat gateway
	ch <- e.netTxRate
	ch <- e.netTxRatePercent
	ch <- e.snatConnections

	// slb dashboard
	ch <- e.activeConnection
	ch <- e.packetRX
	ch <- e.packetTX
	ch <- e.trafficRX
	ch <- e.trafficTX
	ch <- e.newConnection

	// rds dashbaord
	ch <- e.cpuUsage
	ch <- e.connectionUsage
	ch <- e.activeSessions
}

// Collect fetches the metrics from Aliyun cms
// It implements prometheus.Collector.
func (e *CloudmonitorExporter) Collect(ch chan<- prometheus.Metric) {
	natGateway := NewNatGateway(e.client)
	slbDashboard := NewSLBDashboard(e.client)
	rdsDashboard := NewRDSDashboard(e.client)

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

	for _, point := range slbDashboard.retrieveActiveConnection() {
		ch <- prometheus.MustNewConstMetric(
			e.activeConnection,
			prometheus.GaugeValue,
			float64(point.Maximum),
			point.InstanceId+"("+cacheName["slb"][point.InstanceId]+")",
			point.Port,
			point.Vip,
		)
	}

	for _, point := range slbDashboard.retrievePacketRX() {
		ch <- prometheus.MustNewConstMetric(
			e.packetRX,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["slb"][point.InstanceId]+")",
			point.Port,
			point.Vip,
		)
	}

	for _, point := range slbDashboard.retrievePacketTX() {
		ch <- prometheus.MustNewConstMetric(
			e.packetTX,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["slb"][point.InstanceId]+")",
			point.Port,
			point.Vip,
		)
	}

	for _, point := range slbDashboard.retrieveTrafficRX() {
		ch <- prometheus.MustNewConstMetric(
			e.trafficRX,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["slb"][point.InstanceId]+")",
			point.Port,
			point.Vip,
		)
	}

	for _, point := range slbDashboard.retrieveTrafficTX() {
		ch <- prometheus.MustNewConstMetric(
			e.trafficTX,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["slb"][point.InstanceId]+")",
			point.Port,
			point.Vip,
		)
	}

	for _, point := range slbDashboard.retrieveNewConnection() {
		ch <- prometheus.MustNewConstMetric(
			e.newConnection,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["slb"][point.InstanceId]+")",
			point.Port,
			point.Vip,
		)
	}

	for _, point := range rdsDashboard.retrieveCPUUsage() {
		ch <- prometheus.MustNewConstMetric(
			e.cpuUsage,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["rds"][point.InstanceId]+")",
		)
	}

	for _, point := range rdsDashboard.retrieveConnectionUsage() {
		ch <- prometheus.MustNewConstMetric(
			e.connectionUsage,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["rds"][point.InstanceId]+")",
		)
	}

	for _, point := range rdsDashboard.retrieveActiveSessions() {
		ch <- prometheus.MustNewConstMetric(
			e.activeSessions,
			prometheus.GaugeValue,
			float64(point.Average),
			point.InstanceId+"("+cacheName["rds"][point.InstanceId]+")",
		)
	}
}
