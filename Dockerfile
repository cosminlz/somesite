FROM golang:1.13.4 as builder
ENV DATA_DIRECTORY /go/src/cabhelp.ro/backend
WORKDIR $DATA_DIRECTORY
ARG APP_VERSION
ARG CGO_ENABLED=0
COPY . .
RUN go build -ldflags="-X cabhelp.ro/backend/internal/config.Version=$APP_VERSION" cabhelp.ro/backend/cmd/server

FROM alpine:3.10
ENV DATA_DIRECTORY /go/src/cabhelp.ro/backend
RUN apk add --update --no-cache \
    ca-certificates
COPY --from=builder $DATA_DIRECTORY/server /backend
ENTRYPOINT ["/backend"]


