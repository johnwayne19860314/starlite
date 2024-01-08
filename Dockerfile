FROM golang:1.20 as compiler
USER root
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn/,https://goproxy.io/,direct

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
# RUN apk add build-base
WORKDIR /starlite

COPY . .

ARG MODULE_NAME
WORKDIR /starlite/internal/$MODULE_NAME

#-ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -s -w"
RUN GOOS=linux GOARCH=amd64 go build -tags musl  -o main ./cmd/*
RUN mkdir -p dist && \
    cp main dist 
    # && \
    # if [ -d db ]; then cp -r db dist; fi
#RUN ls dist
RUN echo $(ls -1 dist)

FROM alpine:3.16.0

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add ca-certificates tzdata

ARG MODULE_NAME
WORKDIR /starlite/$MODULE_NAME
RUN true

COPY --from=compiler /starlite/internal/$MODULE_NAME/dist .
RUN true

COPY internal/$MODULE_NAME/app.env .
COPY start.sh .

COPY internal/$MODULE_NAME/db/migration ./db/migration

EXPOSE 8080

ENV MODULE_NAME=$MODULE_NAME
#USER 1000
#USER root
RUN chmod 777 /starlite/first/start.sh
#RUN ls /starlite/$MODULE_NAME
RUN echo $(ls -1 /starlite/$MODULE_NAME)
#ENTRYPOINT /starlite/$MODULE_NAME/main
#ENTRYPOINT ["/starlite/$MODULE_NAME/start.sh"]
CMD [ "/starlite/first/main" ]
ENTRYPOINT ["/starlite/first/start.sh"]
#CMD ["sleep" "infinity"]


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