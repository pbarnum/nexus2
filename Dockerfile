FROM golang:1.18.2-alpine3.15 AS builder

LABEL author="Team MSRebirth" maintainer="admin@msrebirth.net"

LABEL org.opencontainers.image.source="https://github.com/MSRevive/nexus2"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.description="Docker image containing nexus2 binaries."

ARG build=develop

ENV CGO_ENABLED=1

WORKDIR /app

RUN go install github.com/silenceper/gowatch@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags "-X main.version=${VERSION}" -o /app/cmd/api/nexus2_api /app/cmd/api

CMD ["/bin/sh"]


FROM alpine:3.15 AS release

ENV PATH="/app/bin:${PATH}"

WORKDIR /app

COPY --from=builder /app/cmd/api/nexus2_api /app/bin/

CMD ["nexus2_api"]
