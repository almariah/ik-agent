MAKEDIR:=$(strip $(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))

.PHONY: build
build: docker-build build-image

.PHONY: build-agent
build-agent:
	CGO_ENABLED=1 go build -ldflags '-w -s' -a -installsuffix cgo -o .build/ik-agent .

.PHONY: docker-build
docker-build:
	docker run -v ${GOPATH}/src:/go/src -v ${MAKEDIR}:/go/src/github.com/almariah/ik-agent golang:1.8 make -C /go/src/github.com/almariah/ik-agent build-agent

.PHONY: build-image
build-image:
	docker build -t ik-agent -f Dockerfile .

.PHONY: clean
clean:
	rm -rf .build
