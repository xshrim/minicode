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
import sqlite3
import logging
import threading
from urllib import request
from urllib import parse
from urllib.parse import quote
from pyquery import PyQuery

# sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf8')


class av(object):
    '''
    code = ''
    title = ''
    issuedate = ''
    length = ''
    mosaic = ''
    director = ''
    manufacturer = ''
    publisher = ''
    series = ''
    category = ''
    actors = ''
    favor = ''
    coverlink = ''
    cover = b''
    link = ''
    '''

    def __init__(self, code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link):
        self.code = code
        self.title = title
        self.issuedate = issuedate
        self.length = length
        self.mosaic = mosaic
        self.director = director
        self.manufacturer = manufacturer
        self.publisher = publisher
        self.series = series
        self.category = category
        self.actors = actors
        self.favor = favor
        self.coverlink = coverlink
        if isinstance(cover, bytes):
            self.cover = cover
        else:
            self.cover = str(cover).encode()
        self.link = link

    def __str__(self):
        '''
        print('番号:'.rjust(5) + self.code)
        print('标题:'.rjust(5) + self.title)
        print('日期:'.rjust(5) + self.issuedate)
        print('时长:'.rjust(5) + self.length)
        print('修正:'.rjust(5) + self.mosaic)
        print('导演:'.rjust(5) + self.director)
        print('制作:'.rjust(5) + self.manufacturer)
        print('发行:'.rjust(5) + self.publisher)
        print('系列:'.rjust(5) + self.series)
        print('类别:'.rjust(5) + self.category)
        print('女优:'.rjust(5) + self.actors)
        print('收藏:'.rjust(5) + self.favor)
        print('预览:'.rjust(5) + self.coverlink)
        print('磁链:'.rjust(5) + self.link)
        '''
        return '番号:'.rjust(5) + self.code + '\n' + '标题:'.rjust(5) + self.title + '\n' + '日期:'.rjust(5) + self.issuedate + '\n' + '时长:'.rjust(5) + self.length + '\n' + '修正:'.rjust(5) + self.mosaic + '\n' + '导演:'.rjust(5) + self.director + '\n' + '制作:'.rjust(5) + self.manufacturer + '\n' + '发行:'.rjust(5) + self.publisher + '\n' + '系列:'.rjust(5) + self.series + '\n' + '类别:'.rjust(5) + self.category + '\n' + '女优:'.rjust(5) + self.actors + '\n' + '收藏:'.rjust(5) + self.favor + '\n' + '预览:'.rjust(5) + self.coverlink + '\n' + '磁链:'.rjust(5) + self.link

    __repr__ = __str__


class avlink(object):
    '''
    code = ''
    title = ''
    time = ''
    hot = ''
    size = ''
    link = ''
    origin = ''
    '''

    def __init__(self, code, title, time, hot, size, link, origin):
        self.code = code
        self.title = title
        self.time = time
        self.hot = hot
        self.size = size
        self.link = link
        self.origin = origin

    def __str__(self):
        return self.code + ' -- ' + self.title + ' -- ' + self.time + ' -- ' + self.hot + ' -- ' + self.size + ' -- ' + self.link + ' -- ' + self.origin

    __repr__ = __str__

######################################### DB START#########################################


def dict_factory(cursor, row):
    '''将数据库查询结果按字典输出的字典工厂'''
    d = {}
    for idx, col in enumerate(cursor.description):
        d[col[0]] = row[idx]
    return d


def get_conn(path):
    '''获取到数据库的连接对象，参数为数据库文件的绝对路径
    如果传递的参数是存在，并且是文件，那么就返回硬盘上面改
    路径下的数据库文件的连接对象；否则，返回内存中的数据接
    连接对象'''
    conn = sqlite3.connect(path)
    if os.path.exists(path) and os.path.isfile(path):
        # print('硬盘上面:[{}]'.format(path))
        return conn
    else:
        conn = None
        # print('内存上面:[:memory:]')
        return sqlite3.connect(':memory:')


def get_cursor(conn):
    '''该方法是获取数据库的游标对象，参数为数据库的连接对象
    如果数据库的连接对象不为None，则返回数据库连接对象所创
    建的游标对象；否则返回一个游标对象，该对象是内存中数据
    库连接对象所创建的游标对象'''
    if conn is not None:
        return conn.cursor()
    else:
        return get_conn('').cursor()


def drop_table(conn, table):
    '''如果表存在,则删除表，如果表中存在数据的时候，使用该
    方法的时候要慎用！'''
    if table is not None and table != '':
        sql = 'DROP TABLE IF EXISTS ' + table
        # print('执行sql:[{}]'.format(sql))
        cu = get_cursor(conn)
        cu.execute(sql)
        conn.commit()
        # print('删除数据库表[{}]成功!'.format(table))
        close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def create_table(conn, sql):
    '''创建数据库表'''
    if sql is not None and sql != '':
        cu = get_cursor(conn)
        # print('执行sql:[{}]'.format(sql))
        cu.execute(sql)
        conn.commit()
        # print('创建数据库表成功!'
        close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def close_all(conn, cu):
    '''关闭数据库游标对象和数据库连接对象'''
    try:
        if cu is not None:
            cu.close()
    finally:
        if cu is not None:
            cu.close()


def save(conn, sql, data):
    '''插入数据'''
    if sql is not None and sql != '':
        if data is not None:
            cu = get_cursor(conn)
            for d in data:
                # print('执行sql:[{}],参数:[{}]'.format(sql, d))
                cu.execute(sql, d)
                conn.commit()
            close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def fetchall(conn, sql):
    '''查询所有数据'''
    if sql is not None and sql != '':
        cu = get_cursor(conn)
        # print('执行sql:[{}]'.format(sql))
        cu.execute(sql)
        r = cu.fetchall()
        if len(r) > 0:
            for e in range(len(r)):
                print(r[e])
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def fetchone(conn, sql, data):
    '''查询一条数据'''
    if sql is not None and sql != '':
        if data is not None:
            # Do this instead
            d = (data,)
            cu = get_cursor(conn)
            # print('执行sql:[{}],参数:[{}]'.format(sql, data))
            cu.execute(sql, d)
            r = cu.fetchall()
            if len(r) > 0:
                for e in range(len(r)):
                    print(r[e])
        else:
            logging.error('the [{}] equal None!'.format(data))
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def update(conn, sql, data):
    '''更新数据'''
    if sql is not None and sql != '':
        if data is not None:
            cu = get_cursor(conn)
            for d in data:
                # print('执行sql:[{}],参数:[{}]'.format(sql, d))
                cu.execute(sql, d)
                conn.commit()
            close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def delete(conn, sql, data):
    '''删除数据'''
    if sql is not None and sql != '':
        if data is not None:
            cu = get_cursor(conn)
            for d in data:
                # print('执行sql:[{}],参数:[{}]'.format(sql, d))
                cu.execute(sql, d)
                conn.commit()
            close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))

