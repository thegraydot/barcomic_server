ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Stage 1: Build
FROM golang:1.24-bookworm AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

# Cache dependency download separately from source build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-s -w" \
    -o /barcomic \
    ./cmd/barcomic/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.21

COPY --from=builder /barcomic /barcomic

ENTRYPOINT ["/barcomic"]
CMD ["-v", "-a", "0.0.0.0", "-p", "80", "-s", "-i=false"]
