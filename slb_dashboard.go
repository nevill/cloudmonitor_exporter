package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// SLBDashboard represents the dashboard of SLB
type SLBDashboard struct {
	project Project
}

// NewSLBDashboard returns a project respresents acs_slb_dashboard
func NewSLBDashboard(c *cms.Client) *SLBDashboard {
	return &SLBDashboard{
		project: Project{
			client: c,
			Name:   "acs_slb_dashboard",
		},
	}
}

// below are Layer-4 metrics

func (db *SLBDashboard) retrieveActiveConnection() []datapoint {
	return retrieve("ActiveConnection", db.project)
}

func (db *SLBDashboard) retrievePacketTX() []datapoint {
	return retrieve("PacketTX", db.project)
}

func (db *SLBDashboard) retrievePacketRX() []datapoint {
	return retrieve("PacketRX", db.project)
}

func (db *SLBDashboard) retrieveTrafficRX() []datapoint {
	return retrieve("TrafficRXNew", db.project)
}

func (db *SLBDashboard) retrieveTrafficTX() []datapoint {
	return retrieve("TrafficTXNew", db.project)
}

func (db *SLBDashboard) retrieveNewConnection() []datapoint {
	return retrieve("NewConnection", db.project)
}

func (db *SLBDashboard) retrieveMaxConnection() []datapoint {
	return retrieve("MaxConnection", db.project)
}

func (db *SLBDashboard) retrieveDropConnection() []datapoint {
	return retrieve("DropConnection", db.project)
}

func (db *SLBDashboard) retrieveDropPacketRX() []datapoint {
	return retrieve("DropPacketRX", db.project)
}

func (db *SLBDashboard) retrieveDropPacketTX() []datapoint {
	return retrieve("DropPacketTX", db.project)
}

func (db *SLBDashboard) retrieveDropTrafficRX() []datapoint {
	return retrieve("DropTrafficRX", db.project)
}

func (db *SLBDashboard) retrieveDropTrafficTX() []datapoint {
	return retrieve("DropTrafficTX", db.project)
}

// below are Layer-7 metrics

func (db *SLBDashboard) retrieveQps() []datapoint {
	return retrieve("Qps", db.project)
}

func (db *SLBDashboard) retrieveRt() []datapoint {
	return retrieve("Rt", db.project)
}

func (db *SLBDashboard) retrieveStatusCode5xx() []datapoint {
	return retrieve("StatusCode5xx", db.project)
}

func (db *SLBDashboard) retrieveUpstreamCode4xx() []datapoint {
	return retrieve("UpstreamCode4xx", db.project)
}

func (db *SLBDashboard) retrieveUpstreamCode5xx() []datapoint {
	return retrieve("UpstreamCode5xx", db.project)
}

func (db *SLBDashboard) retrieveUpstreamRt() []datapoint {
	return retrieve("UpstreamRt", db.project)
}
