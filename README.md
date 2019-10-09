# Cloudmonitor Exporter

This exporter should be used in conjuction with Prometheus. It utilises CMS API to collect metrics from Alibaba cloud.

# Build
```
go get -d
go build
```

# Test
```
go test
```

# Run
```
./cloudmonitor_exporter -id access_id -secret access_secret -region cn-hangzhou
```

or by Docker

```
docker-compose up
```

# API Reference
1. https://help.aliyun.com/document_detail/51939.html
2. https://help.aliyun.com/document_detail/27582.html
3. https://help.aliyun.com/document_detail/26232.html

# License

Apache License 2.0

Same to [alibaba-cloud-sdk-go](https://github.com/aliyun/alibaba-cloud-sdk-go)
