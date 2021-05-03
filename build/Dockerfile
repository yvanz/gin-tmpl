FROM golang:1.14.13-alpine3.11  as builder

LABEL maintainer="yvan.zy@gmail.com"

ENV GO111MODULE=on

WORKDIR /go/src/gin-tmpl

COPY . .

RUN go install -mod vendor -ldflags="-s -w" -v ./...

FROM alpine:3.11

ENV MIRROR_URL=http://mirrors.aliyun.com/alpine/

RUN echo '' > /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.11/main" >> /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.11/community" >> /etc/apk/repositories

COPY --from=builder /go/bin/app /
COPY --from=builder /go/src/gin-tmpl/docs/ /docs
COPY --from=builder /go/src/gin-tmpl/configs/dev.yaml /

# xorm 指定时区时，需要 tzdata 文件，所以 apk add 后不能删除
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

EXPOSE 80

CMD ["/app", "--config", "dev.yaml"]