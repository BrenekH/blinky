FROM golang:1.17-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/blinky

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY . .

# Disable CGO so that we have a static binary and set the platform for multi-arch builds.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o blinkyd cmd/blinkyd/main.go

# Run stage
FROM archlinux

ENV BLINKY_SIGNING_KEY=/sign-db.key \
	BLINKY_CONFIG_DIR=/config \
	BLINKY_GNUPG_DIR=/gnupg

WORKDIR /usr/src/app

COPY --from=builder /go/src/blinky/blinkyd ./blinkyd

ENTRYPOINT ["./blinkyd"]
