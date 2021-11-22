all: kbuildresource

TAG = 1.0.0
PREFIX = docker-hub-green.tools.bryson.com/littleproject/kbuildresource

.PHONY: vet fmt build clean

vet:
    go vet ./...

fmt:
    go fmt ./...

kbuildresource: *.go controller/*.go routers/*.go models/*.go vet fmt
    CGO_ENABLED=0 GOOS=linux go build -a installsuffix cgo -ldflags '-w' -o $@

build:
    docker build -t $(PREFIX):$(TAG) .

clean:
    rm -f kbuildresource