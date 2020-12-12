FROM docker.elastic.co/kibana/kibana:6.8.4
LABEL maintainer="xshrim@yeah.net"

COPY kibana/sentinl-v6.8.4.zip /usr/share/kibana/
RUN kibana-plugin install file:///usr/share/kibana/sentinl-v6.8.4.zip && rm -rf /usr/share/kibana/sentinl-v6.8.4.zip
