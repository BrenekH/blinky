FROM golang:1.17-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/blinky

COPY go.mod ./go.mod
COPY go.sum ./go.sum
# CGO must also be disabled while downloading so BadgerDB doesn't download CGO libs.
RUN CGO_ENABLED=0 go mod download

COPY . .

# Disable CGO so that we have a static binary and set the platform for multi-arch builds.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o blinkyd -ldflags="-s -w -X 'github.com/BrenekH/blinky/vars.Version=${LDFLAGS_VERSION}'" cmd/blinkyd/main.go

# Run stage
FROM archlinux

ENV BLINKY_SIGNING_KEY=/sign-db.key \
	BLINKY_CONFIG_DIR=/config \
	BLINKY_GNUPG_DIR=/gnupg

WORKDIR /usr/src/app

COPY --from=builder /go/src/blinky/blinkyd ./blinkyd

ENTRYPOINT ["./blinkyd"]
