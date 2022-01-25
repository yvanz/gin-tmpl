FROM golang:1.16.9-alpine3.14 as builder

LABEL maintainer="yvan.zy@gmail.com"

ENV GO111MODULE=on

WORKDIR /go/src/gin-tmpl

COPY . .

RUN go install -mod vendor -ldflags="-s -w" -v ./...

FROM alpine:3.14

ENV MIRROR_URL=http://mirrors.aliyun.com/alpine/

RUN echo '' > /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.14/main" >> /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.14/community" >> /etc/apk/repositories

COPY --from=builder /go/bin/app /
COPY --from=builder /go/src/gin-tmpl/docs/ /docs
COPY --from=builder /go/src/gin-tmpl/configs/dev.yaml /

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

EXPOSE 80

CMD ["/app", "--config", "dev.yaml"]