######################################### DB END#########################################


def logInit():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s', datefmt='%a, %d %b %Y %H:%M:%S', filename=os.path.join(curDir(), 'avfetch.log'), filemode='w')

    console = logging.StreamHandler()
    console.setLevel(logging.ERROR)
    formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
    console.setFormatter(formatter)
    logging.getLogger('').addHandler(console)


def curDir():
    try:
        return os.path.split(os.path.realpath(__file__))[0]
    except Exception as ex:
        return os.getcwd()


def charDetect(data):
    charsets = ['UTF-8', 'UTF-8-SIG', 'GBK', 'GB2312', 'GB18030', 'BIG5', 'SHIFT_JIS', 'EUC-CN', 'EUC-TW', 'EUC-JP', 'EUC-KR', 'ASCII', 'HKSCS', 'KOREAN', 'UTF-7', 'TIS-620', 'LATIN-1', 'KOI8-R', 'KOI8-U', 'ISO-8859-5', 'ISO-8859-6', 'ISO-8859-7', 'ISO-8859-11', 'ISO-8859-15', 'UTF-16', 'UTF-32']
    try:
        charinfo = chardet.detect(data)
        data.decode(charinfo['encoding'])
        return str(charinfo['encoding']).upper()
    except Exception as ex:
        logging.debug('charDetect:' + str(ex))
        for chartype in charsets:
            try:
                data.decode(chartype)
                return chartype
            except Exception as ex:
                logging.debug('charDetect:' + str(ex))
                continue
    return ''


def detectPage(url, timeout, retry, sleep, proxy=''):
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

    opener = request.build_opener()
    opener.addheaders = headers
    i = retry
    while i > 0:
        try:
            time.sleep(sleep)
            data = opener.open(quote(url, safe='/:?=%-&'))
            opener.close()
            return True
        except Exception as ex:
            opener.close()
            if '403' in str(ex) or '404' in str(ex) or '11001'in str(ex):
                return False
        i -= 1
    return False


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

    opener = request.build_opener()
    opener.addheaders = headers
    i = retry
    contents = ''
    while i > 0:
        try:
            time.sleep(sleep)
            data = opener.open(quote(url, safe='/:?=%-&'))
            headertype = str(data.info()['Content-Type']).lower()
            contents = data.read()
            if 'text/' in headertype:
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
            opener.close()
            logging.debug('getHTML:' + str(ex))
            if '403' in str(ex) or '404' in str(ex) or '11001'in str(ex):
                break
        i -= 1
    return contents


def avpageFetch(url, engine='javbus', proxy=''):
    avpage = []
    if engine == 'javbus':
        try:
            curl = url
            urldata = PyQuery(getHTML(curl, 5, 5, 1, proxy))
            urlcontent = urldata('div#waterfall')
            for urlitem in urlcontent('a.movie-box').items():
                code = str(urlitem('div.photo-info')('date:first').text()).strip().upper()
                urlitem('div.photo-info')('span')('date').remove()
                title = str(urlitem('div.photo-info')('span').text().replace('/', '')).strip()
                url = str(urlitem.attr('href')).strip()
                avpage.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javbus:" + str(ex))
    return avpage


