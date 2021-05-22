#!/bin/sh

function require() {
  if [ -f $1 ]; then
   chmod a+x $1
   echo ./$1
  else
    if type $1 &>/dev/null; then
      echo $(which $1)
    else
      echo "ERROR: can not find $1 file"
      exit 1
    fi
  fi
}

function chk() {
  # require snapshot.db
  require etcd
  require etcdctl
}

function getEndpoint() {
  local ep="localhost:2379"

  if [[ $1 =~ ":" ]]; then
    ep=$1
  else
    if [ "$1" -gt 0 ] 2>/dev/null ;then
      ep=127.0.0.1:$1
    elif [ -n "$1" ]; then
      ep=$1:2379
    fi
  fi

  if [ ${ep:0:1} = ":" ]; then
    ep="0.0.0.0"$ep
  fi

  echo $ep
}

function eps() {
  if [ "$1" == "mirror" ]; then
    ps -ef|grep 'etcdctl make-mirror'|grep -v 'grep'
  elif [ "$1" == "etcd" ]; then
    ps -ef|grep 'etcd '|grep -v 'etcdctl'|-v grep 'etcdcc'|grep -v 'grep'
  else
    ps -ef|grep 'etcd '|grep -v 'etcdcc'|grep -v 'grep'
  fi
}

function save() {
  local ep=$(getEndpoint $1)
  local snapshot=$2
  if [ -z "$snapshot" ]; then
    snapshot="snapshot.db"
  fi
  
  local etcdctl=$(require etcdctl)

  $etcdctl --endpoints=$ep snapshot save $snapshot 
}

