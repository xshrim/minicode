#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import io
import os
import re
import sys
import time
import socks
import getopt
import socket
from urllib import request
# from urllib import parse
from pyquery import PyQuery

# sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf8')


class av:
    code = ''
    title = ''
    issuedate = ''
    length = ''
    director = ''
    category = ''
    actors = ''
    preview = ''
    links = ''

    def __init__(self, code, title, issuedate, length, director, category, actors, preview, links):
        self.code = code
        self.title = title
        self.issuedate = issuedate
        self.length = length
        self.director = director
        self.category = category
        self.actors = actors
        self.preview = preview
        self.links = links

    def display(self):
        print('番号:'.rjust(5) + self.code)
        print('标题:'.rjust(5) + self.title)
        print('发行:'.rjust(5) + self.issuedate)
        print('时长:'.rjust(5) + self.length)
        print('导演:'.rjust(5) + self.director)
        print('类别:'.rjust(5) + self.category)
        print('女优:'.rjust(5) + self.actors)
        print('预览:'.rjust(5) + self.preview)
        print('磁链:'.rjust(5) + self.links)


class avlinkinfo:
    code = ''
    title = ''
    time = ''
    hot = ''
    size = ''
    link = ''

    def __init__(self, code, title, time, hot, size, link):
        self.code = code
        self.title = title
        self.time = time
        self.hot = hot
        self.size = size
        self.link = link

    def display(self):
        print(self.code + ' -- ' + self.title + ' -- ' + self.time + ' -- ' +
              self.hot + ' -- ' + self.size + ' -- ' + self.link)


def curDir():
    try:
        return os.path.split(os.path.realpath(__file__))[0]
    except Exception as ex:
        return os.getcwd()


