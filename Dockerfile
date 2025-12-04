FROM docker.io/golang:1.25.5-alpine3.22 AS build
WORKDIR /go/markscribe
COPY . .
RUN go build

FROM alpine:3.22 AS prod

COPY --from=build /go/markscribe/markscribe /go/bin/markscribe
ENTRYPOINT ["/go/bin/markscribe"]