function list() {
  local ep=$(getEndpoint $1)
  local etcdctl=$(require etcdctl)

  local host=${ep%%:*}
  local port=${ep##*:}
  if [ "$host" == "0.0.0.0" ]; then
    host="127.0.0.1"
  fi
  ep=$host:$port

  $etcdctl --endpoints=$ep endpoint status --cluster -w table
}

function switch() {
  local ep=$(getEndpoint $1)
  local id=$2

  local etcdctl=$(require etcdctl)

  if [ -z "$id" ]; then
    list $1
    read -p  "Input New Leader ID/No.:" id
  fi

  if [ -z "$id" ]; then
    id=1
  fi
  
  if [ "$id" -gt 0 ] 2>/dev/null ;then
    id=$($etcdctl --endpoints=$ep member list | awk '{print $1}' | sed 's/,//g' | sed -n "${id}p")
  fi

  $etcdctl --endpoints=$ep move-leader $id
}

function keys() {
  local ep=$(getEndpoint $1)
  local kwd=$2

  local etcdctl=$(require etcdctl)

  if [ -n "$kwd" ]; then
    $etcdctl --endpoints=$ep get / --prefix --keys-only|grep -v '^$'|grep "$kwd"
  else
    $etcdctl --endpoints=$ep get / --prefix --keys-only|grep -v '^$'
  fi
}

function rev() {
  local ep=$(getEndpoint $1)
  local etcdctl=$(require etcdctl)

  $etcdctl --endpoints=$ep endpoint status --write-out="json" | egrep -o '"revision":[0-9]*' | egrep -o '[0-9].*'
}

function get() {
  local ep=$(getEndpoint $1)
  local etcdctl=$(require etcdctl)

  shift

  $etcdctl --endpoints=$ep get $@
}

function put() {
  local ep=$(getEndpoint $1)
  local etcdctl=$(require etcdctl)

  shift

  $etcdctl --endpoints=$ep put $@
}

function del() {
  local ep=$(getEndpoint $1)
  local etcdctl=$(require etcdctl)

  shift

  $etcdctl --endpoints=$ep del $@
}

function up() {
  local name
  local ep
  local snapshot
  if [ "$1" -gt 0 ] 2>/dev/null ;then
    name='etcd-'$(head /dev/urandom |cksum |md5sum |cut -c 1-9)
    ep=$(getEndpoint $1)
    snapshot=$2
  else
    name=$1
    ep=$(getEndpoint $2)
    snapshot=$3
  fi

  local host=${ep%%:*}
  local port=${ep##*:}

  if [ -d $name ]; then
    echo "ERROR: cluster $name exists already"
    exit 1
  fi
  
  if [ -n "$snapshot" ]; then
    snapshot=$(require $snapshot)
    local etcdctl=$(require etcdctl)
    $etcdctl snapshot restore $snapshot --name $name --data-dir $name --initial-cluster $name=http://$host:$((port+1)) --initial-advertise-peer-urls http://$host:$((port+1))
    if [ $? -ne 0 ]; then
      echo "ERROR: snapshot restore failed"
      rm -rf $name
      exit 1
    fi
  fi

  local etcd=$(require etcd)
  $etcd --name $name --data-dir $name --listen-client-urls http://$host:$port  --advertise-client-urls http://$host:$port --listen-peer-urls http://$host:$((port+1)) --initial-cluster-state new &

  if [ $? -ne 0 ]; then
    echo "launch cluster $name failed"
    rm -rf $name
    exit 1
  fi

  sleep 1

  disp $ep
}

function down() {
  dirs=$(ps -ef | grep 'etcd --name'|grep ":$1"|awk '{print $12}')
  kill -9 $(ps -ef|grep 'etcd --name'|grep ":$1" |awk '{print $2}') &> /dev/null
  
  for elem in ${dirs[@]}; do
    rm -rf $elem
  done
}

function sync() {
  if [ $# -ne 2 ]; then
    echo "ERROR: need src cluster and dst cluster"
    exit 1
  fi

  local src=$1
  local dst=$2
  local src=$(getEndpoint $1)
  local dst=$(getEndpoint $2)

  if [ "$src" == "$dst" ]; then
    echo "ERROR: sync dest endpoint shoud not be source endpoint"
    exit 1
  fi

  local etcdctl=$(require etcdctl)
  $etcdctl make-mirror --endpoints=$src $dst > /dev/null &

  trap "clean mirror" HUP INT PIPE QUIT TERM

  while true
  do
    clear
    disp $dst
    sleep 2
  done 
}


function disp() {
  local ep=$(getEndpoint $1)
  local host=${ep%%:*}
  local port=${ep##*:}

  if [ "$host" == "0.0.0.0" ]; then
    host="127.0.0.1"
  fi
  ep=$host:$port

  local etcdctl=$(require etcdctl)

  local num=$(printf %09d $($etcdctl --endpoints=http://$ep get / --prefix --keys-only|grep -v '^$'|wc -l))
  echo "======================== Kubernetes Resource Object Total: $num ========================="
  $etcdctl --endpoints=$ep endpoint status -w table
}

function join() {
  local cep=$(getEndpoint $1)
  local ep
  local name
  local chost
  if [ "$2" -gt 0 ] 2>/dev/null ;then
    name='etcd-'$(head /dev/urandom |cksum |md5sum |cut -c 1-9)
    ep=$2
    host=$3
  else
    name=$2
    ep=$3
    host=$4
  fi
  local host=${ep%%:*}
  local port=${ep##*:}

  chost=$host
  if [ "$chost" == "0.0.0.0" ]; then
    chost="127.0.0.1"
  fi

  if [ -d $name ]; then
    echo "ERROR: cluster $name exists already"
    exit 1
  fi

  local etcdctl=$(require etcdctl)
  local envs=$($etcdctl --endpoints=$cep member add $name --peer-urls=http://$chost:$((port+1))|grep ETCD_)

  echo $envs

  local etcd=$(require etcd)
  eval $envs;$etcd --name $name --data-dir $name --listen-client-urls http://$host:$port  --advertise-client-urls http://$host:$port --listen-peer-urls http://$host:$((port+1)) --initial-cluster-state existing &
  
  if [ $? -ne 0 ]; then
    echo "launch cluster $name failed"
    rm -rf $name
    exit 1
  fi

  sleep 1
  list $ep
}

function clean() {
  if [ "$1" == "mirror" ]; then
    echo aa
    kill -9 $(ps -ef|grep 'etcdctl make-mirror'|grep -v 'grep' |awk '{print $2}') &> /dev/null
  else
    dirs=$(ps -ef | grep 'etcd --name'|grep -v 'grep'|awk '{print $12}')
    kill -9 $(ps -ef|grep 'etcdctl make-mirror'|grep -v 'grep' |awk '{print $2}') &> /dev/null
    kill -9 $(ps -ef|grep 'etcd --name'|grep -v 'grep' |awk '{print $2}') &> /dev/null
    for elem in ${dirs[@]}; do
      rm -rf $elem
    done
  fi
  exit 0
}


function gc() {
  local ep=$(getEndpoint $1)

  local etcdctl=$(require etcdctl)

  rev=$($etcdctl --endpoints=$ep endpoint status --write-out="json" | egrep -o '"revision":[0-9]*' | egrep -o '[0-9].*')
  $etcdctl --endpoints=$ep compact $rev
  $etcdctl --endpoints=$ep defrag

  disp $ep
}

function usage() {
  echo '
  **USAGE**:
  # 打印命令帮助
  etcdcc help             

  # 检查必要文件
  etcdcc chk

  # 备份集群快照
  etcdcc save <host>:<port> <snapshot>

  # 切换集群leader
  etcdcc switch <host>:<port> <etcdid>

  # 查看集群节点
  etcdcc list <host>:<port>
  
  # 查看集群key列表
  etcdcc keys <host>:<port> <keyword>

  # 获取集群revision
  etcdcc rev <host>:<port>

  # 读取指定key
  etcdcc get <host>:<port> <key>

  # 写入指定key
  etcdcc put <host>:<port> <key> <value>

  # 删除指定key
  etcdcc del <host>:<port> <key>

  # 启动新的etcd实例
  etcdcc up <name> <port>

  # 从快照恢复实例
  etcdcc up <name> <port> <snapshot>

  # 关闭etcd实例
  etcdcc down <port>

  # 启动集群间镜像同步
  etcdcc sync <src-host>:<src-port> <dst-host>:<dst-port>

  # 显示etcd进程
  etcdcc ps

  # 查看集群状态
  etcdcc disp <host>:<port>

  # 对集群进行数据压缩和碎片整理
  etcdcc gc <host>:<port>

  # 清除所有变更
  etcdcc clean
'
}

# https://www.cnblogs.com/lowezheng/p/10307592.html
case $1 in
"chk")
  chk
  ;;
"save")
  shift
  save $@
  ;;
"switch")
  shift
  switch $@
  ;;
"list")
  shift
  list $@
  ;;
"keys")
  shift
  keys $@
  ;;
"rev")
  shift
  rev $@
  ;;
"get")
  shift
  get $@
  ;;
"put")
  shift
  put $@
  ;;
"del")
  shift
  del $@
  ;;
#"join")
#  shift
#  join $@
#  ;;
"up")
  shift
  up $@
  ;;
"sync")
  shift
  sync $@
  ;;
"ps")
  shift
  eps $@
  ;;
"disp")
  shift
  disp $@
  ;;
"gc")
  shift
  gc $@
  ;;
"down")
  shift
  down $@
  ;;
"clean")
  shift
  clean
  ;;
*)
  usage
  ;;
esac

