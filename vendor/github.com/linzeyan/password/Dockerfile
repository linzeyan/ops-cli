ARG dist="/tmp/password"
ARG projectDir="/password"

FROM golang:1.16-alpine3.14 AS builder
RUN apk add build-base upx
ARG dist
ARG projectDir
WORKDIR ${projectDir}
COPY . .
RUN go build -trimpath -o main cmd/main.go
RUN upx -9 -o ${dist} main

FROM gcr.io/distroless/static:nonroot
COPY --from=builder --chown=nonroot:nonroot /tmp/password /bin/password
ENTRYPOINT ["/bin/password"]