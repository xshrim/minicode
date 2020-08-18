FROM alpine:latest

RUN mkdir /root/.ipfs && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 4001 5001 8080

WORKDIR /

COPY entrypoint.sh /

COPY swarm.key /root/.ipfs/

COPY ipfs /usr/bin/

CMD ["/entrypoint.sh"]
