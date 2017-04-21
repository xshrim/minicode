#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import re
import sys
import time
import getopt
import socket
from urllib import request
from pyquery import PyQuery as pq


class av:
    code = ''
    title = ''
    preview = ''
    
    def __init__(self, code, title, preview):
        self.code = code
        self.title = title
        self.preview = preview
        

def getHTML(url, encode, timeout, retry, sleep):
    socket.setdefaulttimeout(timeout)
    #url = 'https://www.javbus2.com/HIZ-015'
    #url = "http://img0.imgtn.bdimg.com/it/u=4054848240,1657436512&fm=21&gp=0.jpg"
    # headers = [('User-Agent','Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11'),
    # ('Accept','text/html;q=0.9,*/*;q=0.8'),
    # ('Accept-Charset','ISO-8859-1,utf-8;q=0.7,*;q=0.3'),
    # ('Accept-Encoding','gzip,deflate,sdch'),
    # ('Connection','close'),
    # ('Referer',None )]#注意如果依然不能抓取的话，这里可以设置抓取网站的host
    headers = [
    ('Host','img0.imgtn.bdimg.com'),
    ('Connection', 'keep-alive'),
    ('Cache-Control', 'max-age=0'),
    ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'),
    ('User-Agent', 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'),
    ('Accept-Encoding','*'),
    ('Accept-Language', 'zh-CN,zh;q=0.8'),
    ('If-None-Match', '90101f995236651aa74454922de2ad74'),
    ('Referer','http://image.baidu.com/i?tn=baiduimage&ps=1&ct=201326592&lm=-1&cl=2&nc=1&word=%E4%BA%A4%E9%80%9A&ie=utf-8'),
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


def avbus(keyword, type):
    avs = []
    items = []
    try:
        pattern = re.compile(r'[A-Za-z]+-\d+')
        if type == 'file':
            sfile = os.path.join(os.getcwd(), keyword)
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
                url = 'https://www.javbus2.com/' + code
                #data = pq('https://avmo.pw/cn/search/' + code)
                #url = str(data('a.movie-box').attr('href'))
                #avs[code] = url
                data = pq(getHTML(url, 'utf8', 5, 5, 1))
                content = data('div.container')
                title = content('h3').eq(0).text()
                imgurl = content('a.bigImage').attr('href')
                avs.append(av(code, title, imgurl))
                print(title + ' <-> ' + imgurl)
            except Exception as ex:
                print(ex)
    return avs


def avparse(keyword, type, tpath=os.getcwd()):
    for item in avbus(keyword, type):
        try:
            print('Creating image : ' + item.title, end=' ...... ')
            ext = item.preview.split('.')[-1] if '.' in item.preview else 'jpg'
            filename = item.title + '.' + ext
            filename = filename.replace('<', '').replace('>','').replace('/','').replace('\\', '').replace('|', '').replace(':', '').replace('"', '').replace('*', '').replace('?', '')
            dirpath = os.path.join(os.getcwd(), tpath)
            if not os.path.isdir(dirpath):
                os.mkdir(dirpath)
            filepath = os.path.join(dirpath, filename)
            fs = open(filepath, 'wb')  
            fs.write(getHTML(item.preview, 'no', 5, 5, 0))  
            fs.close() 
            print('OK')
        except Exception as ex:
            print(ex)
            print('Failed')

def main(argv):
    codes = []
    sfile = ''
    tpath = ''
    try:
        opts, args = getopt.getopt(argv, "hd:f:s:", ["dir=", "file=", "code="])
    except getopt.GetoptError:
        print ('Usage  : avfetch.py [-d <targetpath>] [-f <filename>] [-s <codes>] [<codes>]\nExample: avfetch.py -d D:/ -f a.txt -s ABP-563 SRS-064 SNIS-862')
        exit(2)
    
    if len(args) > 0:
        codes.extend(args)
    for opt, arg in opts:
        if opt == '-h':
            print ('Usage  : avfetch.py [-d <targetpath>] [-f <filename>] [-s <codes>] [<codes>]\nExample: avfetch.py -d D:/ -f a.txt -s ABP-563 SRS-064 SNIS-862')
            exit()
        elif opt in ("-d", "--dir"):
            tpath = arg
        elif opt in ("-f", "--file"):
            sfile = os.path.join(os.getcwd(), arg)
        elif opt in ("-s", "--code"):
            codes.append(arg)
    try:
        if sfile != '' and os.path.isfile(sfile):
            for line in open(sfile, encoding='utf-8'):
                codes.append(line)
        if len(codes) > 0:
            avparse(' '.join(codes), 'code', tpath)
    except Exception as ex:
        print(ex)
            
    
if __name__ == '__main__':
    main(sys.argv[1:])

#avparse('ABP-563 SRS-064 SNIS-862', 'code', 'imgs')    
    
main(['-d', 'imgs', '-f', 'a.txt'])   

