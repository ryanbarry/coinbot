FROM ryanbarry/coinbot-dev AS buildimg

COPY . .
RUN glide install
RUN go install

FROM alpine

WORKDIR /

RUN apk --no-cache add ca-certificates
COPY --from=buildimg /go/bin/coinbot /usr/bin/coinbot

CMD ["/usr/bin/coinbot"]
