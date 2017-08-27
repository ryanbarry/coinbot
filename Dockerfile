FROM golang:1.7-alpine AS buildimg

WORKDIR /go/src/github.com/ryanbarry/coinbot

COPY . .

RUN go install

FROM alpine

RUN apk --no-cache add ca-certificates
COPY --from=buildimg /go/bin/coinbot /usr/bin/coinbot

CMD ["/usr/bin/coinbot"]
