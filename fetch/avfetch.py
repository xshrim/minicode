#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import concurrent.futures
import getopt
import io
import logging
import os
import re
import socket
import sqlite3
import sys
import threading
import time
import zlib
from urllib import parse, request
from urllib.parse import quote

import chardet
import html2text
import pyperclip
import socks
from pyquery import PyQuery

global InfoEngine, BTEngine
InfoEngine = None
BTEngine = None
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
        print('番号:'.center(5) + self.code)
        print('标题:'.center(5) + self.title)
        print('日期:'.center(5) + self.issuedate)
        print('时长:'.center(5) + self.length)
        print('修正:'.center(5) + self.mosaic)
        print('导演:'.center(5) + self.director)
        print('制作:'.center(5) + self.manufacturer)
        print('发行:'.center(5) + self.publisher)
        print('系列:'.center(5) + self.series)
        print('类别:'.center(5) + self.category)
        print('女优:'.center(5) + self.actors)
        print('收藏:'.center(5) + self.favor)
        print('预览:'.center(5) + self.coverlink)
        print('磁链:'.center(5) + self.link)
        '''
        return '番号:'.center(5) + self.code + '\n' + '标题:'.center(5) + self.title + '\n' + '日期:'.center(5) + self.issuedate + '\n' + '时长:'.center(5) + self.length + '\n' + '修正:'.center(5) + self.mosaic + '\n' + '导演:'.center(5) + self.director + '\n' + '制作:'.center(5) + self.manufacturer + '\n' + '发行:'.center(5) + self.publisher + '\n' + '系列:'.center(5) + self.series + '\n' + '类别:'.center(5) + self.category + '\n' + '女优:'.center(5) + self.actors + '\n' + '收藏:'.center(5) + self.favor + '\n' + '预览:'.center(5) + self.coverlink + '\n' + '磁链:'.center(5) + self.link

    __repr__ = __str__


class actres(object):
    '''
    jname = ''
    cname = ''
    pname = ''
    aname = ''
    birthday = ''
    birthplace = ''
    cup = ''
    stime = ''
    status = ''
    favorate = ''
    codes = ''
    ipage = ''
    notes = ''
    '''

    def __init__(self, jname, cname, pname, aname, birthday, birthplace, cup, stime, status, favorate, codes, ipage, notes):
        self.jname = jname
        self.cname = cname
        self.pname = pname
        self.aname = aname
        self.birthday = birthday
        self.birthplace = birthplace
        self.cup = cup
        self.stime = stime
        self.status = status
        self.favorate = favorate
        self.codes = codes
        self.ipage = ipage
        self.notes = notes

    def __str__(self):
        return '日文原名:'.center(7) + self.jname + '\n' + '中文译名:'.center(7) + self.jname + '\n' + '平片假名:'.center(7) + self.pname + '\n' + '曾用别名:'.center(7) + self.aname + '\n' + '出生日期:'.center(7) + self.birthday + '\n' + '出生地区:'.center(7) + self.birthplace + '\n' + '胸部罩杯:'.center(7) + self.cup + '\n' + '出道时间:'.center(7) + self.stime + '\n' + '当前状态:'.center(7) + self.status + '\n' + '兴趣爱好:'.center(7) + self.favorate + '\n' + '作品代码:'.center(7) + self.codes + '\n' + '简介主页:'.center(7) + self.ipage + '\n' + '备注信息:'.center(7) + self.notes


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
        return '来源:'.center(5) + self.origin + '\n' + '番号:'.center(5) + self.code + '\n' + '标题:'.center(5) + self.title + '\n' + '日期:'.center(5) + self.time + '\n' + '热度:'.center(5) + self.hot + '\n' + '大小:'.center(5) + self.size + ' GB\n' + '磁链:'.center(5) + self.link
        # return self.code + ' -- ' + self.title + ' -- ' + self.time + ' -- ' + self.hot + ' -- ' + self.size + ' -- ' + self.link + ' -- ' + self.origin

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
        return r
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))
    return None


def fetchone(conn, sql, data):
    if sql is not None and sql != '':
        if data is not None:
            # Do this instead
            d = (data,)
            cu = get_cursor(conn)
            # print('执行sql:[{}],参数:[{}]'.format(sql, data))
            cu.execute(sql, d)
            r = cu.fetchall()
            return r
        else:
            logging.error('the [{}] equal None!'.format(data))
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))
    return None


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


def duplicateCheck(code, dbfile):
    sql = 'SELECT code FROM av WHERE code = ?'
    res = fetchone(get_conn(dbfile), sql, code)
    if res is not None and len(res) > 0:
        return True
    else:
        return False


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


def choseEngine(type='ie', ielist=[('javbus', 'https://www.javbus.com'), ('javhoo', 'https://www.javhoo.com')], belist=[('zhongzilou', 'https://www.zhongzilou.com/'), ('btso', 'https://btso.pw/')], proxy=''):
    if type == 'ie':
        for ie in ielist:
            if detectPage(url=ie[1], proxy=proxy):
                logging.info('Using ' + ie[0] + ' Info Engine: ' + ie[1])
                return ie[0]
    if type == 'be':
        for be in belist:
            if detectPage(url=be[1], proxy=proxy):
                logging.info('Using ' + be[0] + ' BT Engine: ' + be[1])
                return be[0]
    return None


def detectPage(url, timeout=2, retry=2, sleep=0, proxy=''):
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
    headers = [('Host', 'img0.imgtn.bdimg.com'), ('Connection', 'close'), ('Cache-Control', 'max-age=0'), ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'), ('User-Agent', 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'), ('Accept-Encoding', '*'), ('Accept-Language', 'zh-CN,zh,en-US,en,*;q=0.8'), ('If-None-Match', '90101f995236651aa74454922de2ad74'), ('Referer', 'http://image.baidu.com/'), ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

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
            logging.debug('getHTML:' + str(ex))
            if '403' in str(ex) or '404' in str(ex) or '502' in str(ex) or '11001'in str(ex):
                return False
        i -= 1
    return False


def render(curl):
    """Fully render HTML, JavaScript and all."""
    import sys
    from PyQt5.QtCore import (QUrl, QEventLoop)
    from PyQt5.QtWidgets import QApplication
    from PyQt5.QtWebEngineWidgets import QWebEngineView

    class Render(QWebEngineView):
        def __init__(self, iurl):
            self.html = None
            self.app = QApplication(sys.argv)
            QWebEngineView.__init__(self)
            self.loadFinished.connect(self._loadFinished)
            self.load(QUrl(iurl))
            while self.html is None:
                self.app.processEvents(QEventLoop.ExcludeUserInputEvents | QEventLoop.ExcludeSocketNotifiers | QEventLoop.WaitForMoreEvents)
            self.app.quit()

        def _callable(self, data):
            self.html = data

        def _loadFinished(self, result):
            self.page().toHtml(self._callable)

    return Render(curl).html


def getHTML(url, timeout=5, retry=3, sleep=0, proxy=''):
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
            opener.close()
            logging.debug('getHTML:' + str(ex))
            if '403' in str(ex) or '404' in str(ex) or '502' in str(ex) or '11001'in str(ex):
                break
        i -= 1
    return contents


# def avpageFetch(url, engine='javbus', proxy=''):
def avpageFetch(url, engine=InfoEngine, proxy=''):
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
            logging.warning("avpageFetch:javbus:" + str(ex))
    if engine == 'javhoo':
        try:
            curl = url
            urldata = PyQuery(getHTML(curl, 5, 5, 1, proxy))
            urlcontent = urldata('div#content')
            for urlitem in urlcontent('div[class="wf-cell iso-item"]').items():
                code = str(urlitem('div.project-list-content')('date:first').text()).strip().split('/')[0].strip().upper()
                title = str(urlitem('div.project-list-content')('a:first').text()).strip()
                url = str(urlitem('div.project-list-media')('a:first').attr('href')).strip()
                print(code.ljust(12) + ' -> ' + title + ' -> ' + url)
                avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avpageFetch:javhoo:" + str(ex))
    return avpage


# def avurlFetch(keyword, engine='javbus', proxy=''):
def avurlFetch(keyword, engine=InfoEngine, proxy=''):
    avpages = []
    if engine == 'javbus':
        pidx = 1
        baseurl = 'https://www.javbus.com/search/' + keyword + '/'
        while True:
            try:
                curl = baseurl + str(pidx)
                print(('Parsing Page (' + curl + ')').center(100, '-'))
                urldata = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                urlcontent = urldata('div#waterfall')
                for urlitem in urlcontent('a.movie-box').items():
                    code = str(urlitem('div.photo-info')('date:first').text()).strip().upper()
                    urlitem('div.photo-info')('span')('date').remove()
                    title = str(urlitem('div.photo-info')('span').text().replace('/', '')).strip()
                    url = str(urlitem.attr('href')).strip()
                    print(code.ljust(12) + ' -> ' + title + ' -> ' + url)
                    avpages.append({'code': code, 'title': title, 'url': url})

                pagination = urldata('ul[class="pagination pagination-lg"]')
                if pagination is not None and str(pagination).strip() != '':
                    pginfo = ''
                    for item in pagination('li').items():
                        if str(pidx) == item.text().strip():
                            pginfo += '[' + str(pidx) + '] '
                        else:
                            pginfo += item.text().strip() + ' '
                    print(pginfo)
                    if '下一页' in pagination.text() or '下一頁' in pagination.text():
                        pidx += 1
                    else:
                        break
                else:
                    break
            except Exception as ex:
                logging.warning("avurlFetch:javbus:" + str(ex))
                break
    if engine == 'javhoo':
        pidx = 1
        baseurl = 'https://www.javhoo.com/search/' + keyword + '/page/'
        while True:
            try:
                curl = baseurl + str(pidx)
                print(('Parsing Page (' + curl + ')').center(100, '-'))
                urldata = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                urlcontent = urldata('div#content')
                for urlitem in urlcontent('div[class="wf-cell iso-item"]').items():
                    code = str(urlitem('div.project-list-content')('date:first').text()).strip().split('/')[0].strip().upper()
                    title = str(urlitem('div.project-list-content')('a:first').text()).strip()
                    url = str(urlitem('div.project-list-media')('a:first').attr('href')).strip()
                    print(code.ljust(12) + ' -> ' + title + ' -> ' + url)
                    avpages.append({'code': code, 'title': title, 'url': url})

                pagination = urldata('ul[class="pagination pagination-lg"]')
                if pagination is not None and str(pagination).strip() != '':
                    pginfo = ''
                    for item in pagination('li').items():
                        if str(pidx) == item.text().strip():
                            pginfo += '[' + str(pidx) + '] '
                        else:
                            pginfo += item.text().strip() + ' '
                    print(pginfo)
                    if '下一页' in pagination.text() or '下一頁' in pagination.text():
                        pidx += 1
                    else:
                        break
                else:
                    break
            except Exception as ex:
                logging.warning("avurlFetch:javbus:" + str(ex))
                break
    '''
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
                    print(code.ljust(12) + ' -> ' + title + ' -> ' + url)
                    avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javhoo:url:" + str(ex))
    '''
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
                    print(code.ljust(12) + ' -> ' + title + ' -> ' + url)
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
                    print(code.ljust(12) + ' -> ' + title + ' -> ' + url)
                    avpages.append({'code': code, 'title': title, 'url': url})
        except Exception as ex:
            logging.warning("avurlFetch:javhoo:url:" + str(ex))
    return avpages


# def avactresFetch(atype, engine='javbus', proxy=''):
def avactresFetch(atype, engine=InfoEngine, proxy=''):
    avactress = []
    pidx = 1
    if engine == 'javbus':
        if atype == '有码' or atype == '有碼' or atype == 'censored':
            baseurl = 'https://www.javbus.com/actresses/'
        else:
            baseurl = 'https://www.javbus.com/uncensored/actresses/'
        while True:
            try:
                curl = baseurl + str(pidx)
                print(('Parsing Page (' + curl + ')').center(100, '-'))
                urldata = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                urlcontent = urldata('div#waterfall')
                for urlitem in urlcontent('a[class="avatar-box text-center"]').items():
                    name = urlitem.text().strip()
                    link = urlitem.attr('href').strip()
                    print(name.ljust(15) + ' -> ' + link)
                    avactress.append((name, link))

                pagination = urldata('ul[class="pagination pagination-lg"]')
                if pagination is not None and str(pagination).strip() != '':
                    pginfo = ''
                    for item in pagination('li').items():
                        if str(pidx) == item.text().strip():
                            pginfo += '[' + str(pidx) + '] '
                        else:
                            pginfo += item.text().strip() + ' '
                    print(pginfo)
                    if '下一页' in pagination.text() or '下一頁' in pagination.text():
                        pidx += 1
                    else:
                        break
                else:
                    break
            except Exception as ex:
                logging.warning("avactresFetch:" + str(ex))
                break
    return avactress


def avkeywordParse(textargs, type):
    lines = []
    keywords = []
    try:
        pattern = re.compile(r'[A-Za-z]{1,7}-?[A-Za-z]?\d{2,4}-?\d{0,3}|\d{6}[-_]\d{4}[-_]\d{2}|\d{6}[-_]\d{2,3}|\d{6}-[A-Za-z]{3,6}|[A-Za-z]{1,3}\d[A-Za-z]{1,3}-\d{2,4}')
        if type == 'file':
            sfile = os.path.join(curDir(), textargs)
            chartype = charDetect(open(sfile, 'rb').read())
            for line in open(sfile, encoding=chartype):
                lines.append(line)
        elif type == 'url':
            data = html2text.html2text(PyQuery(getHTML(textargs)).html())
            for item in data.split(' '):
                lines.append(item)
        else:
            for textarg in textargs.split(' '):
                lines.append(textarg)
                keywords.append(str(textarg).upper())

        while len(lines) > 0:
            for number in pattern.finditer(lines.pop(0)):
                code = str(number.group()).upper()
                if re.match(r'^\d{6}_\d{2,3}$', code):
                    lines.append(code.replace('_', '-'))
                elif re.match(r'([A-Za-z]+)(\d+)', code):
                    p = re.match(r'([A-Za-z]+)(\d+)', code)
                    lines.append(str(p.group(1)) + '-' + str(p.group(2)))
                if code not in keywords:
                    keywords.append(str(code))
        '''
        for line in lines:
            for number in pattern.finditer(line):
                code = str(number.group()).upper()
                if re.match(r'\d{6}[-_]\d{2,3}', code):
                    if code not in keywords:
                        keywords.append(str(code))
                    code = code.replace('-', '_')
                    if code not in keywords:
                        keywords.append(str(code))
                elif re.match(r'([A-Za-z]+).*?(\d+)', code):
                    if code not in keywords:
                        keywords.append(str(code))
                    p = re.match(r'([A-Za-z]+).*?(\d+)', code)
                    code = str(p.group(1)) + '-' + str(p.group(2))
                    if code not in keywords:
                        keywords.append(str(code))
                if code not in keywords:
                    keywords.append(str(code))
        '''
    except Exception as ex:
        logging.error('avkeywordParse:' + str(ex))
    return keywords


# def avfullFetch(keywords, stype, tpath, mthread=10, engine='javbus', proxy='', dbfile=None):
def avfullFetch(keywords, stype, tpath, mthread=10, engine=InfoEngine, proxy='', dbfile=None):
    avs = []
    try:
        print((' [ Collecting Information For ' + str(len(keywords)) + ' keywords ] ').center(100, '/'))
        for keyword in keywords:
            print('// ' + keyword)
        print('/' * 100)
        avs = []
        avpages = []
        with concurrent.futures.ThreadPoolExecutor(max_workers=mthread) as executor:
            tasks = []
            for keyword in keywords:
                print('>' * 20 + ('Getting AV URLs For Keyword (' + keyword + ')').center(60) + '<' * 20)
                tasks.append(executor.submit(avurlFetch, *(keyword, engine, proxy)))
            concurrent.futures.wait(tasks)
            for task in tasks:
                avpages.extend(task.result())

        with concurrent.futures.ThreadPoolExecutor(max_workers=mthread) as executor:
            tasks = []
            for avpage in avpages:
                if dbfile is None or not os.path.exists(dbfile) or not duplicateCheck(avpage['code'].upper(), dbfile):
                    print((' -- Fetching ' + avpage['code'] + ' -- ').center(100, '#'))
                    tasks.append(executor.submit(avinfoFetch, *(avpage['url'], engine, proxy)))
            concurrent.futures.wait(tasks)
            for task in tasks:
                avs.append(task.result())
        print((' [ Saving Information to ' + tpath + ' ] ').center(100, '/'))
        avsave(avs, stype, tpath)
    except Exception as ex:
        logging.error("avfullFetch:" + str(ex))


# def avinfoFetch(url, engine='javbus', proxy=''):
def avinfoFetch(url, engine=InfoEngine, proxy=''):
    cav = None
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

    if engine == 'javbus' or 'javbus' in url:
        try:
            data = PyQuery(getHTML(url, 5, 5, 1, proxy))
            content = data('div.container')
            mosaic = str(data('ul[class="nav navbar-nav"]')('li[class="active"]').text()).strip()
            mosaic = mosaic.replace('無', '无').replace('碼', '码').replace('修正', '码')
            if mosaic == '码':
                mosaic = '有码'
            title = str(content('h3').eq(0).text()).strip()
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
                if '識別碼' in item.text() or '品番' in item.text() or '番號' in item.text()or '识别码' in item.text()or '番号' in item.text():
                    code = str(item('span:eq(1)').text()).strip()
                if '發行日' in item.text() or '発売日' in item.text() or '发行日' in item.text():
                    issuedate = str(re.search(r'\d*-\d*-\d*', item.text()).group()).strip()
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
            title = title.replace(code, '').strip()
            coverlink = str(content('a.bigImage').attr('href')).strip()
            cover = getHTML(coverlink, 5, 5, 0, proxy)
            try:
                link = avlinkFilter(avlinkFetch(code, 'zhongzilou', proxy)).link
            except Exception as ex:
                logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                link = 'page:' + url
            print(url.center(100, '-'))
            print('番号:'.center(5) + code + '\n' + '标题:'.center(5) + title + '\n' + '日期:'.center(5) + issuedate + '\n' + '时长:'.center(5) + length + '\n' + '修正:'.center(5) + mosaic + '\n' + '导演:'.center(5) + director + '\n' + '制作:'.center(5) + manufacturer + '\n' + '发行:'.center(5) + publisher + '\n' + '系列:'.center(5) + series + '\n' + '类别:'.center(5) + category + '\n' + '女优:'.center(5) + actors + '\n' + '收藏:'.center(5) + favor + '\n' + '预览:'.center(5) + coverlink + '\n' + '磁链:'.center(5) + link)

            cav = av(code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link)
        except Exception as ex:
            logging.warning('avinfoFetch:javbus:' + str(ex))
    if engine == 'javhoo' or 'javhoo' in url:
        try:
            data = PyQuery(getHTML(url, 5, 5, 1, proxy))
            content = data('div#content')('div.wf-container')
            avinfo = content('div.project_info')
            mosaic = str(avinfo('span.category-link').text()).strip()
            mosaic = mosaic.replace('無', '无').replace('碼', '码').replace('修正', '码')
            if mosaic == '码':
                mosaic = '有码'
            title = str(data('h1[class="h3-size entry-title"]').text()).strip()
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
                if '識別碼' in item.text() or '品番' in item.text() or '番號' in item.text()or '识别码' in item.text()or '番号' in item.text():
                    code = str(item('span:eq(1)').text()).split('/')[0].strip()
                if '發行日' in item.text() or '発売日' in item.text() or '发行日' in item.text():
                    issuedate = str(re.search(r'\d*-\d*-\d*', item.text()).group()).strip()
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
            title = title.replace(code, '').strip()
            coverlink = str(content('div.project-content')('img[class="alignnone size-full"]').attr('src')).strip()
            cover = getHTML(coverlink, 5, 5, 0, proxy)
            try:
                link = avlinkFilter(avlinkFetch(code, 'zhongzilou', proxy)).link
            except Exception as ex:
                logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                link = 'page:' + url
            print(url.center(100, '-'))
            print('番号:'.center(5) + code + '\n' + '标题:'.center(5) + title + '\n' + '日期:'.center(5) + issuedate + '\n' + '时长:'.center(5) + length + '\n' + '修正:'.center(5) + mosaic + '\n' + '导演:'.center(5) + director + '\n' + '制作:'.center(5) + manufacturer + '\n' + '发行:'.center(5) + publisher + '\n' + '系列:'.center(5) + series + '\n' + '类别:'.center(5) + category + '\n' + '女优:'.center(5) + actors + '\n' + '收藏:'.center(5) + favor + '\n' + '预览:'.center(5) + coverlink + '\n' + '磁链:'.center(5) + link)

            cav = av(code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link)
        except Exception as ex:
            logging.warning('avinfoFetch:javhoo:' + str(ex))
    if engine == 'torrentant' or 'torrentant' in url:
        try:
            data = PyQuery(getHTML(url, 5, 5, 1, proxy))
            content = data('div.container-fluid')('div.movie-view:first')
            avinfo = content('table.movie-view-table')
            mosaic = '未知'
            title = str(content('h1:first').text()).strip()
            code = str(avinfo('tr:eq(0)')('td:eq(1)').text()).strip()
            issuedate = str(avinfo('tr:eq(1)')('td:eq(1)').text()).strip()
            length = str(avinfo('tr:eq(2)')('td:eq(1)').text()).replace('分钟', '').replace('分鐘', '').replace('分', '').strip()
            director = str(avinfo('tr:eq(5)')('td:eq(1)').text()).strip()
            manufacturer = str(avinfo('tr:eq(3)')('td:eq(1)').text()).strip()
            publisher = str(avinfo('tr:eq(4)')('td:eq(1)').text()).strip()
            series = str(avinfo('tr:eq(6)')('td:eq(1)').text()).strip()
            category = str(content('div[class="col-md-12 tags"]').text()).strip()
            actors = str(content('div#avatar-waterfall').text()).strip()
            favor = '0'
            title = title.replace(code, '').strip()
            coverlink = ''
            cover = b''
            try:
                link = avlinkFilter(avlinkFetch(code, 'zhongzilou', proxy)).link
            except Exception as ex:
                logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                link = 'page:' + url
            print(url.center(100, '-'))
            print('番号:'.center(5) + code + '\n' + '标题:'.center(5) + title + '\n' + '日期:'.center(5) + issuedate + '\n' + '时长:'.center(5) + length + '\n' + '修正:'.center(5) + mosaic + '\n' + '导演:'.center(5) + director + '\n' + '制作:'.center(5) + manufacturer + '\n' + '发行:'.center(5) + publisher + '\n' + '系列:'.center(5) + series + '\n' + '类别:'.center(5) + category + '\n' + '女优:'.center(5) + actors + '\n' + '收藏:'.center(5) + favor + '\n' + '预览:'.center(5) + coverlink + '\n' + '磁链:'.center(5) + link)

            cav = av(code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link)
        except Exception as ex:
            logging.warning('avinfoFetch:torrentant:' + str(ex))
    return cav


def avlinkFetch(code, engine=BTEngine, proxy=''):
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
                    continue
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
                    continue
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
                    continue
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
                    time = str(linkdata('table[class="table table-hover"]')('tbody')('tr').eq(0)('td').eq(0).text()).strip()
                    clink = str(linkdata('a[class="btn btn-warning"]').attr('href')).strip()
                    avlinks.append(avlink(code, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('avlinkFetch:torrentant:' + str(ex))
                    continue
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
                    continue
        if engine == 'zhongzilou':
            data = PyQuery(getHTML('https://www.zhongzilou.com/list/' + code + '/1', 5, 5, 1, proxy))
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
                    logging.debug('avlinkFetch:zhongzilou:' + str(ex))
                    continue
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


def updateItemLink(code, dbfile):
    try:
        print(str('Updating ' + code).center(100, '*'))
        link = avlinkFilter(avlinkFetch(code, 'btdb', '')).link
        print(code.rjust(6) + ' <--> ' + link)
        isql = 'UPDATE av set link=? where code=?'
        update(get_conn(dbfile), isql, [(link, code)])
    except Exception as ex:
        logging.warning('updateItemLink:' + str(ex))


def avlinkUpdate(dbfile):
    if os.path.exists(dbfile):
        sql = 'SELECT code, link FROM av where link like "page:%"'
        res = fetchall(get_conn(dbfile), sql)
        if res is not None and len(res) > 0:
            print(str(len(res)).center(100, '+'))
            with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
                tasks = []
                for r in res:
                    tasks.append(executor.submit(updateItemLink, *(r[0], dbfile)))
                concurrent.futures.wait(tasks)


def clipthFetch(regstr, source, stype, tpath, mthread, engine, proxy, dbfile=None):
    tasks = []

    def consumer(regstr, stype, tpath, mthread, engine, proxy, dbfile):
        r = ''
        while True:
            if len(tasks) > 0:
                pattern = re.compile(regstr)
                keywords = list(set(number.group() for number in pattern.finditer(tasks.pop(0).strip())))
                avfullFetch(keywords, stype, tpath, mthread, engine, proxy, dbfile)

    def produce(source, proxy):
        preclip = ''
        while True:
            clipdata = pyperclip.paste()
            if clipdata is not None and clipdata != preclip and clipdata.strip() != '':
                preclip = clipdata
                if source == 'url':
                    hdata = getHTML(clipdata, 5, 3, 0, proxy)
                    if hdata is not None and hdata.strip() != '':
                        data = html2text.html2text(PyQuery(hdata).html())
                    else:
                        continue
                else:
                    data = clipdata
                tasks.append(data)

    # threading.Thread(target=produce).start()
    threading.Thread(target=consumer, args=(regstr, stype, tpath, mthread, engine, proxy, dbfile)).start()
    produce(source, proxy)


def clipFetch(regstr, source, mode, stype, tpath, mthread, engine, proxy, dbfile=None):
    # regstr = r'[A-Za-z]{1,7}-?[A-Za-z]?\d{2,4}-?\d{0,3}|\d{6}[-_]\d{4}[-_]\d{2}|\d{6}[-_]\d{2,3}|\d{6}-[A-Za-z]{3,6}|[A-Za-z]{1,3}\d[A-Za-z]{1,3}-\d{2,4}'
    # regstr = r'\S+'
    # regstr = r'[^\n]+'
    preclip = ''
    while True:
        try:
            keywords = []
            clipdata = pyperclip.paste()
            if source == 'url':
                clipdata = html2text.html2text(PyQuery(getHTML(clipdata, 5, 3, 0, proxy)).html())
            if clipdata is not None and clipdata != preclip and clipdata.strip() != '':
                if source == 'url':
                    hdata = getHTML(clipdata, 5, 3, 0, proxy)
                    if hdata is not None and hdata.strip() != '':
                        data = html2text.html2text(PyQuery(hdata).html())
                    else:
                        continue
                else:
                    data = clipdata
                preclip = clipdata
                pattern = re.compile(regstr)
                keywords = list(set(number.group() for number in pattern.finditer(data.strip())))
                avfullFetch(keywords, stype, tpath, mthread, engine, proxy, dbfile)
            if mode != 'loop':
                break
        except Exception as ex:
            logging.error('clipFetch:' + str(ex))
            continue


def cliploopFetch(mode=None):
    if mode is None or mode.strip() == '':
        print('Please Choose One ClipFetch Mode:')
        print('A. 按番号提取(剪切板)')
        print('B. 按空格分割(剪切板)')
        print('C. 按换行分割(剪切板)')
        print('D. 按番号提取(网址正文)')
        print('E. 按空格分割(网址正文)')
        print('F. 按换行分割(网址正文)')
        mode = input('My Choice: ')

    if mode.strip().lower() == 'a' or mode.strip() == '1' or mode.strip().lower() == 'd' or mode.strip() == '4' or mode.strip() == 'code':
        regstr = r'[A-Za-z]{1,7}-?[A-Za-z]?\d{2,4}-?\d{0,3}|\d{6}[-_]\d{4}[-_]\d{2}|\d{6}[-_]\d{2,3}|\d{6}-[A-Za-z]{3,6}|[A-Za-z]{1,3}\d[A-Za-z]{1,3}-\d{2,4}'
    elif mode.strip().lower() == 'b' or mode.strip() == '2' or mode.strip().lower() == 'e' or mode.strip() == '5' or mode.strip() == 'blank':
        regstr = r'\S+'
    elif mode.strip().lower() == 'c' or mode.strip() == '3' or mode.strip().lower() == 'f' or mode.strip() == '6' or mode.strip() == 'wrap':
        regstr = r'[^\n]+'
    else:
        regstr = mode.strip()
    if mode.strip().lower() in 'abc' or mode.strip().lower() in '123':
        clipthFetch(regstr, 'clip', 'db', 'C:/Users/xshrim/Desktop/imgss', 20, InfoEngine, '', 'C:/Users/xshrim/Desktop/imgss/avinfos.db')
        # clipthFetch(regstr, 'clip', 'file', 'C:/Users/xshrim/Desktop/imgss', 20, InfoEngine, '')
    else:
        clipthFetch(regstr, 'url', 'db', 'C:/Users/xshrim/Desktop/imgss', 20, InfoEngine, '', 'C:/Users/xshrim/Desktop/imgss/avinfos.db')
        # clipthFetch(regstr, 'url', 'file', 'C:/Users/xshrim/Desktop/imgss', 20, InfoEngine, '')


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
                continue
        try:
            link = avlinkFilter(links).link
        except Exception as ex:
            logging.debug('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
            link = 'page:' + 'https://www.javhoo.com/av/' + code
        return av(code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link)
    except Exception as ex:
        logging.error('avquickFetch:' + str(ex))
        return None


def keywordlinkFetch(keyword, engine=BTEngine, limit=-1, proxy='', sprint=True):
    head = ''
    time = ''
    hot = ''
    size = ''
    clink = ''
    klinks = []
    try:
        if engine == 'btgongchang':
            pidx = 1
            baseurl = 'http://btgongchang.org/search/' + keyword + '-first-asc-'
            while True:
                try:
                    if limit != -1 and int(limit) <= len(klinks):
                        break
                    curl = baseurl + str(pidx)
                    if sprint:
                        print(('Parsing Page (' + curl + ')').center(100, '-'))
                    data = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                    content = data('table[class="data mb20"]')
                    if '暂无' in content.text() and '相关资源' in content.text():
                        break
                    items = content('tr:gt(0)')
                    for item in items.items():
                        try:
                            if limit != -1 and int(limit) <= len(klinks):
                                break
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
                            al = avlink(keyword, head, time, hot, size, clink, engine)
                            klinks.append(al)
                            if sprint:
                                print(al.title.center(100, '*'))
                                print(al)
                        except Exception as ex:
                            logging.debug('keywordlinkFetch:btgongchang:' + str(ex))
                            continue
                    pidx += 1
                except Exception as ex:
                    logging.warning("keywordlinkFetch:btgongchang:" + str(ex))
                    break
        if engine == 'btso':
            pidx = 1
            baseurl = 'https://btso.pw/search/' + keyword + '/page/'
            while True:
                try:
                    if limit != -1 and int(limit) <= len(klinks):
                        break
                    curl = baseurl + str(pidx)
                    if sprint:
                        print(('Parsing Page (' + curl + ')').center(100, '-'))
                    data = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                    content = data('div.data-list')
                    items = content('div[class="row"]')
                    for item in items.items():
                        try:
                            if limit != -1 and int(limit) <= len(klinks):
                                break
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
                            al = avlink(keyword, head, time, hot, size, clink, engine)
                            klinks.append(al)
                            if sprint:
                                print(al.title.center(100, '*'))
                                print(al)
                        except Exception as ex:
                            logging.debug('keywordlinkFetch:btso:' + str(ex))
                            continue

                    pagination = data('ul[class="pagination pagination-lg"]')
                    if pagination is not None and str(pagination).strip() != '':
                        pginfo = ''
                        for item in pagination('li').items():
                            if str(pidx) == item.text().strip():
                                pginfo += '[' + str(pidx) + '] '
                            else:
                                if '首页' not in item.text() and '尾页' not in item.text() and '首頁' not in item.text() and '尾頁' not in item.text():
                                    pginfo += item.text().strip() + ' '
                        if sprint:
                            print(pginfo)
                        if '下一页' in pagination.text() or '下一頁' in pagination.text() or 'next' in pagination.text().lower():
                            pidx += 1
                        else:
                            break
                    else:
                        break
                except Exception as ex:
                    logging.warning("keywordlinkFetch:btso:" + str(ex))
                    break
        if engine == 'btdb':
            data = PyQuery(getHTML('https://btdb.in/q/' + keyword + '/?sort=popular', 5, 5, 1, proxy))
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
                    klinks.append(avlink(keyword, head, time, hot, size, clink, engine))
                except Exception as ex:
                    logging.debug('keywordlinkFetch:btdb:' + str(ex))
                    continue
        if engine == 'torrentant':
            pidx = 1
            baseurl = 'http://www.torrentant.com/cn/s/' + keyword + '?p='
            while True:
                try:
                    if limit != -1 and int(limit) <= len(klinks):
                        break
                    curl = baseurl + str(pidx)
                    if sprint:
                        print(('Parsing Page (' + curl + ')').center(100, '-'))
                    data = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                    content = data('ul[class="search-container"]')
                    items = content('li[class="search-item clearfix"]')
                    for item in items.items():
                        try:
                            if limit != -1 and int(limit) <= len(klinks):
                                break
                            head = str(item('div[class="search-content text-left"]')('h2')('a').attr('title')).strip()
                            head = re.sub(r'\/\*.*\*\/', '', head)
                            linkinfo = item('div[class="search-content text-left"]')('div[class="resultsContent"]')('p[class="resultsIntroduction"]')
                            size = str(linkinfo('label').eq(1).text()).strip().lower()
                            hot = str(linkinfo('label').eq(2).text()).strip()
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
                            time = str(linkdata('table[class="table table-hover"]')('tbody')('tr').eq(0)('td').eq(0).text()).strip()
                            clink = str(linkdata('a[class="btn btn-warning"]').attr('href')).strip()
                            al = avlink(keyword, head, time, hot, size, clink, engine)
                            klinks.append(al)
                            if sprint:
                                print(al.title.center(100, '*'))
                                print(al)
                        except Exception as ex:
                            logging.debug('keywordlinkFetch:torrentant:' + str(ex))
                            continue
                    pagination = data('ul[class="pagination"]')
                    if pagination is not None and str(pagination).strip() != '':
                        pginfo = ''
                        for item in pagination('li').items():
                            if str(pidx) == item.text().strip():
                                pginfo += '[' + str(pidx) + '] '
                            else:
                                if '首页' not in item.text() and '尾页' not in item.text() and '首頁' not in item.text() and '尾頁' not in item.text():
                                    pginfo += item.text().strip() + ' '
                        if sprint:
                            print(pginfo)
                        if '下一页' in pagination.text() or '下一頁' in pagination.text() or 'next' in pagination.text().lower():
                            pidx += 1
                        else:
                            break
                    else:
                        break
                except Exception as ex:
                    logging.warning("keywordlinkFetch:torrentant:" + str(ex))
                    break
        if engine == 'zhongzilou':
            pidx = 1
            baseurl = 'https://www.zhongzilou.com/list/' + keyword + '/'
            while True:
                try:
                    if limit != -1 and int(limit) <= len(klinks):
                        break
                    curl = baseurl + str(pidx)
                    if sprint:
                        print(('Parsing Page (' + curl + ')').center(100, '-'))
                    data = PyQuery(getHTML(curl, 5, 5, 1, proxy))
                    content = data('div.inerTop')
                    items = content('table[class="table table-bordered table-striped"]')
                    for item in items.items():
                        try:
                            if limit != -1 and int(limit) <= len(klinks):
                                break
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
                            al = avlink(keyword, head, time, hot, size, clink, engine)
                            klinks.append(al)
                            if sprint:
                                print(al.title.center(100, '*'))
                                print(al)
                        except Exception as ex:
                            logging.debug('keywordlinkFetch:zhongzilou:' + str(ex))
                            continue
                    pagination = data('ul[class="pagination"]')
                    if pagination is not None and str(pagination).strip() != '':
                        pginfo = ''
                        for item in pagination('li').items():
                            if str(pidx) == item.text().strip():
                                pginfo += '[' + str(pidx) + '] '
                            else:
                                if '首页' not in item.text() and '尾页' not in item.text() and '首頁' not in item.text() and '尾頁' not in item.text():
                                    pginfo += item.text().strip() + ' '
                        if sprint:
                            print(pginfo)
                        if '下一页' in pagination.text() or '下一頁' in pagination.text() or 'next' in pagination.text().lower():
                            pidx += 1
                        else:
                            break
                    else:
                        break
                except Exception as ex:
                    logging.warning("keywordlinkFetch:zhongzilou:" + str(ex))
                    break
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
        logging.warning('keywordlinkFetch:' + str(ex))
    return klinks


def av2file(avs, dirpath):
    txtfs = None
    txtname = 'avinfos.txt'
    count = 0
    try:
        print('Saving {0} AV Infomation to Files'.format(len(avs)).center(100, '*'))
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
                imgname = cav.code + ' ' + cav.title + '.' + ext
                imgname = imgname.replace('<', '').replace('>', '').replace('/', '').replace('\\', '').replace('|', '').replace(':', '').replace('"', '').replace('*', '').replace('?', '')
                imgpath = os.path.join(dirpath, imgname)
                imgfs = open(imgpath, 'wb')
                # imgfs.write(getHTML(cav.coverlink, 5, 5, 0, proxy))
                imgfs.write(cav.cover)
                imgfs.close()
                count += 1
                print('READY')
            except Exception as ex:
                logging.debug('av2file:' + str(ex))
                print('FAILED')
                continue
        txtfs.close()
        print(str(count) + '/' + str(len(avs)) + ' COMPLETE')
    except Exception as ex:
        if txtfs is not None:
            txtfs.close()
        logging.error('av2file:' + str(ex))
        print('FAILED')


def av2db(avs, dbfile):
    count = 0
    try:
        print('Saving {0} AV Infomation to Database'.format(len(avs)).center(100, '*'))
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
        # conn = get_conn(dbfile)
        create_table(get_conn(dbfile), sql)

        sql = '''INSERT OR IGNORE INTO av values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'''
        for cav in avs:
            try:
                print('Creating AV Information : ' + cav.title, end=' ...... ')
                save(get_conn(dbfile), sql, [(cav.code, cav.title, cav.issuedate, cav.length, cav.mosaic, cav.director, cav.manufacturer, cav.publisher, cav.series, cav.category, cav.actors, cav.favor, cav.coverlink, sqlite3.Binary(cav.cover), cav.link)])
                count += 1
                print('READY')
            except Exception as ex:
                logging.debug('av2db:' + str(ex))
                print('FAILED')
                continue
        print(str(count) + '/' + str(len(avs)) + ' COMPLETE')
    except Exception as ex:
        if 'UNIQUE constraint failed' in str(ex):
            logging.debug('av2db:' + str(ex))
            print(str(count) + '/' + str(len(avs)) + ' COMPLETE')
        else:
            logging.error('av2db:' + str(ex))
            print('FAILED')


def actres2db(actress, dbfile):
    count = 0
    try:
        print('Saving {0} Actres Infomation to Database'.format(len(actress)).center(100, '*'))
        sql = '''CREATE TABLE  IF NOT EXISTS `actres` (
              `jname` varchar(50) NOT NULL,
              `cname` varchar(50) DEFAULT NULL,
              `pname` varchar(50) DEFAULT NULL,
              `aname` varchar(50) DEFAULT NULL,
              `birthday` varchar(50) DEFAULT NULL,
              `birthplace` varchar(50) DEFAULT NULL,
              `cup` varchar(50) DEFAULT NULL,
              `stime` varchar(50) DEFAULT NULL,
              `status` varchar(50) DEFAULT NULL,
              `favorate` varchar(500) DEFAULT NULL,
              `codes` varchar(10000) DEFAULT NULL,
              `ipage` varchar(100) DEFAULT NULL,
              `notes` varchar(10000) DEFAULT NULL,
               PRIMARY KEY (`jname`)
            )'''

        create_table(get_conn(dbfile), sql)

        sql = '''INSERT OR IGNORE INTO actres values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'''
        for actres in actress:
            try:
                print('Creating Actres Information : ' + actres.jname, end=' ...... ')
                save(get_conn(dbfile), sql, [(actres.jname, actres.cname, actres.pname, actres.aname, actres.birthday, actres.birthplace, actres.cup, actres.stime, actres.status, actres.favorate, actres.codes, actres.ipage, actres.notes)])
                count += 1
                print('READY')
            except Exception as ex:
                logging.debug('actres2db:' + str(ex))
                print('FAILED')
                continue
        print(str(count) + '/' + str(len(actress)) + ' COMPLETE')
    except Exception as ex:
        if 'UNIQUE constraint failed' in str(ex):
            logging.debug('actres2db:' + str(ex))
            print(str(count) + '/' + str(len(actress)) + ' COMPLETE')
        else:
            logging.error('actres2db:' + str(ex))
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
                dbfile = os.path.join(dirpath, 'avinfos.db')
                av2db(avs, dbfile)
            elif savetype.lower() == 'both':
                dbfile = os.path.join(dirpath, 'avinfos.db')
                av2file(avs, dirpath)
                av2db(avs, dbfile)

        except Exception as ex:
            logging.error('avsave:' + str(ex))
    else:
        logging.error('No AV Infomation')


def main(argv):
    texts = []
    stype = 'file'
    surl = ''
    sfile = ''
    tpath = curDir()
    sengine = 'javbus'
    sproxy = ''
    smthread = 0
    autoloop = False
    keyword = ''
    count = -1
    keywords = []
    textwords = []
    filewords = []
    urlwords = []
    logInit()

    if argv is not None and len(argv) > 0:
        try:
            opts, args = getopt.getopt(argv, "had:e:t:p:m:u:f:s:k:c:", ["dir=", "engine=", "type=", "proxy=", "mthread=", "url=", "file=", "code=", "keyword=", "count="])
        except getopt.GetoptError:
            print(
                '''Usage: avfetch.py [-a] [-d <targetpath>] [-e <engine>] [-t <savetype>] [-p <proxy>] [-m <mthread>] [-u <url>] [-f <filename>] [-s <codes>] [<codes>] [<keyword>] [<count>]\n
                Example: avfetch.py -d D:/ -e javbus -t file -p socks5@127.0.0.1:1080 -m 5 -u http://www.baidu.com -f a.txt -s ABP-563 SRS-064 SNIS-862 -k 天海つばさ -c -1'''
            )
            sys.exit(2)

        if len(args) > 0:
            texts.extend(args)
        for opt, arg in opts:
            if opt == '-h':
                print(
                    '''Usage: avfetch.py [-a] [-d <targetpath>] [-e <engine>] [-t <savetype>] [-p <proxy>] [-m <mthread>] [-u <url>] [-f <filename>] [-s <codes>] [<codes>]\n
                    Example: avfetch.py -d D:/ -e javbus -t file -p socks5@127.0.0.1:1080 -m 5 -u http://www.baidu.com -f a.txt -s ABP-563 SRS-064 SNIS-862'''
                )
                exit()
            if opt == '-a':
                autoloop = True
                break
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
                    sys.exit(2)
            elif opt in ("-m", "--mthread"):
                smthread = int(arg)
            elif opt in ("-u", "--url"):
                surl = arg
            elif opt in ("-f", "--file"):
                sfile = os.path.join(curDir(), arg)
                if not os.path.isfile(sfile):
                    print('file is not exist!')
                    sys.exit(2)
            elif opt in ("-s", "--code"):
                texts.append(arg)
            elif opt in ("-k", "--keyword"):
                keyword = arg
            elif opt in ("-c", "--count"):
                count = arg
            else:
                pass
        try:
            if autoloop:
                cliploopFetch()
            elif keyword != '':
                if sengine == '' or sengine == 'javbus':
                    sengine = BTEngine
                with open(os.path.join(tpath, keyword + '.txt'), 'w') as wf:
                    for item in keywordlinkFetch(keyword, sengine, count, sproxy):
                        wf.write(item.link + '\n')
            else:
                if len(texts) > 0:
                    textwords = avkeywordParse(' '.join(texts), 'code')
                if sfile != '':
                    filewords = avkeywordParse(sfile, 'file')
                if surl != '':
                    urlwords = avkeywordParse(surl, 'url')
                keywords.extend(textwords)
                keywords.extend(filewords)
                keywords.extend(urlwords)
                # avfullFetch(keywords, stype, tpath, smthread, sengine, sproxy)
                avfullFetch(keywords, stype, tpath, smthread, sengine, sproxy)
        except Exception as ex:
            logging.error('main:' + str(ex))


