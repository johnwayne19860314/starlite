FROM docker-hub-remote.arf.xxx.cn/golangci/golangci-lint:v1.46.1

ENV GOPROXY=https://goproxy.cn/,https://goproxy.io/,direct
ENV GOPRIVATE="*.xxx.cn,stash.xxxmotors.com"

WORKDIR $GOPATH/src/github.startlite.cn/itapp/startlite/pkg/lines
RUN apt-get install -y make
ADD . .
