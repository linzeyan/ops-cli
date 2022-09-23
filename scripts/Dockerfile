FROM golang:1.18-alpine as builder

ARG Tag
ARG Date
ARG Commit

COPY / /ops-cli
WORKDIR /ops-cli

RUN apk --no-cache add gcc

RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X github.com/linzeyan/ops-cli/cmd.appVersion=$Tag -X github.com/linzeyan/ops-cli/cmd.appBuildTime=$Date -X github.com/linzeyan/ops-cli/cmd.appCommit=$Commit" -o ops-cli main.go


FROM alpine

COPY --from=builder /ops-cli/ops-cli /usr/bin/
# ENTRYPOINT [ "/usr/bin/ops-cli" ]
CMD [ "/usr/bin/ops-cli" ]
