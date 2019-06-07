#!/bin/bash
#!/bin/bash
# Usage:
# dls                                                   => 默认显示
# dls -o --format "table {{.ID}}\t{{.Command}}"         => 原始docker ps
# dls -i                                                => 完整显示(无表格格式)
# dls -f name=cli                                       => 按key过滤
# dls -s                                                => 精简显示
# dls -k                                                => 删除容器
# dls -n cli                                            => 仅对容器名包含指定关键字的容器执行操作

sep="#"
IFS=$'\n'

function append_cell(){
    #对表格追加单元格
    #append_cell col0 "col 1" ""
    #append_cell col3
    local i
    for i in "$@"
    do
        line+="|$i${sep}"
    done
}
function check_line(){
if [ -n "$line" ] 
then
    c_c=$(echo $line|tr -cd "${sep}"|wc -c)
    difference=$((${column_count}-${c_c}))
    if [ $difference -gt 0 ]
    then
        line+=$(seq -s " " $difference|sed -r s/[0-9]\+/\|${sep}/g|sed -r  s/${sep}\ /${sep}/g)
    fi
    content+="${line}|\n"
fi

}
function append_line(){
    check_line
    line=""
    local i
    for i in "$@"
    do
        line+="|$i${sep}"
    done
    check_line
    line=""
}
function segmentation(){
    local seg=""
    local i
    for i in $(seq $column_count)
    do 
        seg+="+${sep}"
    done
    seg+="${sep}+\n"
    echo $seg
}
function set_title(){
    #表格标头，以空格分割，包含空格的字符串用引号，如
    #set_title Column_0 "Column 1" "" Column3
    [ -n "$title" ] && echo "Warring:表头已经定义过,重写表头和内容"
    column_count=0
    title=""
    local i
    for i in "$@"
    do
        title+="|${i}${sep}"
        let column_count++
    done
    title+="|\n"
    seg=`segmentation`
    title="${seg}${title}${seg}"
    content=""
}
function output_table(){
    if [ ! -n "${title}" ] 
    then
        echo "未设置表头，退出" && return 1
    fi
    append_line
    table="${title}${content}$(segmentation)"
    echo -e $table|column -s "${sep}" -t|awk '{if($0 ~ /^+/){gsub(" ","-",$0);print $0}else{gsub("\\(\\*\\)","\033[31m(*)\033[0m",$0);print $0}}'

}

args=""
stop_container="false"
simple_mode="false"
full_mode="false"
container_name="all"
filter=""

if [[ $* =~ "-o" ]]                           # 使用docker ps命令的原始语法, 忽略所有dls的自定义参数语法
then
    args=$*
    args=${args//"-o"/""}
    docker ps $args
    exit 0
fi

while getopts :n:f:sik opt; do
    case $opt in
        n)
            container_name=$OPTARG
            ;;
        f)
            filter="-f "$OPTARG
            ;;
        s)
            simple_mode="true"
            ;;
        i)
            full_mode="true"
            ;;
        k)
            stop_container="true"
            ;;
        *)
            echo "Usage: dls [-o <args>] [-s] [-f <k=v>] [-i] [-k] [-n <container>]"
            exit 1
    esac
done

if [[ $stop_container == "true" ]]
then
    if [[ $container_name == "all" ]]
    then
        docker rm -f $(docker ps -aq $filter)
        exit 0
    else
        # docker rm -f $(docker ps -aqf name=$container_name $filter)
        docker rm -f $(docker ps -a $filter | grep $container_name | awk '{print $1}')
        exit 0
    fi
fi

if [[ $full_mode == "true" ]]
then
    echo "# Container: ""$container_name"
    if [[ $container_name == "all" ]]
    then
        docker inspect $(docker ps -aq $filter)
        exit 0
    else
        # docker inspect $(docker ps -aqf name=$container_name $filter)
        docker inspect $(docker ps -a $filter | grep $container_name | awk '{print $1}')
        exit 0
    fi
fi

if [[ $container_name = "all" ]]
then
    containers=$(docker inspect -f '{{.Name}}✡{{.Config.Hostname}}✡{{.Config.Image}}✡{{.Path}} {{.Args}}✡{{.HostConfig.NetworkMode}}✡{{range .NetworkSettings.Networks}}{{.IPAddress}}/{{.IPPrefixLen}}{{end}}✡{{.NetworkSettings.Ports}}-✡{{.State.Status}}' $(docker ps -aq $filter))
else
    containers=$(docker inspect -f '{{.Name}}✡{{.Config.Hostname}}✡{{.Config.Image}}✡{{.Path}} {{.Args}}✡{{.HostConfig.NetworkMode}}✡{{range .NetworkSettings.Networks}}{{.IPAddress}}/{{.IPPrefixLen}}{{end}}✡{{.NetworkSettings.Ports}}-✡{{.State.Status}}' $(docker ps -a $filter | grep $container_name | awk '{print $1}'))
fi

if [[ $simple_mode == "true" ]]
then
    set_title "容器名称" "主机名称" "启动命令" "网络地址" "端口映射" "运行状态"
else
    set_title "容器名称" "主机名称" "基础镜像" "启动命令" "网络子域" "网络地址" "端口映射" "运行状态"
fi

for container in  $containers
do
    item=$(echo $container) 
    if [[ $item != "" ]]
    then
        cname=$(echo $item | cut -d '✡' -f1)
        cname=${cname//"/"/""}
        hname=$(echo $item | cut -d '✡' -f2)
        image=$(echo $item | cut -d '✡' -f3)
        cmmnd=$(echo $item | cut -d '✡' -f4)
        cmmnd=${cmmnd//"["/""}
        cmmnd=${cmmnd//"]"/""}
        nmode=$(echo $item | cut -d '✡' -f5)
        ipadr=$(echo $item | cut -d '✡' -f6)
        ports=$(echo $item | cut -d '✡' -f7)
        ports=${ports//"0.0.0.0"/""}
        ports=${ports//"map"/""}
        ports=${ports//"["/""}
        ports=${ports//"]"/""}
        ports=${ports//"{"/""}
        ports=${ports//"}"/""}
        ports=${ports//": "/":"}
        ports=${ports//" "/","}
        if [[ $ports != "-" ]]
        then
            ports=${ports//"-"/""}
        fi
        rstat=$(echo $item | cut -d '✡' -f8)

        if [[ $simple_mode == "true" ]]
        then
             append_line $cname $hname $cmmnd $ipadr $ports $rstat
        else
             append_line $cname $hname $image $cmmnd $nmode $ipadr $ports $rstat
        fi
    fi
done
output_table