def avurlFetch(keyword, engine='javbus', proxy=''):
    avpages = []
    if engine == 'javbus':
        surls = []
        try:
            urldata = PyQuery(getHTML('https://www.javbus.com/search/' + keyword, 5, 5, 1, proxy))
            pagination = urldata('ul[class="pagination pagination-lg"]')
            if pagination is not None and str(pagination).strip() != '':
                for subpage in pagination('li').items():
                    if re.match(r'^.*\d+.*$', str(subpage.text())):
                        if subpage('a').attr('href') is not None and str(subpage('a').attr('href')).strip() != '':
                            surls.append(parse.urljoin('https://www.javbus.com/', subpage('a').attr('href')))
            else:
                surls.append('https://www.javbus.com/search/' + keyword)
            for surl in surls:
                urldata = PyQuery(getHTML(surl, 5, 5, 1, proxy))
                urlcontent = urldata('div#waterfall')
                for urlitem in urlcontent('a.movie-box').items():
                    code = str(urlitem('div.photo-info')('date:first').text()).strip().upper()
                    urlitem('div.photo-info')('span')('date').remove()
                    title = str(urlitem('div.photo-info')('span').text().replace('/', '')).strip()
                    url = str(urlitem.attr('href')).strip()
                    avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javbus:" + str(ex))
    if engine == 'javhoo':
        surls = []
        try:
            urldata = PyQuery(getHTML('https://www.javhoo.com/search/' + keyword, 5, 5, 1, proxy))
            pagination = urldata('ul[class="pagination pagination-lg"]')
            if pagination is not None and str(pagination).strip() != '':
                for subpage in pagination('li').items():
                    if re.match(r'^.*\d+.*$', str(subpage.text())):
                        if subpage('a').attr('href') is not None and str(subpage('a').attr('href')).strip() != '':
                            surls.append(parse.urljoin('https://www.javhoo.com/', subpage('a').attr('href')))
            else:
                surls.append('https://www.javhoo.com/search/' + keyword)
            for surl in surls:
                urldata = PyQuery(getHTML(surl, 5, 5, 1, proxy))
                urlcontent = urldata('div#content')
                for urlitem in urlcontent('div[class="wf-cell iso-item"]').items():
                    code = str(urlitem('div.project-list-content')('date:first').text()).strip().split('/')[0].strip().upper()
                    title = str(urlitem('div.project-list-content')('a:first').text()).strip()
                    url = str(urlitem('div.project-list-media')('a:first').attr('href')).strip()
                    avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javhoo:url:" + str(ex))
    if engine == 'torrentant':
        surls = []
        try:
            urldata = PyQuery(getHTML('http://m.torrentant.com/cn/search/' + keyword, 5, 5, 1, proxy))
            pagination = urldata('div.site-index')('div[class="search-pagination text-center"]')('ul.pagination:eq(1)')
            if pagination is not None and str(pagination).strip() != '':
                for subpage in pagination('li').items():
                    if re.match(r'^.*\d+.*$', str(subpage.text())):
                        if subpage('a').attr('href') is not None and str(subpage('a').attr('href')).strip() != '':
                            surls.append(parse.urljoin('http://m.torrentant.com/', subpage('a').attr('href')))
            else:
                surls.append('http://m.torrentant.com/cn/search/' + keyword)
            for surl in surls:
                urldata = PyQuery(getHTML(surl, 5, 5, 1, proxy))
                urlcontent = urldata('div.site-index')('div.co-md-12')
                for urlitem in urlcontent('div.movie-item-in').items():
                    code = str(urlitem('div.meta')('div.movie-tag').text()).split('/')[0].strip()
                    title = str(urlitem('a:first').attr('title')).strip()
                    url = parse.urljoin('http://m.torrentant.com/', str(urlitem('a:first').attr('href')).strip())
                    avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javhoo:url:" + str(ex))
    if engine == 'avmoo':
        surls = []
        try:
            urldata = PyQuery(getHTML('https://avmo.pw/cn/search/' + keyword, 5, 5, 1, proxy))
            pagination = urldata('ul[class="pagination pagination-lg mtb-0"]')
            if pagination is not None and str(pagination).strip() != '':
                for subpage in pagination('li').items():
                    if re.match(r'^.*\d+.*$', str(subpage.text())):
                        if subpage('a').attr('href') is not None and str(subpage('a').attr('href')).strip() != '':
                            surls.append(parse.urljoin('https://avmo.pw/', subpage('a').attr('href')))
            else:
                surls.append('https://avmo.pw/cn/search/' + keyword)
            for surl in surls:
                urldata = PyQuery(getHTML(surl, 5, 5, 1, proxy))
                urlcontent = urldata('div#waterfall')
                for urlitem in urlcontent('a.movie-box').items():
                    code = str(urlitem('div.photo-info')('date:first').text()).strip().upper()
                    urlitem('div.photo-info')('span')('date').remove()
                    title = str(urlitem('div.photo-info')('span').text().replace('/', '')).strip()
                    url = str(urlitem.attr('href')).strip()
                    avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javhoo:url:" + str(ex))
    return avpages


def avkeywordParse(textargs, type):
    lines = []
    keywords = []
    try:
        pattern = re.compile(r'[A-Za-z]{2,5}-?\d{2,4}|\d{6}[-_]\d{2,3}')
        if type == 'file':
            sfile = os.path.join(curDir(), textargs)
            chartype = charDetect(open(sfile, 'rb').read())
            for line in open(sfile, encoding=chartype):
                lines.append(line)
        elif type == 'url':
            lines.append(str(getHTML(textargs, 5, 5, 0)))
        else:
            for textarg in textargs.split(' '):
                lines.append(textarg)
                keywords.append(str(textarg).upper())

        for line in lines:
            for number in pattern.finditer(line):
                code = str(number.group()).upper()
                if re.match(r'\d{6}[-_]\d{2,3}', code):
                    code = code.replace('-', '_')
                else:
                    p = re.match(r'([A-Za-z]+).*?(\d+)', code)
                    code = str(p.group(1)) + '-' + str(p.group(2))
                if code not in keywords:
                    keywords.append(str(code))
    except Exception as ex:
        logging.error('avkeywordParse:' + str(ex))
    return keywords


