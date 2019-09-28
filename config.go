package main

import (
	"flag"
)

var config = struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	ListenAddress   string
}{}

func init() {
	flag.StringVar(&(config.AccessKeyId), "id", "invalid id", "The access key ID")
	flag.StringVar(&(config.AccessKeySecret), "secret", "empty secret", "The access key secret")
	flag.StringVar(&(config.Region), "region", "cn-hangzhou", "The region")
	flag.StringVar(&(config.ListenAddress), "listenaddress", ":8080", "The address it will listen on")
}
