# coding:utf-8
# 爬取斗鱼弹幕


import json
import random
import re
import socket
import sys
import threading
import time
import urllib.parse
import urllib.request


class douYuTVDanmu(object):
    def __init__(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.codeLocalToServer = 689
        self.serverToLocal = 690
        self.gid = -9999
        self.rid = 16789
        self.server = {}

    def log(self, str):
        # str = str.encode()
        now_time = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime())
        log = now_time + '\t\t' + str
        with open('log.txt', 'a', encoding='utf-8')as f:
            f.writelines(log + '\n')
        print(log)

    def sendMsg(self, msg):
        msg = msg.encode('utf-8')
        data_length = len(msg) + 8
        msgHead = int.to_bytes(data_length, 4, 'little') + int.to_bytes(data_length, 4, 'little') + int.to_bytes(self.codeLocalToServer, 4, 'little')
        self.sock.send(msgHead)
        self.sock.sendall(msg)

    def keeplive(self):
        while True:
            msg = 'type@=keeplive/tick@=' + str(int(time.time())) + '/\x00'
            self.sendMsg(msg)
            # keeplive=sock.recv(1024)
            time.sleep(20)

    def getHTML(self, url, timeout=5, retry=3, sleep=0, proxy=''):
        proxyDict = {}
        if proxy is not None and re.match(r'^.+@.+:.+$', proxy, flags=0):
            proxyDict['type'] = proxy.split('@')[0]
            proxy = proxy.split('@')[1]
            proxyDict['host'] = proxy.split(':')[0]
            proxyDict['port'] = proxy.split(':')[1]
        if len(proxyDict) > 0 and proxyDict['type'] is not None and proxyDict['type'].lower() == 'socks5':
            socks.set_default_proxy(socks.SOCKS5, proxyDict['host'], int(proxyDict['port']))
            socket.socket = socks.socksocket
        elif len(proxyDict) > 0 and proxyDict['type'] is not None and proxyDict['type'].lower() == 'socks4':
            socks.set_default_proxy(socks.SOCKS4, proxyDict['host'], int(proxyDict['port']))
            socket.socket = socks.socksocket
        elif len(proxyDict) > 0 and proxyDict['type'] is not None and proxyDict['type'].lower() == 'http':
            socks.set_default_proxy(socks.HTTP, proxyDict['host'], int(proxyDict['port']))
            socket.socket = socks.socksocket
        socket.setdefaulttimeout(timeout)
        # url = 'https://www.javbus2.com/HIZ-015'
        # url = "http://img0.imgtn.bdimg.com/it/u=4054848240,1657436512&fm=21&gp=0.jpg"
        # headers = [('User-Agent','Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.11 (KHTML, like Gecko) \
        # Chrome/23.0.1271.64 Safari/537.11'),
        # ('Accept','text/html;q=0.9,*/*;q=0.8'),
        # ('Accept-Charset','ISO-8859-1,utf-8;q=0.7,*;q=0.3'),
        # ('Accept-Encoding','gzip,deflate,sdch'),
        # ('Connection','close'),
        # ('Referer',None )]#注意如果依然不能抓取的话，这里可以设置抓取网站的host
        headers = [('Host', 'img0.imgtn.bdimg.com'), ('Connection', 'close'), ('Cache-Control', 'max-age=0'), ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'), ('User-Agent', 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'), ('Accept-Encoding', '*'), ('Accept-Language', 'zh-CN,zh;q=0.8'), ('If-None-Match', '90101f995236651aa74454922de2ad74'), ('Referer', 'http://www.deviantart.com/whats-hot/'), ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

        opener = urllib.request.build_opener()
        opener.addheaders = headers
        i = retry
        contents = ''
        while i > 0:
            try:
                time.sleep(sleep)
                data = opener.open(urllib.parse.quote(url, safe='/:?=%-&'))
                headerinfo = data.info()
                headertype = str(headerinfo['Content-Type']).lower()
                contents = data.read()

                if 'text/' in headertype:
                    if str(headerinfo['Content-Encoding']).lower() == 'gzip':
                        contents = zlib.decompress(contents, 16 + zlib.MAX_WBITS)
                    if 'charset' in headertype:
                        for item in ['utf-8', 'utf8', 'gbk', 'gb2312', 'gb18030', 'big5', 'latin-1', 'latin1']:
                            if item in headertype:
                                chartype = item.upper()
                    else:
                        chartype = charDetect(contents)
                    contents = contents.decode(chartype, errors='ignore')
                opener.close()
                break
            except Exception as ex:
                print(str(ex))
                opener.close()
                if '403' in str(ex) or '404' in str(ex) or '11001'in str(ex):
                    break
            i -= 1
        return contents

    # 分析网页中的信息
    def getInfo(self, url):
        self.log("请求网页内容...")
        try:
            data = self.getHTML(url)
            # with urllib.request.urlopen(url)as f:
                # data = f.read().decode()
        except BaseException as e:
            self.log("请求网页内容失败...")
            exit(404)
        self.log("获取房间信息...")
        room = re.search('var \$ROOM = (.*);', data)
        if room:
            room = room.group(1)
            room = json.loads(room)
            self.log("房间名:" + room["room_name"] + '\t\t|\t\t主播:' + room["owner_name"])
            self.rid = room["room_id"]

            if room["show_status"] == 2:
                self.log("未开播!\t\t" + str(self.rid))
                exit(1)
        else:
            print('不存在')
        # self.log("获取弹幕服务器信息...")
        # server_config = re.search('server_config":"(.*)","de',data)
        # if server_config:
        #     server_config = server_config.group(1)
        #     server_config = urllib.parse.unquote(server_config)
        #     server_config = json.loads(server_config)
        #     self.server = server_config

    def connectToDanMuServer(self):
        # index = random.randint(0,len(self.server)-1)     #选择一个服务器
        # HOST = self.server[index]['ip']
        # PORT = self.server[index]['port']
        HOST = 'openbarrage.douyutv.com'
        PORT = 8601

        self.log("连接弹幕服务器..." + HOST + ':' + str(PORT))
        self.sock.connect((HOST, PORT))
        self.log("连接成功,发送登录请求...")

        # 抓包:msg = 'type@=loginreq/username@=qq_NEUBLK/password@=1234567890123456/roomid@=335166/'
        # msg = 'type@=loginreq/username@=/password@=/roomid@=' + str(self.rid) + '/\x00'
        msg = 'type@=loginreq/roomid@=' + str(self.rid) + '/\x00'
        self.sendMsg(msg)
        data = self.sock.recv(1024)
        self.log('Received from login\t\t' + repr(data))
        a = re.search(b'type@=(\w*)', data)
        if a.group(1) != b'loginres':
            self.log("登录失败,程序退出...")
            exit(0)
        self.log("登录成功")

        # data = self.sock.recv(1024)   #type@=msgrepeaterlist  服务器列表
        # self.log('msgrepeaterlist\t\t'+ repr(data))
        # gid = re.search('gid@=(.*)\/',data)
        # # data = self.sock.recv(1024)
        # # self.log('Received', repr(data))
        # if gid:
        #     self.gid = gid.group(1)
        #     self.log("找到弹幕服务器"+str(self.gid))

        msg = 'type@=joingroup/rid@=' + str(self.rid) + '/gid@=-9999/\x00'
        # print(msg)
        self.sendMsg(msg)
        self.log("进入弹幕服务器...")
        threading.Thread(target=douYuTVDanmu.keeplive, args=(self,)).start()
        self.log("心跳包机制启动...")
        data = self.sock.recv(1024)
        # print('Received', repr(data))

    def danmuWhile(self):
        self.log("监听中")
        while True:
            data = self.sock.recv(9999)
            print(data)
            # self.log(repr(data))
            a = re.search(b'type@=(\w*)', data)
            if a:
                if a.group(1) == b'chatmsg':
                    danmu = re.search(b'nn@=(.*)/txt@=(.*?)/', data)
                    # self.log(danmu.group(1).decode()+'\t:\t'+danmu.group(2).decode())
                    try:
                        self.log(danmu.group(2).decode())
                    except BaseException as e:
                        self.log("\t\t_________解析弹幕信息失败:" + str(data))


def main(url):
    danmu = douYuTVDanmu()
    danmu.getInfo(url)
    danmu.connectToDanMuServer()
    danmu.danmuWhile()


if __name__ == '__main__':
    if len(sys.argv) > 1:
        main('https://www.douyu.com/' + str(sys.argv[1]))

main('https://www.douyu.com/99039')
