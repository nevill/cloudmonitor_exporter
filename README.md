# Cloudmonitor Exporter

通过阿里云提供的 API 来搜集云上的监控数据

同其它 exporter 一样，需要配合 prometheus 来使用

# Build
```
go get -d
go build
```

# Test
```
ACCESS_KEY_ID=fakeid go test
```

# API Reference
1. https://help.aliyun.com/document_detail/28617.html
2. https://help.aliyun.com/document_detail/51939.html

# License

Apache License 2.0

Same to [alibaba-cloud-sdk-go](https://github.com/aliyun/alibaba-cloud-sdk-go)