def avinfoFetch(keywords, engine='javbus', proxy=''):
    avs = []
    title = ''
    issuedate = ''
    length = ''
    mosaic = ''
    director = ''
    manufacturer = ''
    publisher = ''
    series = ''
    category = ''
    actors = ''
    favor = ''
    coverlink = ''
    cover = b''
    link = ''

    for keyword in keywords:
        try:
            if engine == 'javbus':
                print('>' * 20 + ('Getting AV URLs For Keyword (' + keyword + ')').center(60) + '<' * 20)
                for avpage in avurlFetch(keyword, 'javbus', proxy):
                    # curl = 'https://www.javbus.com/' + avpage['code']
                    # curl = parse.urljoin('https://www.javbus.com/', avpage['code'])
                    # data = PyQuery('https://avmo.pw/cn/search/' + avpage['code'])
                    # url = str(data('a.movie-box').attr('href'))
                    # avs[avpage['code']] = url
                    print((' -- Fetching ' + avpage['code'] + ' -- ').center(100, '#'))
                    try:
                        data = PyQuery(getHTML(avpage['url'], 5, 5, 1, proxy))
                        content = data('div.container')
                        mosaic = str(data('ul[class="nav navbar-nav"]')('li[class="active"]').text()).strip()
                        mosaic = mosaic.replace('無', '无').replace('碼', '码').replace('修正', '码')
                        if mosaic == '码':
                            mosaic = '有码'
                        title = str(content('h3').eq(0).text().replace(avpage['code'], '')).strip()
                        avinfo = content('div[class="col-md-3 info"]')
                        '''
                        issuedate = str(re.search(r'\d*-\d*-\d*', avinfo('p:eq(1)').text()).group())
                        length = str(avinfo('p:eq(2)').text().split(' ')[-1]).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
                        director = str(avinfo('p:eq(3)')('a').text()).strip()
                        manufacturer = str(avinfo('p:eq(4)')('a').text()).strip()
                        publisher = str(avinfo('p:eq(5)')('a').text()).strip()
                        if avinfo('p:eq(6)')('a') is not None and str(avinfo('p:eq(6)')('a')).strip() != '':
                            series = str(avinfo('p:eq(6)')('a').text()).strip()
                            category = str(avinfo('p:eq(8)').text()).strip()
                            actors = str(avinfo('p:eq(10)').text()).strip()
                        else:
                            series = ''
                            category = str(avinfo('p:eq(7)').text()).strip()
                            actors = str(avinfo('p:eq(9)').text()).strip()
                        '''
                        for item in avinfo('p').items():
                            if '發行日' in item.text() or '発売日' in item.text() or '发行日' in item.text():
                                issuedate = str(re.search(r'\d*-\d*-\d*', item.text()).group())
                            if '長度' in item.text() or '時間' in item.text() or '长度' in item.text() or '时间' in item.text() or '时长' in item.text():
                                length = str(item.text().split(' ')[-1]).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
                            if '監督' in item.text() or '導演' in item.text() or '监督' in item.text() or '导演' in item.text():
                                director = str(item('a').text()).strip()
                            if 'メーカー' in item.text() or '製作商' in item.text() or '制作商' in item.text():
                                manufacturer = str(item('a').text()).strip()
                            if 'レーベル' in item.text() or '發行商' in item.text() or '发行商' in item.text():
                                publisher = str(item('a').text()).strip()
                            if 'シリーズ' in item.text() or '系列' in item.text():
                                series = str(item('a').text()).strip()
                            if 'ジャンル' in item.text() or '類別' in item.text() or '类别' in item.text():
                                category = str(item.next().text()).strip()
                            if '演員' in item.text() or '出演者' in item.text() or '演员' in item.text():
                                actors = str(item.next().text()).strip()
                        favor = '0'
                        coverlink = str(content('a.bigImage').attr('href')).strip()
                        print('番号:'.rjust(5) + avpage['code'] + '\n' + '标题:'.rjust(5) + title + '\n' + '日期:'.rjust(5) + issuedate + '\n' + '时长:'.rjust(5) + length + '\n' + '修正:'.rjust(5) + mosaic + '\n' + '导演:'.rjust(5) + director + '\n' + '制作:'.rjust(5) + manufacturer + '\n' + '发行:'.rjust(5) + publisher + '\n' + '系列:'.rjust(5) + series + '\n' + '类别:'.rjust(5) + category + '\n' + '女优:'.rjust(5) + actors + '\n' + '收藏:'.rjust(5) + favor + '\n' + '预览:'.rjust(5) + coverlink)

                        cover = getHTML(coverlink, 5, 5, 0, proxy)
                        # magnets = content('table#magnet-table')
                        try:
                            link = avlinkFilter(avlinkFetch(avpage['code'], 'zhongziso', proxy)).link
                        except Exception as ex:
                            logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                            link = 'page:' + avpage['url']
                        print('磁链:'.rjust(5) + link)
                        avs.append(av(avpage['code'], title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link))
                    except Exception as ex:
                        logging.warning('avinfoFetch:javbus:' + str(ex))
            if engine == 'javhoo':
                print('>' * 20 + ('Getting AV URLs For Keyword (' + keyword + ')').center(60) + '<' * 20)
                for avpage in avurlFetch(keyword, 'javhoo', proxy):
                    # curl = 'https://www.javhoo.com/av/' + avpage['code']
                    print((' -- Fetching ' + avpage['code'] + ' -- ').center(100, '#'))
                    try:
                        data = PyQuery(getHTML(avpage['url'], 5, 5, 1, proxy))
                        content = data('div#content')('div.wf-container')
                        avinfo = content('div.project_info')
                        mosaic = str(avinfo('span.category-link').text()).strip()
                        mosaic = mosaic.replace('無', '无').replace('碼', '码').replace('修正', '码')
                        if mosaic == '码':
                            mosaic = '有码'
                        title = str(data('h1[class="h3-size entry-title"]').text().replace(avpage['code'], '')).strip()
                        '''
                        issuedate = str(re.search(r'\d*-\d*-\d*', avinfo('p:eq(1)').text()).group())
                        length = str(avinfo('p:eq(2)').text().split(' ')[-1]).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
                        director = str(avinfo('p:eq(3)')('a').text()).strip()
                        manufacturer = str(avinfo('p:eq(4)')('a').text()).strip()
                        publisher = str(avinfo('p:eq(5)')('a').text()).strip()
                        series = str(avinfo('p:eq(6)')('a').text()).strip()
                        category = str(avinfo('p:eq(8)').text()).strip()
                        actors = str(avinfo('p:eq(10)').text()).strip()
                        '''
                        for item in avinfo('p').items():
                            if '發行日' in item.text() or '発売日' in item.text() or '发行日' in item.text():
                                issuedate = str(re.search(r'\d*-\d*-\d*', item.text()).group())
                            if '長度' in item.text() or '時間' in item.text() or '长度' in item.text() or '时间' in item.text() or '时长' in item.text():
                                length = str(item.text().split(' ')[-1]).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
                            if '監督' in item.text() or '導演' in item.text() or '监督' in item.text() or '导演' in item.text():
                                director = str(item('a').text()).strip()
                            if 'メーカー' in item.text() or '製作商' in item.text() or '制作商' in item.text():
                                manufacturer = str(item('a').text()).strip()
                            if 'レーベル' in item.text() or '發行商' in item.text() or '发行商' in item.text():
                                publisher = str(item('a').text()).strip()
                            if 'シリーズ' in item.text() or '系列' in item.text():
                                series = str(item('a').text()).strip()
                            if 'ジャンル' in item.text() or '類別' in item.text() or '类别' in item.text():
                                category = str(item.next().text()).strip()
                            if '演員' in item.text() or '出演者' in item.text() or '演员' in item.text():
                                actors = str(item.next().text()).strip()
                        favor = '0'
                        coverlink = str(content('div.project-content')('img[class="alignnone size-full"]').attr('src')).strip()
                        print('番号:'.rjust(5) + avpage['code'] + '\n' + '标题:'.rjust(5) + title + '\n' + '日期:'.rjust(5) + issuedate + '\n' + '时长:'.rjust(5) + length + '\n' + '修正:'.rjust(5) + mosaic + '\n' + '导演:'.rjust(5) + director + '\n' + '制作:'.rjust(5) + manufacturer + '\n' + '发行:'.rjust(5) + publisher + '\n' + '系列:'.rjust(5) + series + '\n' + '类别:'.rjust(5) + category + '\n' + '女优:'.rjust(5) + actors + '\n' + '收藏:'.rjust(5) + favor + '\n' + '预览:'.rjust(5) + coverlink)

                        cover = getHTML(coverlink, 5, 5, 0, proxy)
                        try:
                            link = avlinkFilter(avlinkFetch(avpage['code'], 'zhongziso', proxy)).link
                        except Exception as ex:
                            logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                            link = 'page:' + avpage['url']
                        print('磁链:'.rjust(5) + link)
                        avs.append(av(avpage['code'], title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link))
                    except Exception as ex:
                        logging.warning('avinfoFetch:javhoo:' + str(ex))
            if engine == 'torrentant':
                print('>' * 20 + ('Getting AV URLs For Keyword (' + keyword + ')').center(60) + '<' * 20)
                for avpage in avurlFetch(keyword, 'torrentant', proxy):
                    # curl = 'https://www.javhoo.com/av/' + avpage['code']
                    print((' -- Fetching ' + avpage['code'] + ' -- ').center(100, '#'))
                    try:
                        data = PyQuery(getHTML(avpage['url'], 5, 5, 1, proxy))
                        content = data('div.container-fluid')('div.movie-view:first')
                        avinfo = content('table.movie-view-table')
                        mosaic = '未知'
                        title = str(content('h1:first').text().replace(avpage['code'], '')).strip()
                        issuedate = str(avinfo('tr:eq(1)')('td:eq(1)').text()).strip()
                        length = str(avinfo('tr:eq(2)')('td:eq(1)').text()).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
                        director = str(avinfo('tr:eq(5)')('td:eq(1)').text()).strip()
                        manufacturer = str(avinfo('tr:eq(3)')('td:eq(1)').text()).strip()
                        publisher = str(avinfo('tr:eq(4)')('td:eq(1)').text()).strip()
                        series = str(avinfo('tr:eq(6)')('td:eq(1)').text()).strip()
                        category = str(content('div[class="col-md-12 tags"]').text()).strip()
                        actors = str(content('div#avatar-waterfall').text()).strip()
                        favor = '0'
                        coverlink = ''
                        print('番号:'.rjust(5) + avpage['code'] + '\n' + '标题:'.rjust(5) + title + '\n' + '日期:'.rjust(5) + issuedate + '\n' + '时长:'.rjust(5) + length + '\n' + '修正:'.rjust(5) + mosaic + '\n' + '导演:'.rjust(5) + director + '\n' + '制作:'.rjust(5) + manufacturer + '\n' + '发行:'.rjust(5) + publisher + '\n' + '系列:'.rjust(5) + series + '\n' + '类别:'.rjust(5) + category + '\n' + '女优:'.rjust(5) + actors + '\n' + '收藏:'.rjust(5) + favor + '\n' + '预览:'.rjust(5) + coverlink)

                        cover = b''
                        try:
                            link = avlinkFilter(avlinkFetch(avpage['code'], 'zhongziso', proxy)).link
                        except Exception as ex:
                            logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                            link = 'page:' + avpage['url']
                        print('磁链:'.rjust(5) + link)
                        avs.append(av(avpage['code'], title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link))
                    except Exception as ex:
                        logging.warning('avinfoFetch:torrentant:' + str(ex))
        except Exception as ex:
            logging.error('avinfoFetch:' + str(ex))
    return avs


