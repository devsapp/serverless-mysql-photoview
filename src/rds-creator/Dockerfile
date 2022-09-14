FROM golang:1.18.1 AS builder


WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

COPY code/ code/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o server ./code/server.go

FROM registry.cn-hangzhou.aliyuncs.com/serverless_devs/terraform-binary:0.0.2 as binary

FROM node:16.15

WORKDIR /

COPY --from=builder /workspace/server .
COPY --from=binary bin/.terraform.d /root/.terraform.d
COPY --from=binary bin/terraform /usr/bin/terraform
COPY --from=binary dev_bin/.terraformrc /root/.terraformrc

COPY data/serverless-rds.tf /data/serverless-rds.tf


ENTRYPOINT ["/server"]
