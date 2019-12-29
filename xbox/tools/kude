#!/usr/bin/sh

kubecfg=$HOME/.kube/config
image=xo
port=16022
timeout=60
user=root
idx=0
context=""
chroot=""
vol=""

mandatory=()

function usage() {
echo "USAGE:
  kude [<options>] podname

OPTIONS:
  -c [test] kubernetes连接上下文(可选)
  -d [log, logf, inspect, describe, wide, yaml, json] 非交互式指令(可选)
  -f [./kube/config] kubeonfig文件路径(可选)
  -i [0] 匹配pod内容器index(可选)
  -m [xo] 自定义工具容器(可选)
  -n [default] 查询的kubernetes命名空间(可选)
  -o [xos] podname模糊匹配(可选)
  -p [16022] 目标主机ssh端口(可选)
  -s [admin123] 目标主机ssh密码(可选)
  -t [60] 目标主机ssh连接超时时间(可选)
  -u [root] 目标主机ssh账号(可选)
  -r Debug容器中是否需要chroot到目标容器的文件系统
  -v Debug容器中是否需要挂载目标主机的docker.sock文件
  -h 使用说明

EXAMPLES:
# 命令在带kubectl的主机运行
# 不带ssh相关参数的命令默认已经配置了免密登录

# 最简, 所有命名空间中寻找名称带xos字段的pod, 部署Debug容器共享该pod下第一个容器的命名空间并进入交互模式
kude xos

# xshrim命名空间中寻找名称带xos字段的pod, 部署Debug容器共享该pod下第一个容器的命名空间并进入交互模式
kude -n xshrin xos

# xshrim命名空间中寻找名称带xos字段的pod, 部署Debug容器共享该pod下第一个容器的命名空间进入交互模式并chroot到该容器的文件系统
kude -n xshrim -i 0 -r xos

# 所有命名空间中寻找名称带xos字段的pod下的第二个容器并显示容器日志
kude -d log -i 1 xos

# 所有命名空间中寻找名称带xos字段的pod并显示其描述信息
kude -d describe xos

# 所有命名空间中寻找名称带xos字段的pod并将其配置以json格式输出
kude -d json xos

# xshrim命名空间中寻找名称带xos字段的pod, 部署Debug容器共享该pod下第一个容器的命名空间并进入交互模式, 目标主机ssh账号密码端口指定, Debug容器使用netshoot
kude -n xshrim -u root -p 22 -s admin123 -m netshoot xos
"
}

while [ $# -gt 0 ] && [ "$1" != "--" ]; do
  while getopts "c:d:f:i:m:n:o:p:s:t:u:hrv" opt; do
    case $opt in
      c) 
        context="--context $OPTARG"
        ;;
      d) 
        darg=$OPTARG
        image=""
        ;;
      f) 
        kubeconfig=$OPTARG
        ;;
      i)
        idx=$OPTARG
        ;;
      m) 
        image=$OPTARG
        ;;
      n) 
        namespace=$OPTARG
        ;;
      o) 
        pod=$OPTARG
        ;;
      p) 
        port=$OPTARG
        ;;
      s) 
        passwd=$OPTARG
        ;;
      t) 
        timeout=$OPTARG
        ;;
      u) 
        user=$OPTARG
        ;;
      r)
        chroot="-c 'chroot /proc/1/root/'"
        ;;
      v)
        vol="-v /var/run/docker.sock:/var/run/docker.sock"
        ;;
      h)
        usage
        exit 0
        ;;
      \?) 
        echo "Invalid option: -$OPTARG" >&2; 
        exit 1
        ;;
    esac
  done
  shift $((OPTIND-1))

  while [ $# -gt 0 ] && ! [[ "$1" =~ ^- ]]; do
    mandatory=("${mandatory[@]}" "$1")
    shift
  done
done

if [ "$1" == "--" ]; then
  shift
  mandatory=("${mandatory[@]}" "$@")
fi

if [ -z "$pod" ] && [ ${#mandatory[@]} -gt 0 ];
then
  pod=${mandatory[0]}
fi

if [ -n "$kubecfg" ] && [ -f "$HOME/.kube/config" ];
then
  kubecfg=$HOME/.kube/config
fi

if ! type docker &>/dev/null || ! type kubectl &>/dev/null;
then
  echo "docker and kubectl required"
  exit 1
fi

if [ -z "$namespace" ];
then
  ns="--all-namespaces"
else
  ns="-n $namespace"
fi

kc=""
if [ -n "$kubecfg" ];
then
  kc="--kubeconfig $kubecfg"
fi

pods=$(kubectl $kc $context get pod $ns 2>/dev/null | grep $pod)
if [ $? -ne 0 ] || [ ${#pods[@]} -eq 0 ] ;
then 
  echo "No resources found"
  exit 1
fi

if [ ${#pods[@]} -eq 1 ];
then
  podname=$(echo $pods | awk '{print $2}')
  namespace=$(echo $pods | awk '{print $1}')
  ns="-n $namespace"
else
  echo "Warning: more than one pod! selected the first one."
  podname=$(echo $pods[0] | awk '{print $2}')
  namespace=$(echo $pods[0] | awk '{print $1}')
  ns="-n $namespace"
fi

podinfo=$(kubectl $kc $context get pod $podname $ns -o=custom-columns=NODE:.spec.nodeName,HOST:.status.hostIP,NAME:.metadata.name,ADDR:.status.podIP,ID:.status.containerStatuses[$idx].containerID 2>/dev/null | grep $podname)

if [ $? -ne 0 ];
then
  echo "$podname[$idx] container not found"
  exit 1
fi

hname=$(echo $podinfo | awk '{print $1}')
haddr=$(echo $podinfo | awk '{print $2}')
cname=$(echo $podinfo | awk '{print $3}')
caddr=$(echo $podinfo | awk '{print $4}')
cid=$(echo $podinfo | awk '{print $5}')
cid=${cid#*docker://}
cid=${cid:0:12}

echo "====================================================================================================================================="
echo "Namespace: "$namespace "HostName: "$hname "HostAddr:"$haddr "PodName:"$cname "PodAddr:"$caddr "containerID:"$cid
echo "====================================================================================================================================="

if [ -z "$image" ];
then
  if [ $darg == "log" ];
  then
    cmd="docker logs $cid"
  elif [ $darg == "logf" ];
  then
    cmd="docker logs -f $cid"
  elif [ $darg == "inspect" ];
  then
    cmd="docker inspect $cid"
  elif [ $darg == "describe" ];
  then
    kubectl $kc $context describe pod $podname $ns
    exit 0
  elif [ $darg == "wide" ];
  then
    kubectl $kc $context get pod $podname $ns -o wide
    exit 0
  elif [ $darg == "yaml" ];
  then
    kubectl $kc $context get pod $podname $ns -o yaml
    exit 0
  elif [ $darg == "json" ];
  then
    kubectl $kc $context get pod $podname $ns -o json
    exit 0
  fi
else
  if [ $image == "xo" ];
  then
    imgsh="$image zsh $chroot"
  else
    imgsh="$image sh $chroot"
  fi

  cmd="docker run --rm -it $vol --network=container:$cid --pid=container:$cid --ipc=container:$cid $imgsh"
fi

if [ -n "$passwd" ] && type sshpass &>/dev/null;
then
  echo $passwd
  sshpass -p $passwd ssh -t -p $port -o StrictHostKeyChecking=no -o ConnectTimeout=$timeout $user@$haddr "$cmd"
else
  ssh -t -p $port -o StrictHostKeyChecking=no -o ConnectTimeout=$timeout $user@$haddr "$cmd"
fi