def avlinkFetch(code, engine='btso', proxy=''):
    head = ''
    time = ''
    hot = ''
    size = ''
    clink = ''
    code = code.upper()
    avlinks = []
    try:
        if engine == 'btgongchang':
            data = PyQuery(getHTML('http://btgongchang.org/search/' + code + '-first-asc-1', 5, 5, 1, proxy))
            content = data('table[class="data mb20"]')
            items = content('tr:gt(0)')
            for item in items.items():
                try:
                    head = str(item('td:eq(0)')('div.item-title')('a').text()).strip()
                    head = re.sub(r'\/\*.*\*\/', '', head)
                    time = str(item('td:eq(1)').text()).strip()
                    hot = str(item('td:eq(2)').text()).strip()
                    size = str(item('td:eq(3)').text()).strip().lower()
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    clink = item('td:eq(4)')('a:first').attr('href')
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:btgongchang:' + str(ex))
        if engine == 'btso':
            data = PyQuery(getHTML('https://btso.pw/search/' + code + '/', 5, 5, 1, proxy))
            content = data('div.data-list')
            items = content('div[class="row"]')
            for item in items.items():
                try:
                    head = str(item('a').attr('title')).strip()
                    head = re.sub(r'\/\*.*\*\/', '', head)
                    time = str(item('div[class="col-sm-2 col-lg-2 hidden-xs text-right date"]').text()).strip()
                    size = str(item('div[class="col-sm-2 col-lg-1 hidden-xs text-right size"]').text()).strip().lower()
                    hot = '100'
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    tmplink = str(item('a').attr('href')).strip()
                    tmplink = parse.urljoin('https://btso.pw/', tmplink)
                    linkdata = PyQuery(getHTML(tmplink, 5, 5, 1, proxy))
                    clink = str(linkdata('textarea#magnetLink').text()).strip()
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:btso:' + str(ex))
        if engine == 'btdb':
            data = PyQuery(getHTML('https://btdb.in/q/' + code + '/?sort=popular', 5, 5, 1, proxy))
            content = data('ul.search-ret-list')
            items = content('li.search-ret-item')
            for item in items.items():
                try:
                    head = str(item('h2.item-title')('a').attr('title')).strip()
                    head = re.sub(r'\/\*.*\*\/', '', head)
                    # head = head.encode('latin-1').decode('utf-8')
                    linkinfo = item('div.item-meta-info')
                    clink = str(linkinfo('a.magnet').attr('href')).strip()
                    linkinfodata = str(linkinfo.text()).lower()
                    linkinfos = re.match(r'^.*size:(.*)files:(.*)addtime:(.*)popularity:(.*)$', linkinfodata).groups()
                    size = str(linkinfos[0]).strip().replace('  ', ' ')
                    time = str(linkinfos[2]).strip()
                    hot = str(linkinfos[3]).strip()
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:btdb:' + str(ex))
        if engine == 'torrentant':
            data = PyQuery(getHTML('http://www.torrentant.com/cn/s/' + code + '?sort=hot', 5, 5, 1, proxy))
            content = data('ul[class="search-container"]')
            items = content('li[class="search-item clearfix"]')
            for item in items.items():
                try:
                    head = str(item('div[class="search-content text-left"]')('h2')('a').attr('title')).strip()
                    head = re.sub(r'\/\*.*\*\/', '', head)
                    linkinfo = item('div[class="search-content text-left"]')('div[class="resultsContent"]')('p[class="resultsIntroduction"]')
                    size = str(linkinfo('label').eq(1).text()).strip().lower()
                    hot = str(linkinfo('label').eq(2).text()).strip()
                    time = str(linkdata('table[class="table table-hover"]')('tbody')('tr').eq(0)('td').eq(0).text()).strip()
                    clink = str(linkdata('a[class="btn btn-warning"]').attr('href')).strip()
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    tmplink = str(item('div[class="search-content text-left"]')('h2')('a').attr('href')).strip()
                    tmplink = parse.urljoin('http://www.torrentant.com/', tmplink)
                    linkdata = PyQuery(getHTML(tmplink, 5, 5, 1, proxy))
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:torrentant:' + str(ex))
        if engine == 'javhoo':
            data = PyQuery(getHTML('https://www.javhoo.com/av/' + code + '/', 5, 5, 1, proxy))
            content = data('table#magnet-table')
            items = content('tr:gt(0)')
            for item in items.items():
                try:
                    head = str(item('td:eq(0)')('a').text()).strip()
                    head = re.sub(r'\/\*.*\*\/', '', head)
                    size = str(item('td:eq(1)')('a').text()).strip().lower()
                    time = str(item('td:eq(2)')('a').text()).strip()
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    hot = '100'
                    clink = str(item('td:eq(0)')('a').attr('href')).strip()
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:javhoo:' + str(ex))
        if engine == 'zhongziso':
            data = PyQuery(getHTML('http://www.zhongziso.com/list/' + code + '/1', 5, 5, 1, proxy))
            content = data('div.inerTop')
            items = content('table[class="table table-bordered table-striped"]')
            for item in items.items():
                try:
                    head = str(item('tr:eq(0)')('div.text-left').text()).strip()
                    head = re.sub(r'\/\*.*\*\/', '', head)
                    time = str(item('tr:eq(1)')('td:eq(0)')('strong:first').text()).strip()
                    size = str(item('tr:eq(1)')('td:eq(1)')('strong:first').text()).strip().lower()
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    hot = str(item('tr:eq(1)')('td:eq(2)')('strong:first').text()).strip()
                    clink = str(item('tr:eq(1)')('td:eq(3)')('a').attr('href')).strip()
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:zhongziso:' + str(ex))
        '''
        if engine == 'btago':
            data = PyQuery(getHTML('http://www.btago.com/e/' + code + '/', 5, 5, 1, proxy))
            content = data('div#container')('div.listLoader')
            items = content('div.item')
            for item in items.items():
                try:
                    head = str(item('div.t').text()).strip()
                    linkinfo = str(item('div.info').text()).strip()
                    print(linkinfo)
                    tmplink = str(item('div.t')('a:first').attr('href')).strip()
                    tmplink = parse.urljoin('http://www.btago.com/', tmplink)
                    time = str(linkinfo.split('|')[2].split('：')[1]).strip()
                    hot = '100'
                    size = str(linkinfo.split('|')[0].split('：')[1]).strip().lower()
                    if 'g' in size:
                        size = str(size.replace('gb', '').replace('g', '').strip())
                    elif 'm' in size:
                        size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                    elif 'k' in size:
                        size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                    else:
                        size = str(size).replace('b', '').strip()
                    size = str("%.2f" % float(size))
                    linkdata = PyQuery(getHTML(tmplink, 5, 5, 1, proxy))
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:btago:' + str(ex))
        '''
    except Exception as ex:
        logging.warning('avlinkFetch:' + str(ex))
    return avlinks


