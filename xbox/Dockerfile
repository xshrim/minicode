FROM busybox:1.28.4

MAINTAINER xshrim <xshrim@gmail.com>

WORKDIR /home

ADD tools/* /bin/

ADD curl-7.30.0.ermine.tar.bz2 .

RUN mv /home/curl-7.30.0.ermine/curl.ermine /bin/curl && rm -Rf /home/curl-7.30.0.ermine

CMD ["/bin/sh"]
