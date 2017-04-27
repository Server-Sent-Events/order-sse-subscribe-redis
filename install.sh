#!/usr/bin/env bash
export GOPATH="$PWD"

go get -v -d || {
echo '### ERROR go get ###'
exit 1
}

CGO_ENABLED=0 GOOS=linux go build -v -ldflags '-s' -a -installsuffix cgo \
-o target/bin/order-sse . || {
echo '### ERROR build ###'
exit 1
}

docker rmi -f arthurstockler/order-sse
docker build -t arthurstockler/order-sse .
docker push arthurstockler/order-sse
#docker-compose up