def avlinkSort(cav):
    rate = 100
    if len(cav.title) < 20:
        rate -= 20
    elif len(cav.title) < 80:
        rate -= 50

    if float(cav.size) < 1:
        rate -= 30
    elif float(cav.size) < 2:
        rate -= 10
    elif float(cav.size) > 5:
        rate -= 60

    rate = int(cav.hot) * rate * (1 + float(cav.size) / 10) * (1 + float(len(cav.title)) / 300) / 100
    return rate


def avlinkFilter(avlinks):
    '''
    selectedav = None
    for cavlink in avlinks:
        if selectedav is None or int(cavlink.hot) > int(selectedav.hot):
            selectedav = cavlink
    return selectedav
    '''
    return sorted(avlinks, key=avlinkSort)[-1]


def av2file(avs, dirpath):
    txtfs = None
    txtname = 'avinfos.txt'
    try:
        print('Saving AV Infomation to Files'.center(100, '*'))
        txtpath = os.path.join(dirpath, txtname)
        txtfs = open(txtpath, 'a', encoding='utf8')
        for cav in avs:
            try:
                print('Creating AV Information : ' + cav.title, end=' ...... ')
                txtfs.write('番号:'.center(5) + cav.code + '\n')
                txtfs.write('标题:'.center(5) + cav.title + '\n')
                txtfs.write('日期:'.center(5) + cav.issuedate + '\n')
                txtfs.write('时长:'.center(5) + cav.length + '\n')
                txtfs.write('修正:'.center(5) + cav.mosaic + '\n')
                txtfs.write('导演:'.center(5) + cav.director + '\n')
                txtfs.write('制作:'.center(5) + cav.manufacturer + '\n')
                txtfs.write('发行:'.center(5) + cav.publisher + '\n')
                txtfs.write('系列:'.center(5) + cav.series + '\n')
                txtfs.write('类别:'.center(5) + cav.category + '\n')
                txtfs.write('女优:'.center(5) + cav.actors + '\n')
                txtfs.write('收藏:'.center(5) + cav.favor + '\n')
                txtfs.write('预览:'.center(5) + cav.coverlink + '\n')
                txtfs.write('磁链:'.center(5) + cav.link + '\n')
                txtfs.write('#' * 100 + '\n')
                ext = cav.coverlink.split('.')[-1] if '.' in cav.coverlink else 'jpg'
                imgname = cav.code + '_' + cav.title + '.' + ext
                imgname = imgname.replace('<', '').replace('>', '').replace('/', '').replace('\\', '').replace('|', '').replace(':', '').replace('"', '').replace('*', '').replace('?', '')
                imgpath = os.path.join(dirpath, imgname)
                imgfs = open(imgpath, 'wb')
                # imgfs.write(getHTML(cav.coverlink, 5, 5, 0, proxy))
                imgfs.write(cav.cover)
                imgfs.close()
                print('READY')
            except Exception as ex:
                logging.debug('av2file:' + str(ex))
                print('FAILED')
        txtfs.close()
        print('COMPLETE')
    except Exception as ex:
        if txtfs is not None:
            txtfs.close()
        logging.error('av2file:' + str(ex))
        print('FAILED')


