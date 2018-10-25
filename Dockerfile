FROM golang:alpine AS builder

ENV GOPATH="/go" CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go build -o icqag icqag.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/icqag .
EXPOSE 8888
CMD ["./icqag"]