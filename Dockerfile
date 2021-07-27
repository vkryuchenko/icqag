FROM golang:alpine AS builder

ENV GOPATH="/go" CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod tidy
RUN go build -ldflags "-w -s -linkmode internal" -o icqag icqag.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/icqag .
EXPOSE 8888
CMD ["./icqag"]