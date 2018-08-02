package main

import "os"

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	ListenAddress   string
}

func getConfigFromEnv() *Config {
	config := &Config{
		AccessKeyId:     os.Getenv("ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ACCESS_KEY_SECRET"),
		Region:          os.Getenv("REGION"),
		ListenAddress:   os.Getenv("LISTEN_ADDRESS"),
	}
	if config.AccessKeyId == "" {
		panic("Cannot get ACCESS_KEY_ID from environment variables")
	}
	return config
}
