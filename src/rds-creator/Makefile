# 修改这里
VERSION ?= 0.0.1
# BUILD_TIME is used for daemon BuildTime in go build.
BUILD_TIME=$(shell date --rfc-3339 s 2> /dev/null | sed -e 's/ /T/')

IMG = serverless-rds-creator-image

fmt:
	go fmt ./...

build-img:
	docker build -t ${IMG}:latest -f Dockerfile .

run:
	docker run -p 9000:9000 --name myTerraform --rm -it ${IMG}:latest

push:
	docker tag ${IMG}:latest registry.cn-hangzhou.aliyuncs.com/serverless_devs/serverless-rds-creator-image:${VERSION}
	docker push registry.cn-hangzhou.aliyuncs.com/serverless_devs/serverless-rds-creator-image:${VERSION}

vendor:
	go mod vendor


build&push: build-img push
