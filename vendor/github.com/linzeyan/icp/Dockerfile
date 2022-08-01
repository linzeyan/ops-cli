FROM golang:1.17-alpine3.14 AS builder
RUN apk add build-base upx
WORKDIR /go/src/icp
COPY . .
RUN go mod vendor
RUN CGO_ENABLED=0 go build -trimpath -mod vendor -o icp cmd/main.go &&\
  upx -9 -o /go/bin/checkicp icp

FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=builder /go/bin/checkicp /bin/checkicp
ENTRYPOINT ["/bin/checkicp"]