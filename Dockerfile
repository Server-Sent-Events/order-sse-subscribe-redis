FROM golang:alpine

WORKDIR /app

ADD order-sse-subscribe-redis /app/

ENTRYPOINT ["./order-sse-subscribe-redis"]

EXPOSE 8080




#docker run -i -t --rm -v $PWD:/go/src/order-sse-subscribe-redis -w /go/src/order-sse-subscribe-redis golang go build

#docker build -t arthurstockler/order-sse:latest .
#docker push arthurstockler/order-sse
#docker run --rm -p 8080:8080 arthurstockler/order-sse
