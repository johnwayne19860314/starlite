FROM 19860314/starlite-front:latest AS frontendBuilder
WORKDIR /starlite
COPY ui/ .

#WORKDIR /starlite/ui
RUN npm run build


FROM golang:1.20 as compiler
USER root
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn/,https://goproxy.io/,direct

WORKDIR /starlite

COPY . .

WORKDIR /starlite/internal/first

#-ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn -s -w"
#RUN GOOS=linux GOARCH=amd64 go build -tags musl  -o main ./cmd/*
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/*
RUN mkdir -p dist && \
    cp main dist 

RUN echo $(ls -1 dist)

FROM alpine:3.16.0

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add ca-certificates tzdata curl

#ARG MODULE_NAME
RUN addgroup -S deploy && adduser -S deploy -G deploy
ARG ROOT_DIR=/starlite/first
WORKDIR ${ROOT_DIR}
RUN chown deploy:deploy ${ROOT_DIR}
RUN true

COPY --from=frontendBuilder --chown=deploy:deploy /starlite/dist ./ui/dist
COPY --from=compiler --chown=deploy:deploy /starlite/internal/first/dist .
RUN true

#COPY internal/first/app.env .
COPY --chown=deploy:deploy start.sh .

COPY --chown=deploy:deploy internal/first/db/migration ./db/migration

EXPOSE 8080

#ENV MODULE_NAME=first
USER deploy
# RUN chmod 777 /starlite/first/start.sh
ENTRYPOINT /starlite/first/main
# ENTRYPOINT ["/starlite/first/start.sh"]
# CMD [ "/starlite/first/main" ]
# ENTRYPOINT ["/starlite/first/start.sh"]
#CMD ["sleep infinity"]
#CMD ["sh", "-c", "tail -f /dev/null"]

