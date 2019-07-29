package main

import (
	"fmt"
	"os"

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
		fmt.Fprintf(os.Stderr, "Create SLB Client error from Aliyun: %s\n", err)
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
		fmt.Fprintf(os.Stderr, "Create RDS Client error from Aliyun: %s\n", err)
	}

	return client
}
