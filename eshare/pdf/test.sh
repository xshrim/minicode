#!/bin/bash
case $# in [!4]) printf "usage : ${0##*/} stPage endPage File\n" >&2 ;; esac
stPage=$1
endPage=$2
# ((endPage++))
file=$3
mode=$4

if [[ $mode=="single" ]]; then
    gs -sstdout=/dev/null -sDEVICE=pnggray -dBATCH -dNOPAUSE -dFirstPage=${stPage} -dLastPage=${endPage} -r600 -sOutputFile=image_%d.png ${file}

else
    i=$stPage
    while ((i < endPage)); do
        gs -sstdout=/dev/null -sDEVICE=pnggray -dBATCH -dNOPAUSE -dFirstPage=${i} -dLastPage=${i} -r600 -sOutputFile=image_$i.png ${file}

        ((i++))
    done
fi
