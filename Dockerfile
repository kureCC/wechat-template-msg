# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /wechat-template-msg

COPY . .

RUN go mod tidy

RUN go build -o wechat-template-msg

# Deploy
FROM alpine:3.15

RUN adduser -D nonroot

WORKDIR /wechat-template-msg

COPY --from=builder /wechat-template-msg /wechat-template-msg

USER nonroot:nonroot
