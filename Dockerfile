FROM ryanbarry/coinbot-dev AS buildimg

WORKDIR /go/src/github.com/ryanbarry/coinbot

COPY . .
RUN glide install
RUN go install

FROM alpine

RUN apk --no-cache add ca-certificates
COPY --from=buildimg /go/bin/coinbot /usr/bin/coinbot

CMD ["/usr/bin/coinbot"]