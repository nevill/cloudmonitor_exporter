package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

func newCmsClient() *cms.Client {
	cmsClient, err := cms.NewClientWithAccessKey(
		config.Region,
		config.AccessKeyId,
		config.AccessKeySecret,
	)

	if err != nil {
		panic(err)
	}

	return cmsClient
}

func newSLBClient() *slb.Client {
	client, err := slb.NewClientWithAccessKey(
		config.Region,
		config.AccessKeyId,
		config.AccessKeySecret,
	)
	if err != nil {
		panic(err)
	}

	return client
}

func newRDSClient() *rds.Client {
	client, err := rds.NewClientWithAccessKey(
		config.Region,
		config.AccessKeyId,
		config.AccessKeySecret,
	)
	if err != nil {
		panic(err)
	}

	return client
}
