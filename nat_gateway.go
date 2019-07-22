package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// NatGateway represents nat gateway dashboard
type NatGateway struct {
	project Project
}

// NewNatGateway returns a project respresents acs_nat_gateway
func NewNatGateway(c *cms.Client) *NatGateway {
	return &NatGateway{
		project: Project{
			client:    c,
			Namespace: "acs_nat_gateway",
		},
	}
}

func (db *NatGateway) retrieveNetTxRate() []datapoint {
	return retrieve("net_tx.rate", db.project)
}

func (db *NatGateway) retrieveNetTxRatePercent() []datapoint {
	return retrieve("net_tx.ratePercent", db.project)
}

func (db *NatGateway) retrieveSnatConn() []datapoint {
	return retrieve("SnatConnection", db.project)
}
