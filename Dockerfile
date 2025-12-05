FROM golang:alpine3.20 AS build
WORKDIR /go/markscribe
COPY . .
RUN go build

FROM alpine:3.23 AS prod

COPY --from=build /go/markscribe/markscribe /go/bin/markscribe
ENTRYPOINT ["/go/bin/markscribe"]