# Build OG from alpine based golang environment
FROM golang:1.16-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

ADD . /src
WORKDIR /src
RUN go build ./app/main.go

# Copy OG into basic alpine image
FROM alpine:latest

RUN apk add --no-cache curl iotop busybox-extras

COPY --from=builder src/nodedata/config ./nodedata
COPY --from=builder src/build/main .

# for a temp running folder. This should be mounted from the outside
RUN mkdir /rw

EXPOSE 8080

WORKDIR /

CMD ["./main", "run"]
