ARG dist="/tmp/checkicp"
ARG projectDir="/icp"

FROM golang:1.16-alpine3.14 AS builder
RUN apk add build-base upx
ARG dist
ARG projectDir
WORKDIR ${projectDir}
COPY . .
RUN go build -trimpath -mod vendor -o icp cmd/main.go &&\
  upx -9 -o ${dist} icp


FROM alpine
ARG dist
COPY --from=builder ${dist} /usr/local/bin/checkicp