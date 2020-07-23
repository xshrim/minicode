#!/bin/bash

function usage() {
  echo "**USAGE**:
  bench [<options>] <command>
  **OPTIONS**:
  -c run times
  -d run duration
  -h this help
  "
}
count=0
duration=0
while getopts "c:d:h" opt; do
  case $opt in
  c)
    count=$OPTARG
    ;;
  d)
    duration=$OPTARG
    ;;
  h)
    usage
    exit 0
    ;;
  \?)
    echo "Invalid option: -$OPTARG" >&2
    exit 1
    ;;
  esac
done

shift $((OPTIND - 1))

if [ $duration -eq 0 ]; then
  if [ -n $SECS ]; then
    duration=$SECS
  fi
fi

if [ $count -eq 0 ]; then
  if [ -n $COUNT ]; then
    count=$COUNT
  fi
fi

if [ $count -eq 0 ]; then
  count=1
fi

startTimestamp=$(date +%s%3N)
duration=$((1000 * $duration))
endTimestamp=$((startTimestamp + duration))

num=0

while [ 1 ]; do
  if [ $duration -ne 0 ] && [ $endTimestamp -lt $startTimestamp ] && [ $num -ge $count]; then
    break
  fi

  if [ $duration -ne 0 ] && [ $num -ge $count ]; then
    break
  fi

  eval "$*"

  sleep $((duration / 1000 / count))

  ((num++))

  startTimestamp=$(date +%s%3N)
done
