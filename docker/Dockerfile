FROM golang:1-alpine3.10 AS builder

WORKDIR /cloudmonitor_exporter
COPY . .

ARG GOPROXY

# A way to fetch source without git command
RUN GOPROXY=${GOPROXY} go get -d

RUN GOOS=linux go build

FROM alpine:3.10

LABEL maintainer="Nevill <nevill.dutt@gmail.com>"

COPY --from=builder /cloudmonitor_exporter/cloudmonitor_exporter /bin/

USER nobody
EXPOSE 8080/tcp
CMD ["/bin/cloudmonitor_exporter"]
