FROM golang:1.20 as compiler
USER root
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn/,https://goproxy.io/,direct

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
# RUN apk add build-base
WORKDIR /starlite

COPY . .

#ARG MODULE_NAME
WORKDIR /starlite/internal/first

#-ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -s -w"
#RUN GOOS=linux GOARCH=amd64 go build -tags musl  -o main ./cmd/*
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/*
RUN mkdir -p dist && \
    cp main dist 
    # && \
    # if [ -d db ]; then cp -r db dist; fi
#RUN ls dist
RUN echo $(ls -1 dist)

FROM alpine:3.16.0

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add ca-certificates tzdata

#ARG MODULE_NAME
WORKDIR /starlite/first
RUN true

COPY --from=compiler /starlite/internal/first/dist .
RUN true

COPY internal/first/app.env .
COPY start.sh .

COPY internal/first/db/migration ./db/migration

EXPOSE 8080

#ENV MODULE_NAME=first
#USER 1000
#USER root
RUN chmod 777 /starlite/first/start.sh
#RUN ls /starlite/first
RUN echo $(ls -1 /starlite/first)
#ENTRYPOINT /starlite/first/main
ENTRYPOINT ["/starlite/first/start.sh"]
CMD [ "/starlite/first/main" ]
# ENTRYPOINT ["/starlite/first/start.sh"]
#CMD ["sleep infinity"]
#CMD ["sh", "-c", "tail -f /dev/null"]


# FROM golang:1.20 as builder

# USER root
# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn/,https://goproxy.io/,direct

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
# RUN apk add build-base
# WORKDIR /first

# COPY . .

# RUN go mod init myapp
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /app
# COPY --from=builder /app/main /app/

# EXPOSE 8080
# ENTRYPOINT ["/app/main"]