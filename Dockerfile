# Based on https://github.com/chemidy/smallest-secured-golang-docker-image/blob/master/go_module/Dockerfile
FROM golang:1.14.3-alpine3.11 AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/github.com/manat/microproxy

COPY go.mod .

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod download
RUN go mod verify

COPY . .

# Build the binary
RUN cd $GOPATH/src/github.com/manat/microproxy/cmd/microproxy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/microproxy .


############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/microproxy /go/bin/
COPY --from=builder --chown=appuser /go/src/github.com/manat/microproxy/config.json /tmp/config.json

# Use an unprivileged user.
USER appuser:appuser

# Default Env
ENV MICROPROXY_FILEPATH=/tmp/config.json
ENV MICROPROXY_SERVER_PORT=1338

# Run the microproxy binary.
ENTRYPOINT ["/go/bin/microproxy"]