def av2db(avs, dirpath):
    data = []
    dbname = 'avinfos.db'
    try:
        print('Saving AV Infomation to Database'.center(100, '*'))
        dbpath = os.path.join(dirpath, dbname)
        sql = '''CREATE TABLE  IF NOT EXISTS `av` (
              `code` varchar(100) NOT NULL,
              `title` varchar(500) NOT NULL,
              `issuedate` varchar(100) DEFAULT NULL,
              `length` varchar(100) DEFAULT NULL,
              `mosaic` varchar(100) DEFAULT NULL,
              `director` varchar(100) DEFAULT NULL,
              `manufacturer` varchar(100) DEFAULT NULL,
              `publisher` varchar(100) DEFAULT NULL,
              `series` varchar(100) DEFAULT NULL,
              `category` varchar(500) DEFAULT NULL,
              `actors` varchar(500) DEFAULT NULL,
              `favor` varchar(20) DEFAULT '0',
              `coverlink` varchar(300) DEFAULT NULL,
              `cover` BLOB DEFAULT NULL,
              `link` varchar(10000) DEFAULT NUll,
               PRIMARY KEY (`code`)
            )'''
        # conn = get_conn(dbpath)
        create_table(get_conn(dbpath), sql)

        sql = '''INSERT OR IGNORE INTO av values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'''
        for cav in avs:
            try:
                print('Creating AV Information : ' + cav.title, end=' ...... ')
                # data.append((cav.code, cav.title, cav.issuedate, cav.length, cav.mosaic, cav.director, cav.manufacturer, cav.publisher, cav.series, cav.category, cav.actors, cav.favor, cav.coverlink, sqlite3.Binary(bytes(cav.cover)), cav.link))
                save(get_conn(dbpath), sql, [(cav.code, cav.title, cav.issuedate, cav.length, cav.mosaic, cav.director, cav.manufacturer, cav.publisher, cav.series, cav.category, cav.actors, cav.favor, cav.coverlink, sqlite3.Binary(cav.cover), cav.link)])
                print('READY')
            except Exception as ex:
                logging.debug('av2db:' + str(ex))
                print('FAILED')
        print('COMPLETE')
    except Exception as ex:
        if 'UNIQUE constraint failed' in str(ex):
            logging.debug('av2db:' + str(ex))
            print('COMPLETE')
        else:
            logging.error('av2db:' + str(ex))
            print('FAILED')


def avsave(avs, savetype='file', tpath=curDir()):
    if avs is not None and len(avs) > 0:
        try:
            dirpath = os.path.join(curDir(), tpath)
            if not os.path.isdir(dirpath):
                os.mkdir(dirpath)
            if savetype.lower() == 'file':
                av2file(avs, dirpath)
            elif savetype.lower() == 'db':
                av2db(avs, dirpath)
            elif savetype.lower() == 'both':
                av2file(avs, dirpath)
                av2db(avs, dirpath)

        except Exception as ex:
            logging.error('avsave:' + str(ex))
    else:
        logging.error('No AV Infomation')


def avquickFetch(code, proxy=''):
    title = ''
    issuedate = ''
    length = ''
    mosaic = ''
    director = ''
    manufacturer = ''
    publisher = ''
    series = ''
    category = ''
    actors = ''
    favor = ''
    coverlink = ''
    links = []
    link = ''
    code = code.upper()
    # curl = 'https://www.javhoo.com/av/' + avpage['code']
    try:
        data = PyQuery(getHTML('https://www.javhoo.com/av/' + code, 5, 5, 1, proxy))
        content = data('div#content')('div.wf-container')
        avinfo = content('div.project_info')
        mosaic = str(avinfo('span.category-link').text()).strip()
        mosaic = mosaic.replace('無', '无').replace('碼', '码').replace('修正', '码')
        if mosaic == '码':
            mosaic = '有码'
        title = str(data('h1[class="h3-size entry-title"]').text().replace(code, '')).strip()
        for item in avinfo('p').items():
            if '發行日' in item.text() or '発売日' in item.text() or '发行日' in item.text():
                issuedate = str(re.search(r'\d*-\d*-\d*', item.text()).group())
            if '長度' in item.text() or '時間' in item.text() or '长度' in item.text() or '时间' in item.text() or '时长' in item.text():
                length = str(item.text().split(' ')[-1]).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
            if '監督' in item.text() or '導演' in item.text() or '监督' in item.text() or '导演' in item.text():
                director = str(item('a').text()).strip()
            if 'メーカー' in item.text() or '製作商' in item.text() or '制作商' in item.text():
                manufacturer = str(item('a').text()).strip()
            if 'レーベル' in item.text() or '發行商' in item.text() or '发行商' in item.text():
                publisher = str(item('a').text()).strip()
            if 'シリーズ' in item.text() or '系列' in item.text():
                series = str(item('a').text()).strip()
            if 'ジャンル' in item.text() or '類別' in item.text() or '类别' in item.text():
                category = str(item.next().text()).strip()
            if '演員' in item.text() or '出演者' in item.text() or '演员' in item.text():
                actors = str(item.next().text()).strip()
        favor = '0'
        coverlink = str(content('div.project-content')('img[class="alignnone size-full"]').attr('src')).strip()
        cover = getHTML(coverlink, 5, 5, 0, proxy)

        linkcontent = data('table#magnet-table')
        linkinfo = linkcontent('tr:gt(0)')
        for linkitem in linkinfo.items():
            try:
                head = str(linkitem('td:eq(0)')('a').text()).strip()
                size = str(linkitem('td:eq(1)')('a').text()).strip().lower()
                time = str(linkitem('td:eq(2)')('a').text()).strip()
                if 'g' in size:
                    size = str(size.replace('gb', '').replace('g', '').strip())
                elif 'm' in size:
                    size = str(float(size.replace('mb', '').replace('m', '').strip()) / 1024)
                elif 'k' in size:
                    size = str(float(size.replace('kb', '').replace('k', '').strip()) / 1024 / 1024)
                else:
                    size = str(size).strip()
                size = str("%.2f" % float(size))
                hot = '100'
                clink = str(linkitem('td:eq(0)')('a').attr('href')).strip()
                links.append(avlinkinfo(code, head, time, hot, size, clink, 'javhoo'))
            except Exception as ex:
                logging.debug('avlinkFetch:javhoo:' + str(ex))
        try:
            link = avlinkFilter(links).link
        except Exception as ex:
            logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
            link = 'page:' + 'https://www.javhoo.com/av/' + code
        return av(code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link)
    except Exception as ex:
        logging.error('avquickFetch:' + str(ex))
        return None


