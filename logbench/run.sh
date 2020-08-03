#!/bin/bash

benchname="logbench"  # 压测deployment资源名
label="GGGGGGGG"  # 本次压测标识, 用于找出所有压测日志
num=1  # 循环次数
rate=100 # 每次循环的副本增长率
base=0   # pod副本数基数 (副本总数 = 当前次数 * 增长率 + 副本基数)
count="30000"   # 每个pod的输出日志的条数
duration=600  # 每个pod输出日志的持续时间(second)
logsize=200      # 每条日志的大致大小(byte)
esurl="http://25.4.32.9:9200"  # elasticsearch api地址
cdate=""   # 当前时间(用于找到日志索引)
estype="log"     # 索引下的类型(acp当前是按天建立索引, 每天的索引下有auth, event和log三种类型)
user="alaudaes"  # elasticsearch访问用户
passwd="cWwMJJxsxvcxjC"  # elasticsearch访问密码
context="pro"    # 当前kubectl上下文
namespace=""     # 运行压测pod的命名空间
deleteleaved=false  # 是否删除指定索引下已存在的相同label的日志

cur=0

if [ -z "$cdate" ]; then
  cdate=$(date "+%Y%m%d")
fi

if [ -n "$context" ]; then
  context="--context $context"
fi

if [ -n "$namespace" ]; then
  namespace="-n $namespace"
fi

if [ -n "$user" ]; then
  auth="-u $user:$passwd"
fi

trap "kubectl $context $namespace delete -f bench_tmp.yaml;rm -rf bench_tmp.yaml;exit 1" SIGINT SIGQUIT

while [ $(kubectl $context $namespace get pod | grep "$benchname-" | wc -l) -gt 0 ]; do
  echo "waiting pod to be deleted"
  sleep 1
done

sleep 3

echo "deleting es index log-$cdate"
curl -XDELETE -s $auth  $esurl/log-$cdate > /dev/null

sleep 3

