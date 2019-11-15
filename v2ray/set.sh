#!/bin/sh

curl -k https://portal.caipin.pro/modules/servers/V2raySocks/osubscribe.php?sid=52418\&token=B2BmsbVngJQT | base64 -d

sed -i '/\/\//d' config.json

jsdata=$(echo $1 | sed "s/vmess:\/\///g" | base64 -d)

tag=$2
: ${tag:="proxy"}

echo "正在解析配置信息..."

addr=$(echo $jsdata | jq '.add')

port=$(echo $jsdata | jq '.port' | jq 'tonumber')

aid=$(echo $jsdata | jq '.aid' | jq 'tonumber')

host=$(echo $jsdata | jq '.host')

id=$(echo $jsdata | jq '.id')

net=$(echo $jsdata | jq '.net')

#jq ".outbounds[0].settings.vnext[0].address=$addr" config.json

#jq ".outbounds=(.outbounds[] | select(.tag==\"proxy\").settings.vnext[0].address=$addr)" config.json > config.tmp && mv config.tmp config.json

echo "正在更新配置信息..."

jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].address=$addr)]" config.json >config.tmp && mv config.tmp config.json

jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].port=$port)]" config.json >config.tmp && mv config.tmp config.json

jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].users[0].id=$id)]" config.json >config.tmp && mv config.tmp config.json

jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].users[0].alterId=$aid)]" config.json >config.tmp && mv config.tmp config.json

jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").streamSettings.network=$net)]" config.json >config.tmp && mv config.tmp config.json

jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").streamSettings.wsSettings.headers.Host=$host)]" config.json >config.tmp && mv config.tmp config.json

echo "正在启动v2ray服务..."
./v2ray --config config.json
