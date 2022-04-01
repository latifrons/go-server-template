FROM golang:1.16-alpine as builder

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN cat /etc/apk/repositories
RUN apk add --no-cache curl iotop
RUN apk add --no-cache make gcc musl-dev linux-headers git

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

WORKDIR /src
ADD ./go.mod /src/
RUN go mod download

ADD . /src/

RUN go build ./app/main.go

# Copy OG into basic alpine image
FROM alpine:latest

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN cat /etc/apk/repositories
RUN apk add --no-cache curl iotop tzdata

WORKDIR /www

COPY --from=builder src/nodedata/config ./nodedata/config/
COPY --from=builder src/nodedata/data ./nodedata/data/
COPY --from=builder src/main .

# COPY --from=builder src/nodedata/private_template/* ./nodedata/private/

EXPOSE 8080 8082 8084 8085

CMD ["./main", "run"]
