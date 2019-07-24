FROM busybox:latest

MAINTAINER xshrim <xshrim@gmail.com>

WORKDIR /home

ADD curl-7.30.0.ermine.tar.bz2 .

ADD jq /bin/

ADD bat /bin/

ADD websocketd /bin/

RUN mv /home/curl-7.30.0.ermine/curl.ermine /bin/curl && rm -Rf /home/curl-7.30.0.ermine

CMD ["/bin/sh"]