def main(argv):
    texts = []
    stype = 'file'
    surl = ''
    sfile = ''
    tpath = curDir()
    sengine = ''
    sproxy = ''
    keywords = []
    textwords = []
    filewords = []
    urlwords = []
    logInit()

    if argv is not None and len(argv) > 0:
        try:
            opts, args = getopt.getopt(argv, "hd:e:t:p:u:f:s:", ["dir=", "engine=", "type=", "proxy=", "url=", "file=", "code="])
        except getopt.GetoptError:
            print(
                '''Usage: avfetch.py [-d <targetpath>] [-e <engine>] [-t <savetype>] [-p <proxy>] [-u <url>] [-f <filename>] [-s <codes>] [<codes>]\n
                Example: avfetch.py -d D:/ -e javbus -t file -p socks5@127.0.0.1:1080 -u http://www.baidu.com -f a.txt -s ABP-563 SRS-064 SNIS-862'''
            )
            exit(2)

        if len(args) > 0:
            texts.extend(args)
        for opt, arg in opts:
            if opt == '-h':
                print(
                    '''Usage: avfetch.py [-d <targetpath>] [-e <engine>] [-t <savetype>] [-p <proxy>] [-u <url>] [-f <filename>] [-s <codes>] [<codes>]\n
                    Example: avfetch.py -d D:/ -e javbus -t file -p socks5@127.0.0.1:1080 -u http://www.baidu.com -f a.txt -s ABP-563 SRS-064 SNIS-862'''
                )
                exit()
            elif opt in ("-d", "--dir"):
                tpath = arg
            elif opt in ("-e", "--engine"):
                sengine = arg
            elif opt in ("-t", "--type"):
                stype = arg
            elif opt in ("-p", "--proxy"):
                sproxy = arg
                if not re.match(r'^.+@.+:.+$', sproxy, flags=0):
                    print('proxy format is illegal!')
                    exit(2)
            elif opt in ("-u", "--url"):
                surl = arg
            elif opt in ("-f", "--file"):
                sfile = os.path.join(curDir(), arg)
                if not os.path.isfile(sfile):
                    print('file is not exist!')
                    exit(2)
            elif opt in ("-s", "--code"):
                texts.append(arg)
            else:
                pass
        try:
            if len(texts) > 0:
                textwords = avkeywordParse(' '.join(texts), 'code')
            if sfile != '':
                filewords = avkeywordParse(sfile, 'file')
            if surl != '':
                urlwords = avkeywordParse(surl, 'url')
            keywords.extend(textwords)
            keywords.extend(filewords)
            keywords.extend(urlwords)
            avs = avinfoFetch(keywords, engine=sengine, proxy=sproxy)
            avsave(avs, stype, tpath)

        except Exception as ex:
            logging.error('main:' + str(ex))


if __name__ == "__main__":
    main(sys.argv[1:])

main(['-d', 'C:/Users/xshrim/Desktop/imgsss', '-e', 'javbus', '-t', 'both', '-s', 'ipz-137', 'FSET-337'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgsss', '-e', 'javbus', '-t', 'both', '-s', 'ipz-137', 'ipz-371 midd-791 fset-337 sw-140'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javhoo', '-t', 'file', '-s', '天海つばさ'])
# main(['-d', 'imgss', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-u', 'http://btgongchang.org/'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgs', '-e', 'javbus', '-t', 'db', '-s', 'IPZ-137', 'IPZ820 MDS-825 FSET-337 F-123 FS-1'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'both', '-f', 'C:/Users/xshrim/Desktop/av.txt'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'file', '-s', 'IPZ-137', 'IPZ820 MDS-825 FSET-337 F-123 FS-1'])
# print(avquickFetch('ipz-371'))

# for cav in avlinkFetch('ipz-371', 'zhongziso'):
#    print(cav)
# print(avlinkFilter(avlinkFetch('ipz-101', 'btso')).title)

# 搜索引擎：
# btso:https://btso.pw/search/ipz-137/
# javhoo:https://www.javhoo.com/av/ipz-137/
# btdb:https://btdb.in/q/ipz-137/
# sukebei:https://sukebei.nyaa.se/?page=search&term=ipz-137&sort=4 (only torrent)
# torrentant:http://www.torrentant.com/cn/s/ipz-137?sort=hot (inaccuracy)
# btgongchang:http://btgongchang.org/search/MDS-825-first-asc-1
# zhongziso:http://www.zhongziso.com/list/ipz-137/1
# btago:http://www.btago.com/e/ipz-371/
