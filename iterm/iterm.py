#!/usr/bin/env python
# -*- coding:utf-8 -*-
import paramiko
import sys
import os
import socket
import getpass
import termios
import tty
import select
from paramiko.py3compat import u


def interactive_shell(chan):
    # 获取原tty属性
    oldtty = termios.tcgetattr(sys.stdin)
    try:
        # 为tty设置新属性
        # 默认当前tty设备属性：
        # 输入一行回车，执行
        # CTRL+C 进程退出，遇到特殊字符，特殊处理。

        # 这是为原始模式，不认识所有特殊符号
        # 放置特殊字符应用在当前终端，如此设置，将所有的用户输入均发送到远程服务器
        tty.setraw(sys.stdin.fileno())
        tty.setcbreak(sys.stdin.fileno())
        chan.settimeout(0.0)
        while True:
            r, w, e = select.select([chan, sys.stdin], [], [])
            if chan in r:
                try:
                    x = u(chan.recv(1024))
                    if len(x) == 0:
                        # sys.stdout.write('\r\n*** EOF\r\n')
                        break
                    sys.stdout.write(x)
                    sys.stdout.flush()
                except socket.timeout:
                    pass
            if sys.stdin in r:
                x = sys.stdin.read(1)
                if len(x) == 0:
                    break
                chan.send(x)

    finally:
        # 重新设置终端属性
        termios.tcsetattr(sys.stdin, termios.TCSADRAIN, oldtty)


def run():
    db_dict = {
        '127.0.0.1': {
            'root': {
                'user': 'root',
                'auth': 'p',
                "cert": 'root'
            },
        },
        '172.16.201.189': {
            # 'root': {'user': 'root', 'auth': 'r', "cert": 'key路径'},
            'root': {
                'user': 'root',
                'auth': 'p',
                "cert": 'szyz#233'
            },
        },
        '172.16.201.191': {
            'root': {
                'user': 'root',
                'auth': 'p',
                "cert": 'szyz#233'
            },
        },
        '172.16.201.192': {
            'root': {
                'user': 'root',
                'auth': 'p',
                "cert": 'szyz#233'
            },
        },
        '172.16.201.193': {
            'root': {
                'user': 'root',
                'auth': 'p',
                "cert": 'szyz#233'
            },
        },
    }

    for row in db_dict.keys():
        print(row)

    hostname = input('请选择主机: ')
    tran = paramiko.Transport((
        hostname,
        22,
    ))
    tran.start_client()

    for item in db_dict[hostname].keys():
        print(item)

    username = input('请输入用户: ')

    user_dict = db_dict[hostname][username]
    if user_dict['auth'] == 'r':
        key = paramiko.RSAKey.from_private_key_file(user_dict['cert'])
        tran.auth_publickey(username, key)
    else:
        pw = user_dict['cert']
        tran.auth_password(username, pw)

    # 打开一个通道
    chan = tran.open_session()
    # 获取一个终端
    chan.get_pty()
    # 激活器
    chan.invoke_shell()

    interactive_shell(chan)

    chan.close()
    tran.close()


if __name__ == '__main__':
    run()
