FROM golang:alpine

RUN mkdir -p /app

WORKDIR /app

ADD target/bin/order-sse /app

EXPOSE 8080

CMD ["./order-sse"]



#docker run -i -t --rm -v $PWD:/go/src/order-sse-subscribe-redis -w /go/src/order-sse-subscribe-redis golang go build

#docker build -t arthurstockler/order-sse .
#docker push arthurstockler/order-sse
#docker run --rm -p 8080:8080 arthurstockler/order-sse
