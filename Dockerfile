FROM golang AS builder
ADD . /go/smtp2http/
WORKDIR /go/smtp2http
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/smtp2http

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /app
COPY --from=builder /go/bin/smtp2http .
RUN chmod +x /app/smtp2http
EXPOSE 25

ENTRYPOINT ["/app/smtp2http"]