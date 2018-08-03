package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// NewNatGateway returns a project respresents acs_nat_gateway
func NewNatGateway(c *cms.Client) *Project {
	return &Project{
		client: c,
		Name:   "acs_nat_gateway",
	}
}

func (p *Project) retrieveSnatConn() []datapoint {
	return p.retrieve("SnatConnection")
}

func (p *Project) retrieveNetRx() []datapoint {
	return p.retrieve("net_tx.rate")
}
