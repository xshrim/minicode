#!/bin/sh

set -x

if [ "$HOST" ]; then
  sed -i "s/127.0.0.1/$HOST/g" /root/static/index.html
fi

if [ "$PORT" ]; then
  sed -i "s/22/$PORT/g" /root/static/index.html
fi

if [ "$USER" ]; then
  sed -i "s/root/$USER/g" /root/static/index.html
fi

if [ "$PASSWD" ]; then
  sed -i "s/admin/$PASSWD/g" /root/static/index.html
fi

/usr/sbin/sshd -D &

/root/webssh