if __name__ == "__main__":
    InfoEngine = choseEngine('ie')
    BTEngine = choseEngine('be')
    main(sys.argv[1:])

main(['-d', 'C:/Users/xshrim/Desktop/', '-e', 'javbus', '-p', 'socks5@127.0.0.1:12345', '-t', 'db', '-m', '20', '-s', ' 敗戦国の女'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-t', 'db', '-m', '20', '-s', ' 敗戦国の女'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'db', '-m', '20', '-s', 'ipz-137', 'ipz-371 midd-791 fset-337 sw-140'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javhoo', '-t', 'file', '-s', '天海つばさ'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgs', '-e', 'javbus', '-t', 'db', '-s', 'IPZ-137', 'IPZ820 MDS-825 FSET-337 F-123 FS-1'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'db', '-m', '20', '-f', 'C:/Users/xshrim/Desktop/av.txt'])
# main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'file', '-s', 'IPZ-137', 'IPZ820 MDS-825 FSET-337 F-123 FS-1'])

'''
InfoEngine = choseEngine('ie')
BTEngine = choseEngine('be')
for item in keywordlinkFetch('天海つばさ', BTEngine):
    print(item)
'''

'''
# 搜索引擎：
▶btso:https://btso.pw/search/ipz-137/
▶btsoso:http://www.sosobtt.com/s/ipz-371/
javhoo:https://www.javhoo.com/av/ipz-137/
▶btdb:https://btdb.in/q/ipz-137/
sukebei:https://sukebei.nyaa.se/?page=search&term=ipz-137&sort=4 (only torrent)
▶torrentant:http://www.torrentant.com/cn/s/ipz-137?sort=hot (inaccuracy)
▶btgongchang:http://btgongchang.org/search/MDS-825-first-asc-1
▶zhongzilou:http://www.zhongzilou.com/list/ipz-137/1
btago:http://www.btago.com/e/ipz-371/
btbook http://www.btwhat.org/
▶磁力搜 http://www.cilisou.cn/s.php?q=i~p~z~-~3~7~1~&encode_=1
▶bt蚂蚁 http://www.btanv.com/search/ipz-371-first-asc-1
▶屌丝懂 http://www.diaosisou.org/list/ipz-371/1
蜘蛛磁力 http://www.zzcili.org/so/ipz-371-first-asc-1?f=h
▶bt樱桃 http://www.btcerise.net/search?keyword=ipz-371
磁力吧 https://www.ciliba.org/s/ipz-371.html
bt磁力链 http://www.btcar.net/search/ipz-371_ctime_1.html
▶磁力链接 http://cililiana.com/list/ipz-371/1.html
▶BTKitty http://btkitty.pet/
▶btKitty http://www.bttit.com/torrent/ipz-371.html
▶磁力猪 http://www.cilizhuzhu.org/torrent/ipz-371.html
bt海 http://www.bthai.net/list/ipz-371-s1d-1.html
遨游搜 http://www.aoyoso.com/search/ipz-371_ctime_1.html
▶磁力屋 http://www.ciliwu.top/search-kw-ipz-371-p-1.html
▶bt兔子 http://www.btrabbit.cc/search/ipz-371.html
磁力之家 http://www.cilihome.net/word/ipz-371.html
bt快搜 http://www.btkuai.cc/word/ipz-371.html
▶阿狸磁力 http://alicili.net/list/ipz-371/1-0-0/
'''

# for cav in avlinkFetch('ipz-371', 'zhongzilou'):
#    print(cav)
# print(avlinkFilter(avlinkFetch('ipz-101', 'btso')).title)

'''
InfoEngine = choseEngine('ie')
BTEngine = choseEngine('be')
for item in avurlFetch('前田かおり', engine=InfoEngine, proxy=''):
    print(item)
'''

'''
for i in range(250, 300):
    url = 'http://hk-pic.xyz/forum-2-' + str(i) + '.html'
    print(url.center(100, '='))
    data = PyQuery(getHTML(url))
    content = html2text.html2text(data.html())
    regstr = r'[A-Za-z]{1,7}-?[A-Za-z]?\d{2,4}-?\d{0,3}|\d{6}[-_]\d{4}[-_]\d{2}|\d{6}[-_]\d{2,3}|\d{6}-[A-Za-z]{3,6}|[A-Za-z]{1,3}\d[A-Za-z]{1,3}-\d{2,4}'
    pattern = re.compile(regstr)
    keywords = list(set(number.group() for number in pattern.finditer(content.strip())))
    avfullFetch(keywords, 'db', 'C:/Users/xshrim/Desktop/imgss', 20, 'javbus', '', 'C:/Users/xshrim/Desktop/imgss/avinfos.db')
'''

# clipthFetch(regstr, 'db', 'C:/Users/xshrim/Desktop/imgss', 20, 'javbus', '', 'C:/Users/xshrim/Desktop/imgss/avinfos.db')

# cliploopFetch()
# links = keywordlinkFetch('ipz-371', 'btso')


'''
actress = []
with open('C:/Users/xshrim/Desktop/actres.txt', 'r', encoding='utf8') as rf:
    for line in rf.readlines():
        infos = line.split('*')
        cname = infos[0].strip()
        jname = infos[1].strip()
        pname = infos[2].strip()
        stime = infos[3].strip()
        notes = infos[4].strip()
        if '引退' in notes:
            status = '引退'
        else:
            status = '现役'
        actress.append(actres(jname, cname, pname, '', '', '', '', stime, status, '', '', '', notes))
        print(actress[-1])

with open('C:/Users/xshrim/Desktop/jb.txt', 'r', encoding='utf8') as rf:
    for line in rf.readlines():
        names = line.split('*')[0].strip()
        if '（' in names:
            jname = names.split('（')[0].strip()
            aname = names.split('（')[1].replace('）', '').strip()
        else:
            jname = names
            aname = ''
        ipage = line.split('*')[1]

        idx = -1
        for i in range(0, len(actress)):
            if actress[i].jname == jname:
                idx = i
                break
        if idx == -1:
            actress.append(actres(jname, '', '', aname, '', '', '', '', '', '', '', ipage, ''))
        else:
            actress[idx].aname = ''
            actress[idx].ipage = ipage
        print(actress[-1])

actres2db(actress, 'C:/Users/xshrim/Desktop/imgss/avinfos.db')
'''


'''
avactress = []
with open('C:/Users/xshrim/Desktop/b.txt', 'r', encoding='utf8') as rf:
    html = rf.read()
    data = PyQuery(html)
    items = data('table:gt(0)')('tr')
    for item in items.items():
        if item('td').length > 0:
            cname = item('td:eq(0)').text().strip().replace(' ', '')
            jnames = item('td:eq(1)').text().strip().replace(' ', '')
            if '（' in jnames:
                jname = jnames.split('（')[0].strip()
                janame = jnames.split('（')[1].replace('）', '').strip()
            else:
                jname = jnames.replace(' ', '')
                janame = ''
            stime = item('td:eq(2)').text().strip()
            if 'cup' in stime.lower():
                stime = item('td:eq(3)').text().strip()
                notes = item('td:eq(4)').text().strip()
            else:
                notes = item('td:eq(3)').text().strip()
            avactress.append((cname, jname, janame, stime, notes))
with open('C:/Users/xshrim/Desktop/actres.txt', 'w', encoding='utf8') as wf:
    for actres in avactress:
        wf.write(actres[0] + '*' + actres[1] + '*' + actres[2] + '*' + actres[3] + '*' + actres[4] + '\n')
'''


'''
codes = []
with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
    tasks = []
    for i in range(1, 40):
        url = 'http://www.javzoo.org/javmoo/index_' + str(i) + '.html'
        tasks.append(executor.submit(avkeywordParse, *(url, 'url')))
    concurrent.futures.wait(tasks)
    for task in tasks:
        codes.extend(task.result())
print(codes)
main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'db', '-m', '30', '-s', ' '.join(codes)])
'''


'''
urls = []
avs = []
data = PyQuery(getHTML('https://www.javbus.com/uncensored/genre'))
content = data('div[class="row genre-box"]')
avkinds = content('a')
for item in avkinds.items():
    urls.append(item.attr('href'))

with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
    tasks = []
    for url in urls:
        for avpage in avpageFetch(url, 'javbus', ''):
            tasks.append(executor.submit(avinfoFetch, *(avpage['url'], 'javbus', '')))
    concurrent.futures.wait(tasks)
    for task in tasks:
        avs.append(task.result())
avsave(avs, 'db', 'C:/Users/xshrim/Desktop/imgss/')
'''


'''
for i in range(1, 26):
    url = 'http://www.javlibrary.com/tw/vl_newrelease.php?list&mode=&page=1' + str(i)
    main(['-d', 'C:/Users/xshrim/Desktop/imgss', '-e', 'javbus', '-t', 'db', '-m', '10', '-u', url])
'''


'''
def tmpFetch(url, codes):
    print('Fetching ' + url)
    for avpage in avpageFetch(url, 'javbus', ''):
        codes.append(avpage['code'])


codes = []

threads = []
for i in range(1, 40):
    url = 'https://www.javbus.com/search/倶楽部/' + str(i)
    t = threading.Thread(target=tmpFetch, args=(url, codes))
    t.setDaemon(True)
    threads.append(t)
for t in threads:
    t.start()
for t in threads:
    t.join()

with open('avpage.txt', 'w') as f:
    for c in codes:
        f.write(c + '\n')
print(len(codes))
'''

'''
with open('C:/Users/xshrim/Desktop/av.txt', 'w') as f:
    for i in range(1, 10):
        url = 'https://www.javbus.com/page/' + str(i)
        for avpage in avpageFetch(url, 'javbus', ''):
            print(('Fetching Page ' + avpage['url']).center(100, '*'))
            f.write(avpage['code'] + '\n')
            # avs.append(avinfoFetch(avpage['url'], 'javbus', ''))
'''


'''
avs = []
with concurrent.futures.ThreadPoolExecutor(max_workers=20) as executor:
    tasks = []
    for i in range(1, 40):
        url = 'https://www.javbus.com/page/' + str(i)
        for avpage in avpageFetch(url, 'javbus', ''):
            tasks.append(executor.submit(avinfoFetch, *(avpage['url'], 'javbus', '')))
    concurrent.futures.wait(tasks)
    for task in tasks:
        avs.append(task.result())
avsave(avs, 'db', 'C:/Users/xshrim/Desktop/imgss/')
'''


'''
with open(r'D:/Git/minicode/avfetch/avs.txt', 'r', encoding='utf-8') as f:
    start = 6000
    keywords = []
    data = f.readlines()
    end = len(data)
    print(len(data))
    print('#' * 100)
    for i in range(start, end):
        code = str(data[i].split('#')[0]).strip()
        keywords.append(code)

    tasks = []
    pool = concurrent.futures.ThreadPoolExecutor(10)
    for i in range(0, len(keywords), 10):
        tasks.append(pool.submit(avfullFetch, *(keywords[i:i + 10], 'db', r'C:/Users/xshrim/Desktop/imgss', 'javbus', '')))
    concurrent.futures.wait(tasks)
    print('ALL OK')
'''
# avfullFetch(keywords, 'both', r'C:/Users/xshrim/Desktop/imgss', 'javbus', '')

# cliploopFetch()
