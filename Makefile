PROJECT?=faceit
PATH_ROJECT?=github.com/nakiner/faceit
#PATH_ROJECT?=${PROJECT}
APP?=bin/${PROJECT}

VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: proto
proto:
	sh ./scripts/protoc-gen.sh

.PHONY: clean
clean:
	rm -f ${APP}

.PHONY: build
build: clean
	cd ./cmd/app && \
  	CGO_ENABLED=0 go build -ldflags "-s -w \
    	-X ${PATH_ROJECT}/pkg/health.Version=${VERSION} \
    	-X ${PATH_ROJECT}/pkg/health.Commit=${COMMIT} \
    	-X ${PATH_ROJECT}/pkg/health.BuildTime=${BUILD_TIME}" \
    	-a -installsuffix cgo -o ../../${APP} .

.PHONY: run
run: build
	./${APP}

.PHONY: test
test:
	go test -v -race ./...

docker:
	docker build -t registry.devous.ru/faceit/app:latest .

docker-run:
	docker-compose up -d