def getHTML(url, encode, timeout, retry, sleep, proxy=''):
    proxyDict = {}
    if proxy is not None and re.match(r'^.+@.+:.+$', proxy, flags=0):
        proxyDict['type'] = proxy.split('@')[0]
        proxy = proxy.split('@')[1]
        proxyDict['host'] = proxy.split(':')[0]
        proxyDict['port'] = proxy.split(':')[1]
    if len(proxyDict) > 0 and proxyDict['type'] is not None and proxyDict['type'].lower() == 'socks5':
        socks.set_default_proxy(socks.SOCKS5, proxyDict['host'], int(proxyDict['port']))
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
    headers = [('Host', 'img0.imgtn.bdimg.com'), ('Connection', 'keep-alive'), ('Cache-Control', 'max-age=0'),
    ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'),
    ('User-Agent', 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'), ('Accept-Encoding', '*'), ('Accept-Language', 'zh-CN,zh;q=0.8'), ('If-None-Match', '90101f995236651aa74454922de2ad74'),
    ('Referer', 'http://image.baidu.com/i?tn=baiduimage&ps=1&ct=201326592&lm=-1&cl=2&nc=1&word=%E4%BA%A4%E9%80%9A&ie=utf-8'),
    ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

    opener = request.build_opener()
    opener.addheaders = headers
    i = retry
    contents = ''
    while i > 0:
        try:
            time.sleep(sleep)
            data = opener.open(url)
            contents = data.read()
            if encode != 'no':
                contents = contents.decode(encode)
            break
        except Exception as ex:
            print(ex)
        i -= 1
    return contents


def avinfoFetch(keyword, type, engine='javbus', proxy=''):
    avs = []
    items = []
    try:
        pattern = re.compile(r'[A-Za-z]+-\d+')
        if type == 'file':
            sfile = os.path.join(curDir(), keyword)
            for line in open(sfile, encoding='utf-8'):
                items.append(line)
        else:
            items.append(keyword)
    except Exception as ex:
        print(ex)
    for item in items:
        for number in pattern.finditer(item):
            try:
                code = str(number.group()).upper()
                title = ''
                issuedate = ''
                length = ''
                director = ''
                category = ''
                actors = ''
                preview = ''
                links = ''
                if engine == 'javbus':
                    curl = 'https://www.javbus.com/' + code
                    # culr = parse.urljoin('https://www.javbus.com/', code)
                    # data = PyQuery('https://avmo.pw/cn/search/' + code)
                    # url = str(data('a.movie-box').attr('href'))
                    # avs[code] = url
                    data = PyQuery(getHTML(curl, 'utf8', 5, 5, 1, proxy))
                    content = data('div.container')
                    title = content('h3').eq(0).text()
                    avinfo = content('div[class="col-md-3 info"]')
                    issuedate = avinfo('p:eq(1)').text()
                    issuedate = str(re.search(r'\d*-\d*-\d*', issuedate).group())
                    length = avinfo('p:eq(2)').text().split(' ')[-1]
                    director = avinfo('p:eq(3)').text().split(' ')[-1]
                    category = avinfo('p:eq(7)').text()
                    actors = avinfo('p:eq(9)').text()
                    preview = content('a.bigImage').attr('href')
                    # magnets = content('table#magnet-table')
                    # links = ''
                    links = avfilter(avlinkFetch(code, 'BT工厂', proxy)).link
                avs.append(av(code, title, issuedate, length, director, category, actors, preview, links))
                print('#' * 100)
                avs[-1].display()
                # print(title + ' <-> ' + actors + ' <-> ' + preview + ' <-> ' + links)
            except Exception as ex:
                print(ex)
    return avs


def avparse(avs, tpath=curDir(), proxy=''):
    try:
        txtname = 'avinfos.txt'
        dirpath = os.path.join(curDir(), tpath)
        if not os.path.isdir(dirpath):
            os.mkdir(dirpath)
        txtpath = os.path.join(dirpath, txtname)
        txtfs = open(txtpath, 'a', encoding='utf8')
        for cav in avs:
            try:
                print('*' * 100)
                print('Creating image : ' + cav.title, end=' ...... ')
                txtfs.write('番号:'.rjust(5) + cav.code + '\n')
                txtfs.write('标题:'.rjust(5) + cav.title + '\n')
                txtfs.write('发行:'.rjust(5) + cav.issuedate + '\n')
                txtfs.write('时长:'.rjust(5) + cav.length + '\n')
                txtfs.write('导演:'.rjust(5) + cav.director + '\n')
                txtfs.write('类别:'.rjust(5) + cav.category + '\n')
                txtfs.write('女优:'.rjust(5) + cav.actors + '\n')
                txtfs.write('预览:'.rjust(5) + cav.preview + '\n')
                txtfs.write('磁链:'.rjust(5) + cav.links + '\n')
                txtfs.write('#' * 100 + '\n')
                ext = cav.preview.split('.')[-1] if '.' in cav.preview else 'jpg'
                imgname = cav.title + '.' + ext
                imgname = imgname.replace('<', '').replace('>', '').replace('/', '').replace('\\', '').replace('|', '').replace(':', '').replace('"', '').replace('*', '').replace('?', '')
                imgpath = os.path.join(dirpath, imgname)
                imgfs = open(imgpath, 'wb')
                imgfs.write(getHTML(cav.preview, 'no', 5, 5, 0, proxy))
                imgfs.close()
                print('OK')
            except Exception as ex:
                print(ex)
                print('Failed')
        txtfs.close()
    except Exception as ex:
        print(ex)


def avlinkFetch(code, engine='BT工厂', proxy=''):
    avlinks = []
    try:
        if engine == 'BT工厂':
            data = PyQuery(getHTML('http://btgongchang.org/search/' + code + '-first-asc-1', 'gbk', 5, 5, 1, proxy))
            content = data('table[class="data mb20"]')
            items = content('tr')
            # print(items)
            for item in items.items():
                try:
                    if item.attr('class') != 'firstr':
                        tds = item('td')
                        if len(tds) < 2:
                            continue
                        else:
                            title = PyQuery(tds[0])
                            time = PyQuery(tds[1])
                            hot = PyQuery(tds[2])
                            size = PyQuery(tds[3])
                            link = PyQuery(tds[4])

                            title = str(title('div.item-title')('a').text()).replace(' ', '')
                            time = str(time.text()).replace(' ', '')
                            hot = str(hot.text()).replace(' ', '')
                            size = str(size.text()).lower()
                            if 'gb' in size:
                                size = str(float(size.split(' ')[0]) * 1024)
                            elif 'kb' in size:
                                size = str(float(size.split(' ')[0]) / 1024)
                            else:
                                size = str(float(size.split(' ')[0]))
                            link = link('a:first').attr('href')
                            avlinks.append(avlinkinfo(code, title, time, hot, size, link))
                except Exception as ex:
                    print(ex)
    except Exception as ex:
        print(ex)
    return avlinks


def avfilter(avlinks):
    selectedav = None
    for avlink in avlinks:
        if selectedav is None or int(avlink.hot) > int(selectedav.hot):
            selectedav = avlink
    return selectedav


def main(argv):
    codes = []
    sfile = ''
    tpath = ''
    sengine = ''
    sproxy = ''
    try:
        opts, args = getopt.getopt(argv, "hd:e:p:f:s:", ["dir=", "engine=", "proxy=", "file=", "code="])
    except getopt.GetoptError:
        print(
            'Usage  : avfetch.py [-d <targetpath>] [-e <engine>] [-p <proxy>] [-f <filename>] [-s <codes>] [<codes>]\nExample: avfetch.py -d D:/ -e javbus -p socks5@127.0.0.1:1080 -f a.txt -s ABP-563 SRS-064 SNIS-862'
        )
        exit(2)

    if len(args) > 0:
        codes.extend(args)
    for opt, arg in opts:
        if opt == '-h':
            print(
                'Usage  : avfetch.py [-d <targetpath>] [-e <engine>] [-p <proxy>] [-f <filename>] [-s <codes>] [<codes>]\nExample: avfetch.py -d D:/ -e javbus -p socks5@127.0.0.1:1080 -f a.txt -s ABP-563 SRS-064 SNIS-862'
            )
            exit()
        elif opt in ("-d", "--dir"):
            tpath = arg
        elif opt in ("-e", "--engine"):
            sengine = arg
        elif opt in ("-p", "--proxy"):
            sproxy = arg
            if not re.match(r'^.+@.+:.+$', sproxy, flags=0):
                print('proxy format is illegal!')
                exit(2)
        elif opt in ("-f", "--file"):
            sfile = os.path.join(curDir(), arg)
        elif opt in ("-s", "--code"):
            codes.append(arg)
    try:
        if sfile != '' and os.path.isfile(sfile):
            for line in open(sfile, encoding='utf-8'):
                codes.append(line)
        if len(codes) > 0:
            avs = avinfoFetch(' '.join(codes), 'code', engine=sengine, proxy=sproxy)
            avparse(avs, tpath, sproxy)
    except Exception as ex:
        print(ex)


if __name__ == '__main__':
    main(sys.argv[1:])

# avparse('ABP-563 SRS-064 SNIS-862', 'code', 'imgs')

# main(['-d', 'imgs', '-e', 'javbus', '-s', 'EDRG-002'])
main(['-d', 'imgs', '-e', 'javbus', '-f', 'a.txt'])
#main(['-d', 'imgs', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-f', 'a.txt'])
