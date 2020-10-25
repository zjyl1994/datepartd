FROM golang:1.15.3-alpine3.12 AS builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add build-base
COPY . /code
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
RUN mkdir -p /usr/local/go/src/github.com/zjyl1994 && \
    ln -s /code /usr/local/go/src/github.com/zjyl1994/datepartd && \
    cd /usr/local/go/src/github.com/zjyl1994/datepartd && \
    CGO_ENABLED=1 go build -a -o datepartd
FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
COPY --from=builder /usr/local/go/src/github.com/zjyl1994/datepartd/datepartd /app/datepartd
COPY --from=builder /usr/local/go/src/github.com/zjyl1994/datepartd/config.toml /app/config.toml
CMD ["/app/datepartd"]