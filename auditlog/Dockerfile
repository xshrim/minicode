FROM golang:alpine as builder
WORKDIR /root/
COPY ./ ./
RUN echo "http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
  && echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
  && echo "http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
  && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
  && sed -i 's/nl.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
  && apk update \
  && apk upgrade \
  && apk add --no-cache gcc musl-dev && GOOS=linux go build -mod vendor -a -o auditlog ./main.go
# RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -o auditlog ./main.go

FROM alpine:latest
LABEL maintainer="xshrim@yeah.net"
WORKDIR /root/
COPY --from=builder /root/auditlog .
EXPOSE 9090
CMD ["./auditlog"]