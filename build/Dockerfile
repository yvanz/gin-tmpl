FROM golang:1.18.10-alpine3.17 as builder

LABEL maintainer="yvan.zy@gmail.com"

WORKDIR /go/src/gin-demo

COPY . .

RUN apk add make git && \
    make install

FROM alpine:3.17.3

ENV MIRROR_URL=http://mirrors.aliyun.com/alpine/

RUN echo '' > /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.17/main" >> /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.17/community" >> /etc/apk/repositories

COPY --from=builder /go/bin/gin-demo /
COPY --from=builder /go/src/gin-demo/docs/ /docs
COPY --from=builder /go/src/gin-demo/configs/dev.yaml /

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

EXPOSE 9901

CMD ["/gin-demo", "--config", "dev.yaml"]