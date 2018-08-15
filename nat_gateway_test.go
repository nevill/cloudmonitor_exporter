package main

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

func TestRetrieveSnatConn(t *testing.T) {
	t.Log("Should retrieve maximum number of snat connections \n")
	client := newCmsClient()
	gw := NewNatGateway(client)
	gw.getResponse = func(client *cms.Client, request *cms.QueryMetricLastRequest) string {
		return `[
			{
				"timestamp":25336,
				"userId":"1234567",
				"instanceId":"ngw-instance1",
				"Maximum":7157
			},
			{
				"timestamp":25336,
				"userId":"1234567",
				"instanceId":"ngw-instance2",
				"Maximum":5
			}
		]`
	}

	datapoints := gw.retrieveSnatConn()
	expected := 2
	if len(datapoints) != expected {
		t.Errorf("Want to return %d data points, but get %d\n", expected, len(datapoints))
	}

	dp1 := datapoints[0]
	dp2 := datapoints[1]

	tests := []struct {
		input float64
		want  float64
		field string
	}{
		{
			input: dp1.Maximum,
			want:  7157,
			field: "Maximum",
		},
		{
			input: float64(dp1.Timestamp),
			want:  25336,
			field: "timestamp",
		},
		{
			input: dp2.Maximum,
			want:  5,
			field: "Maximum",
		},
		{
			input: float64(dp2.Timestamp),
			want:  25336,
			field: "timestamp",
		},
	}

	for _, test := range tests {
		value := test.input
		if value != test.want {
			t.Errorf("Field %s hould return %f, but get %f\n", test.field, test.want, value)
		}
	}
}

func TestRetrieveNetTxRatePercent(t *testing.T) {
	t.Log("Should retrieve outbound of gateway used in percent\n")
	client := newCmsClient()
	gw := NewNatGateway(client)
	gw.getResponse = func(client *cms.Client, request *cms.QueryMetricLastRequest) string {
		return `[
			{
				"timestamp":1820000,
				"userId":"78901234",
				"instanceId":"bwp-instance1",
				"Value":62.11
			},
			{
				"timestamp":1820000,
				"userId":"78901234",
				"instanceId":"bwp-instance2",
				"Value":0.015
			}
		]`
	}

	datapoints := gw.retrieveNetTxRatePercent()
	expected := 2
	if len(datapoints) != expected {
		t.Errorf("Want to return %d data points, but get %d\n", expected, len(datapoints))
	}

	dp1 := datapoints[0]
	dp2 := datapoints[1]

	tests := []struct {
		input float64
		want  float64
		field string
	}{
		{
			input: dp1.Value,
			want:  62.11,
			field: "Value",
		},
		{
			input: float64(dp1.Timestamp),
			want:  1820000,
			field: "timestamp",
		},
		{
			input: dp2.Value,
			want:  0.015,
			field: "Value",
		},
		{
			input: float64(dp2.Timestamp),
			want:  1820000,
			field: "timestamp",
		},
	}

	for _, test := range tests {
		value := test.input
		if value != test.want {
			t.Errorf("Field %s hould return %f, but get %f\n", test.field, test.want, value)
		}
	}
}

func TestRetrieveNetTxRate(t *testing.T) {
	t.Log("Should retrieve outbound of gateway\n")
	client := newCmsClient()
	gw := NewNatGateway(client)
	gw.getResponse = func(client *cms.Client, request *cms.QueryMetricLastRequest) string {
		return `[
			{
				"timestamp":10001,
				"userId":"2345678",
				"instanceId":"bwp-instance1",
				"Value":3256328
			},
			{
				"timestamp":10001,
				"userId":"2345678",
				"instanceId":"bwp-instance4",
				"Value":2.2853296E7
			}
		]`
	}

	datapoints := gw.retrieveNetTxRate()
	expected := 2
	if len(datapoints) != expected {
		t.Errorf("Want to return %d data points, but get %d\n", expected, len(datapoints))
	}

	dp1 := datapoints[0]
	dp2 := datapoints[1]

	tests := []struct {
		input float64
		want  float64
		field string
	}{
		{
			input: dp1.Value,
			want:  3256328,
			field: "Value",
		},
		{
			input: float64(dp1.Timestamp),
			want:  10001,
			field: "timestamp",
		},
		{
			input: dp2.Value,
			want:  22853296,
			field: "Value",
		},
		{
			input: float64(dp2.Timestamp),
			want:  10001,
			field: "timestamp",
		},
	}

	for _, test := range tests {
		value := test.input
		if value != test.want {
			t.Errorf("Field %s hould return %f, but get %f\n", test.field, test.want, value)
		}
	}
}
