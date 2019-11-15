#!/bin/sh

#https://portal.caipin.pro/modules/servers/V2raySocks/osubscribe.php?sid=52418\&token=B2BmsbVngJQT

function delay() {
  rtt=$(\ping -c 2 -w 3 $1 |grep rtt |awk '{print $4}' |awk -F'/' '{print $2}')
  if [ -z $rtt ]; then
    echo 10000
  else
    echo ${rtt%.*}
  fi
}

function config() {

  ps=$(echo $1 | jq '.ps')

  echo "[${ps}]正在解析配置信息..."
  echo "$1"

  addr=$(echo $1 | jq '.add')

  port=$(echo $1 | jq '.port' | jq 'tonumber')

  aid=$(echo $1 | jq '.aid' | jq 'tonumber')

  host=$(echo $1 | jq '.host')

  id=$(echo $1| jq '.id')

  net=$(echo $1 | jq '.net')

  adr=$(echo $addr | sed 's/"//g')
  # 连接失败
  # if ! \ping -c 1 -w 3 $adr > /dev/null; then
  #   return 1
  # fi

  if [ $strategy == "fast" ];then
    printf "[${ps}]正在测试延迟数据<$2 -- "
    rtt=$(delay $adr)
    printf "$rtt>...\n"
    if [ $rtt -ge $2 ];then
      return $2
    fi
  else
    printf "[${ps}]正在测试延迟数据<"
    rtt=$(delay $adr)
    printf "$rtt>...\n"
    if [ $rtt -eq 10000 ];then
      return 1
    fi
  fi
  # if [ $(echo "$rtt > $2" | bc) -eq 1 ];then
  #   return $2
  # fi

  #jq ".outbounds[0].settings.vnext[0].address=$addr" config.json

  #jq ".outbounds=(.outbounds[] | select(.tag==\"proxy\").settings.vnext[0].address=$addr)" config.json > config.tmp && mv config.tmp config.json

  echo "[${ps}]正在更新配置信息..."

  jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].address=$addr)]" config.json >config.tmp && mv config.tmp config.json

  jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].port=$port)]" config.json >config.tmp && mv config.tmp config.json

  jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].users[0].id=$id)]" config.json >config.tmp && mv config.tmp config.json

  jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").settings.vnext[0].users[0].alterId=$aid)]" config.json >config.tmp && mv config.tmp config.json

  jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").streamSettings.network=$net)]" config.json >config.tmp && mv config.tmp config.json

  jq ".outbounds=[(.outbounds[] | select(.tag==\"$tag\").streamSettings.wsSettings.headers.Host=$host)]" config.json >config.tmp && mv config.tmp config.json

  cps=$ps

  if [ $strategy == "fast" ];then
    return $rtt
  else
    return 0
  fi
}

cps="已配置"

if [ ! -f config.json ];then
  echo "未找到config.json文件"
  exit 1
fi

if [ $# -gt 0 ];then

  sed -i '/\/\//d' config.json

  strategy="fast"  # fast
  link=$1
  tag=$2
  : ${tag:="proxy"}

  if ! type jq &>/dev/null;then
    echo "未发现jq工具"
    exit 1
  fi

  cps=$(\cat config.json | jq ".outbounds[] | select(.tag==\"$tag\").settings.vnext[0].address" | sed 's/"//g')

  if [[ $link == "https://"* ]];then
    echo "[HTTPS]正在获取订阅信息..."
    data=$(curl -k -s "$1" | base64 -d)
  elif [[ $link == "vmess://"* ]];then
    echo "[VMESS]正在获取订阅信息..."
    data=($1)
  fi

  if [ $strategy == "fast" ];then
    minrtt=$(delay $cps)
  fi
  for vmess in ${data[@]}
  do
    cfg=$(echo $vmess | sed "s/vmess:\/\///g" | base64 -d)
    if [ $strategy == "fast" ];then
      config "$cfg" "$minrtt"
      minrtt=$?
    else
      if config "$cfg";then
        break
      fi
    fi
  done
fi

echo "[${cps}]正在启动v2ray服务..."
echo "============================================="
./v2ray --config config.json
