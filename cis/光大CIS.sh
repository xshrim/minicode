#! /bin/bash

BLUE_INFO="\033[0;34m [信息] \033[0m"
GREEN_PAAS="\033[0;32m [通过] \033[0m"
RED_FAIL="\033[0;31m [失败] \033[0m"

echo -e "$BLUE_INFO 光大CIS安全检测系统启动中。。。"
echo -e "$BLUE_INFO 1.1 API Server"
echo -e "$RED_FAIL 1.1.1 确保没有启用允许特权容器 --allow-priviledged [计分]"
echo -e "$GREEN_PAAS 1.1.2 确保没有启用匿名认证 --anonymous-auth [计分]"
echo -e "$GREEN_PAAS 1.1.3 确保没有启用允许任意token访问 --insecure-allow-any-token [计分]"
echo -e "$GREEN_PAAS 1.1.4 确保没有启用hubelet https通信 --kubelet-https [计分]"
echo -e "$RED_FAIL 1.1.5 确保没有启用不安全的绑定地址 --insecure-bind-address [计分]"
echo -e "$GREEN_PAAS 1.1.6 确保没有启用不安全的端口 --insecure-port [计分]"
echo -e "$GREEN_PAAS 1.1.7 确保没有启用了绑定安全端口 --secure-port [计分]"
echo -e "$GREEN_PAAS 1.1.8 确保审计日志保存时间不小于30 --audit-log-maxage [计分]"
echo -e "$GREEN_PAAS 1.1.9 确保审计日志备份书不小于10 --allow-log-maxbackup [计分]"