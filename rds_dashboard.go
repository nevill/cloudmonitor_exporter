package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// RDSDashboard represents the dashboard of RDS
type RDSDashboard struct {
	project Project
}

// NewRDSDashboard returns a project respresents acs_rds_dashboard
func NewRDSDashboard(c *cms.Client) *RDSDashboard {
	return &RDSDashboard{
		project: Project{
			client:    c,
			Namespace: "acs_rds_dashboard",
		},
	}
}

func (db *RDSDashboard) retrieveCPUUsage() []datapoint {
	return retrieve("CpuUsage", db.project)
}

func (db *RDSDashboard) retrieveConnectionUsage() []datapoint {
	return retrieve("ConnectionUsage", db.project)
}

func (db *RDSDashboard) retrieveActiveSessions() []datapoint {
	return retrieve("MySQL_ActiveSessions", db.project)
}
