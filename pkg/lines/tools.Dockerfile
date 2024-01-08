FROM china-devops-docker-local.arf.xxx.cn/base-images/golang:1.16.3-alpine3.12 AS builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn/,https://goproxy.io/,direct
RUN go get golang.org/x/tools/cmd/goimports
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add build-base git
WORKDIR /app
ADD . .
RUN cd /app/tools/pen && make install
RUN cd /app/tools/mergeopenapi && make install
RUN cd /app/tools/importformat && make install

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add ca-certificates tzdata curl busybox make
