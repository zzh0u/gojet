FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 构建应用（剥离符号表与 DWARF 信息以减小二进制体积）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -o main .

FROM alpine:3.20

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates


WORKDIR /root/

# 从构建阶段复制二进制与配置文件
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

EXPOSE 8080

# 启动二进制
CMD ["./main"]