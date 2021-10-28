FROM golang:1.17.2-alpine3.14 AS builder

ENV GOPATH="/go" CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -ldflags "-w -s -linkmode internal" -o icqag icqag.go

FROM alpine:3.14
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/icqag .
EXPOSE 8888
CMD ["./icqag"]