package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// NewNatGateway returns a project respresents acs_nat_gateway
func NewNatGateway(c *cms.Client) *Project {
	return &Project{
		client:      c,
		getResponse: defaultGetResponseFunc,
		Name:        "acs_nat_gateway",
	}
}

func (p *Project) retrieveNetTxRate() []datapoint {
	return p.retrieve("net_tx.rate")
}

func (p *Project) retrieveNetTxRatePercent() []datapoint {
	return p.retrieve("net_tx.ratePercent")
}

func (p *Project) retrieveSnatConn() []datapoint {
	return p.retrieve("SnatConnection")
}
