ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM registry.yygu.cn/library/golang:1.24-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /build/singer ./backend/cmd/singer/main.go

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /build/skill-proxy ./backend/cmd/skill-proxy/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS singer

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/singer /singer

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/singer"]

FROM gcr.io/distroless/static-debian12:nonroot AS skill-proxy

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/skill-proxy /skill-proxy

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/skill-proxy"]