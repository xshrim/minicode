#!/bin/bash

startTimestamp=$(date +%s%3N)
duration=$((1000*$SECS))
endTimestamp=$((startTimestamp+duration))
count=0

echo $startTimestamp
echo $endTimestamp

while [ 1 ]
do
  if [ $endTimestamp -lt $startTimestamp ]; then
    break
  fi
  
  echo "Kubernetes is an open-source platform for managing containerized services. This is a benchmark container log: $count"
  
  ((count++))
  
  startTimestamp=$(date +%s%3N)
done
