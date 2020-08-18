FROM alpine:latest

LABEL xshrim "xshrim@yeah.net"

WORKDIR /root/

ADD logbench /root/
ADD bench.sh /root/
ADD hyperfine /root/
ADD run.sh /root/
ADD bench.yaml /root/

CMD ["/root/logbench"]
