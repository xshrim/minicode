#!/bin/bash

count=1

echo "running"

while [ 1 ]; do
  if [ $count -gt 10 ]; then
    break
  fi
  
  repls=$((count*100))
  
  echo "============================== time: $(date +%T) | replicas: $repls =============================="
  echo "elasticsearch logs: "
  echo "$(curl -s -u alaudaes:cWwMJJxsxvcxjC  http://25.4.32.9:9200/_cat/indices|grep log)"
  
  sed -i "s/^  replicas: .*/  replicas: $repls/" bench.yaml # 更新副本数量
  
  kubectl --context pro apply -f bench.yaml
  
  sleep 600
  echo "********** done **********"
  echo "elasticsearch logs: "
  echo "$(curl -s -u alaudaes:cWwMJJxsxvcxjC  http://25.4.32.9:9200/_cat/indices|grep log)"
  
  kubectl --context pro delete -f bench.yaml
  
  echo "wating 300s"
  sleep 300  # 每轮暂停5分钟, 等待Deployment中的pod执行完, 监控曲线能够明显区分每轮测试
  
  ((count++))
done

echo "completed"