while [ 1 ]; do
  
  repls=$((cur*rate+base))
  if [ $repls -eq 0 ]; then
    ((cur++))
    if [ $cur -gt $num ]; then
      break
    else
      continue
    fi
  fi
  
  echo "============================== round: $cur | time: $(date +%T) | replicas: $repls | rate/pod: $((count/duration))/s=============================="
  echo "elasticsearch logs: [$(curl -k -s $auth $esurl/_cat/indices|grep log-$cdate)]"
  
  lvdlogs=$(curl -k -s -H "Content-Type: application/json" -u $user:$passwd  -d "{ \"query\":{\"match\":{\"log_data\":\"[$label]\"}}}"  $esurl/log-$cdate/$estype/_search?format=json|jq -r .hits.total)
  echo "leaved logs : $lvdlogs"
  if [ "$deleteleaved" = true ]; then
    lvdlogs_total=$lvdlogs
    echo -ne "deleting leaved log documents [$lvdlogs_total/$lvdlogs_total] ...\r"

    while [ $lvdlogs -gt 0 ]; do
      for lvdlogid in $(curl -k -s $auth "$esurl/log-$cdate/$estype/_search?q=$label&size=10000" | jq -r .hits.hits[]._id);
      do
        curl -XDELETE -k -s $auth $esurl/log-$cdate/$estype/$lvdlogid > /dev/null
      done
      lvdlogs=$(curl -k -s -H "Content-Type: application/json" $auth -d "{ \"query\":{\"match\":{\"log_data\":\"[$label]\"}}}"  $esurl/log-$cdate/$estype/_search?format=json|jq -r .hits.total)
      echo -ne "deleting leaved log documents [$lvdlogs/$lvdlogs_total] ...\r"
    done
    echo -n "deleting leaved log documents [0/$lvdlogs_total] ..."
    echo " done"
  fi

  # 配置压测资源清单
  cp bench.yaml bench_tmp.yaml

  if [ ! -f bench_tmp.yaml ]; then
    echo "bench_tmp.yaml file not found"
    exit 1
  fi
  
  sed -i "s/{{name}}/$benchname/g" bench_tmp.yaml # 设置压测deployment资源名
  sed -i "s/^  replicas: .*/  replicas: $repls/" bench_tmp.yaml # 更新副本数量
  sed -i "s/{{label}}/$label/" bench_tmp.yaml # 设置日志label
  sed -i "s/{{count}}/$count/" bench_tmp.yaml # 设置日志量
  sed -i "s/{{duration}}/${duration}s/" bench_tmp.yaml # 设置日志输出持续时间
  sed -i "s/{{size}}/$logsize/" bench_tmp.yaml # 设置日志大小
  
  kubectl $context $namespace apply -f bench_tmp.yaml
  
  # 等待所有pod启动
  while [ 1 ]; do
    #want=$(kubectl $context $namespae get deployment $benchname -o json|jq .status.replicas)
    #ready=$(kubectl $context $namespae get deployment $benchname -o json|jq .status.readyReplicas)
    want=$repls
    ready=$(kubectl --context pro get pod | grep "$benchname-" | grep "Running" | wc -l)
    if [ $ready ==  "null" ]; then
      ready=0
    fi
    echo -ne "---------- deployment $benchname replicas: $want want / $ready ready ----------\r"
    if [ $ready -eq $want ]; then
      break
    else
      sleep 1
    fi
  done
  echo -ne "\n"

  # 等待所有pod执行完成
  completed=0
  pods=$(kubectl --context pro get pod | grep "$benchname-" | grep "Running" | awk '{print $1}')
  podtrim=$(echo $pods|xargs)
  while [ ${#podtrim} -gt 0 ]; do
    for pod in $pods; do
      lastlog=$(kubectl --context pro logs $pod --tail 1 | awk '{print $NF}')
      if [[ "$lastlog" =~ "$((count-1))" ]]; then
        delete=$pod
        pods=("${pods[@]/$delete}")
        ((completed++))
        echo -ne "---------- deployment $benchname replicas: $ready running / $completed completed ----------\r"
      fi
    done
    podtrim=$(echo $pods|xargs)
    sleep 1
  done
  echo -ne "\n"
  
  echo "********** [logs output] done ********** time: $(date +%T) | bench logs: $((repls*count)) **********"

  logs=$lvdlogs
  mglogs=0
  lastlogcount=$((logs-lvdlogs-mglogs-1))
  while [ $((logs-lvdlogs-mglogs)) -lt $((repls*count+1)) ]; do
    logs=$(curl -k -s -H "Content-Type: application/json" $auth -d "{ \"query\":{\"match\":{\"log_data\":\"$label\"}}}"  $esurl/log-$cdate/$estype/_search?format=json|jq -r .hits.total)
    mglogs=$(curl -s -H "Content-Type: application/json" $auth  -d "{\"query\": {\"bool\": {\"must\": [{\"match\": {\"log_data\":\"$label\"}}, {\"match\": {\"pod_name\":\"morgans\"}}]}}}" $esurl/log-$(date +%Y%m%d)/log/_search|jq -r .hits.total)
    echo -ne "---------- bench logs in elasticsearch: $((repls*count)) desired / $((logs-lvdlogs-mglogs-1)) written ----------\r"
    if [ $lastlogcount -eq $((logs-lvdlogs-mglogs-1)) ]; then
      break
    else
      lastlogcount=$((logs-lvdlogs-mglogs-1))
    fi
    sleep 3
  done
  echo -ne "\n"
  echo "********** [write to es] done ********** time: $(date +%T) | bench logs: $((logs-lvdlogs-mglogs-1)) **********"

  echo "elasticsearch logs: [$(curl -k -s $auth $esurl/_cat/indices|grep log-$cdate)]"

  kubectl $context $namespace delete -f bench_tmp.yaml
  rm -rf bench_tmp.yaml

  # 删除压测日志
  if [ "$deleteleaved" = true ]; then
    logs_total=$logs
    echo -ne "deleting bench log documents [$logs_total/$logs_total] ...\r"
    while [ $logs -gt 0 ]; do
      for logid in $(curl -k -s $auth "$esurl/log-$cdate/$estype/_search?q=$label&size=10000" | jq -r .hits.hits[]._id);
      do
        curl -XDELETE -k -s $auth $esurl/log-$cdate/$estype/$logid > /dev/null
      done
      logs=$(curl -k -s -H "Content-Type: application/json" $auth -d "{ \"query\":{\"match\":{\"log_data\":\"[$label]\"}}}"  $esurl/log-$cdate/$estype/_search?format=json|jq -r .hits.total)
      echo -ne "deleting bench log documents [$logs/$logs_total] ...\r"
    done
    echo -n "deleting bench log documents [0/$logs_total] ..."
    echo " done"
  fi

  echo "### round $cur completed [$(date +%T)]"

  ((cur++))
  if [ $cur -gt $num ]; then
    break
  fi

  echo "waiting 60s for next round"
  sleep 60  # 每轮暂停1分钟, 等待Deployment中的pod执行完, 监控曲线能够明显区分每轮测试
  
done

echo ">>> log bench completed [$(date +%T)]"
