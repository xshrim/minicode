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
import chardet
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


def charDetect(data):
    charsets = ['UTF-8', 'UTF-8-SIG', 'GBK', 'GB2312', 'GB18030', 'UTF-16', 'UTF-32', 'BIG5', 'LATIN-1', 'ASCII', 'SHIFT_JIS', 'EUC-CN', 'EUC-TW', 'EUC-JP', 'EUC-KR', 'HKSCS', 'KOREAN', 'KOI8-R', 'KOI8-U', 'UTF-7', 'ISO-8859-1', 'ISO-8859-1', 'ISO-8859-5', 'ISO-8859-6', 'ISO-8859-7', 'ISO-8859-11', 'ISO-8859-15', 'TIS-620']
    try:
        charinfo = chardet.detect(data)
        data.decode(charinfo['encoding'])
        return str(charinfo['encoding']).upper()
    except Exception as ex:
        for chartype in charsets:
            try:
                data.decode(chartype)
                return chartype
            except:
                continue
    return ''
        

def getHTML(url, timeout, retry, sleep, proxy=''):
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
    ('Referer', 'http://www.deviantart.com/whats-hot/'),
    ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

    opener = request.build_opener()
    opener.addheaders = headers
    i = retry
    contents = ''

    while i > 0:
        try:
            time.sleep(sleep)
            data = opener.open(url)
            headertype = str(data.info()['Content-Type']).lower()
            contents = data.read()
            #if 'text/html' in headertype or 'text/xml' in headertype or 'text/asp' in headertype or 'text/css' in headertype or 'text/plain' in headertype or 'text/scriptlet' in headertype or 'text/h323' in headertype or 'text/asa' in headertype:
            if 'text/' in headertype:
                chartype = charDetect(contents)
                contents = contents.decode(chartype)
            break
        except Exception as ex:
            print('getHTML:' + str(ex))
            if '403:' in str(ex) or '404:' in str(ex):
                break
        i -= 1
    return contents


def avinfoFetch(keyword, type, engine='javbus', proxy=''):
    avs = []
    items = []
    codes = []
    try:
        pattern = re.compile(r'[A-Za-z]+-\d+')
        if type == 'file':
            sfile = os.path.join(curDir(), keyword)
            chartype = charDetect(open(sfile,'rb').read())
            for line in open(sfile, encoding=chartype):
                items.append(line)
        elif type == 'url':
            items.append(str(getHTML(keyword, 5, 5, 0, proxy=sproxy)))
        else:
            items.append(keyword)
    except Exception as ex:
        print('avinfoFetch:' + str(ex))
    try:
        for item in items:
            for number in pattern.finditer(item):
                code = str(number.group()).upper()
                if code not in codes:
                    codes.append(code)
    except Exception as ex:
        print('avinfoFetch:' + str(ex))
    for code in codes:
        try:
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
                data = PyQuery(getHTML(curl, 5, 5, 1, proxy))
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
                try:
                    links = avfilter(avlinkFetch(code, 'BT工厂', proxy)).link
                except Exception as ex:
                    print('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                    links = 'page:' + curl
            avs.append(av(code, title, issuedate, length, director, category, actors, preview, links))
            print('#' * 100)
            avs[-1].display()
            # print(title + ' <-> ' + actors + ' <-> ' + preview + ' <-> ' + links)
        except Exception as ex:
            print('avinfoFetch:' + str(ex))
    return avs


def avparse(avs, tpath=curDir(), proxy=''):
    if avs is not None and len(avs) > 0:
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
                    imgfs.write(getHTML(cav.preview, 5, 5, 0, proxy))
                    imgfs.close()
                    print('OK')
                except Exception as ex:
                    print('avparse:' + str(ex))
                    print('Failed')
            txtfs.close()
        except Exception as ex:
            print('avparse:' + str(ex))
    else:
        print('No AV Infomation!')


def avlinkFetch(code, engine='BT工厂', proxy=''):
    avlinks = []
    try:
        if engine == 'BT工厂':
            data = PyQuery(getHTML('http://btgongchang.org/search/' + code + '-first-asc-1', 5, 5, 1, proxy))
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
                    print('avlinkFetch:' + str(ex))
    except Exception as ex:
        print('avlinkFetch:' + str(ex))
    return avlinks


def avfilter(avlinks):
    selectedav = None
    for avlink in avlinks:
        if selectedav is None or int(avlink.hot) > int(selectedav.hot):
            selectedav = avlink
    return selectedav


def main(argv):
    codes = []
    surl = ''
    sfile = ''
    tpath = ''
    sengine = ''
    sproxy = ''
    try:
        opts, args = getopt.getopt(argv, "hd:e:p:u:f:s:", ["dir=", "engine=", "proxy=", "url=", "file=", "code="])
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
        elif opt in ("-u", "--url"):
            surl = arg
        elif opt in ("-f", "--file"):
            sfile = os.path.join(curDir(), arg)
        elif opt in ("-s", "--code"):
            codes.append(arg)
    try:
        if surl != '':
            codes.append(str(getHTML(surl, 5, 5, 0, proxy=sproxy)))
        if sfile != '' and os.path.isfile(sfile):
            chartype = charDetect(open(sfile,'rb').read())
            for line in open(sfile, encoding=chartype):
                codes.append(line)
        if len(codes) > 0:
            avs = avinfoFetch(' '.join(codes), 'code', engine=sengine, proxy=sproxy)
            avparse(avs, tpath, sproxy)
    except Exception as ex:
        print('main:' + str(ex))


if __name__ == '__main__':
    main(sys.argv[1:])

# avparse('ABP-563 SRS-064 SNIS-862', 'code', 'imgs')

main(['-d', 'imgss', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-u', 'http://btgongchang.org/'])
#main(['-d', 'imgss', '-e', 'javbus', '-f', 'a.txt'])
# main(['-d', 'imgs', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-f', 'a.txt'])
