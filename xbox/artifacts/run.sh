#!/bin/sh

webhost=""
webport="80"

kubexpport="2666"
kubexpsport="4666"

polarisport="2555"
polarispath="/tmp/"

if [ "$WEBHOST" ]; then
  webhost="$WEBHOST"
fi

if [ "$WEBPORT" ]; then
  webport="$WEBPORT"
fi

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

if [ "$TIMEOUT" ]; then
  sed -i "s/300/$TIMEOUT/g" /root/static/index.html
fi

if [ "$SSL" == "true" ]; then
  sed -i "s#ws://#wss://#g" /root/static/index.html
fi

hsdir="/root/"
if [ "$HSDIR" ]; then
  hsdir=$HSDIR
fi

if [ "$HSUSER" ] && [ "$HSPASSWD" ]; then
  hsauth="--auth-type http --auth-http $HSUSER:$HSPASSWD"
fi

if [ "$POLARISPORT" ]; then
  polarisport="$POLARISPORT"
fi

if [ "$POLARISPATH" ]; then
  polarispath="$POLARISPATH"
fi

if [ "$KUBEXPPORT" ]; then
  kubexpport="$KUBEXPPORT"
fi

if [ "$KUBEXPSPORT" ]; then
  kubexpsport="$KUBEXPSPORT"
fi

trap "echo SIGTERM/SIGINT detected; exit 1" SIGTERM SIGINT

set -x

if [ "$SSH" == "true" ]; then
  /usr/sbin/sshd -D &
fi

if [ "$SAMBA" == "true" ]; then
  smbd --no-process-group --configfile /root/smb.conf &
fi

if [ "$GOFS" == "true" ]; then
  WEBHOST=$webhost WEBPORT=$webport gofs -d $hsdir &
fi

if [ "$HS" == "true" ]; then
  gohttpserver -r $hsdir --port 2444 $hsauth --upload --delete --xheaders --cors --theme green --google-tracker-id "" &
fi

if [ "$WEBSSH" == "true" ]; then
  /root/webssh &
fi

if [ "$POLARIS" == "true" ]; then
  polaris dashboard --port $polarisport --audit-path $polarispath &
fi

if [ "$KUBEXP" == "true" ]; then
  kubexp --http-listen-port $kubexpport --https-listen-port $kubexpsport &
fi

tail -f /dev/null
