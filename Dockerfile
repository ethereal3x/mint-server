FROM golang:1.25-alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o mint-server ./cmd/

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/mint-server /usr/local/bin/
# config.yaml 从 Consul KV 拉取，不进镜像
WORKDIR /etc/mint-server
EXPOSE 8888 9999
CMD ["mint-server"]
