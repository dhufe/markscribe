# Build auf der nativen Plattform des Runners (amd64)
FROM golang:1.25.5-alpine3.23 AS build

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /go/markscribe

COPY . .
RUN apk update --no-cache \
    && apk add --no-cache make zip \
    && go mod tidy \
    && CGO_ENABLED=0 \
       GOOS=${TARGETOS} \
       GOARCH=${TARGETARCH} \
       go build -ldflags="-w -s" -o markscribe ./cmd/...

# Prod-Stage mit Zielplattform
FROM alpine:3.23 AS prod

LABEL org.opencontainers.image.source=https://gitlab.hufschlaeger.net/
LABEL org.opencontainers.image.description="markscribe"
LABEL org.opencontainers.image.licenses=MIT

# Nur das Binary kopieren (kein apk add nötig wenn QEMU fehlt)
WORKDIR /app
COPY --from=build /go/markscribe/markscribe /app/

ENTRYPOINT ["/app/markscribe"]
