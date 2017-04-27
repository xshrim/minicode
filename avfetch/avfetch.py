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
from urllib import request
# from urllib import parse
from pyquery import PyQuery

# sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf8')


class av(object):
    '''
    code = ''
    title = ''
    issuedate = ''
    length = ''
    director = ''
    category = ''
    actors = ''
    coverlink = ''
    cover = b''
    links = ''
    '''

    def __init__(self, code, title, issuedate, length, director, category, actors, coverlink, cover, links):
        self.code = code
        self.title = title
        self.issuedate = issuedate
        self.length = length
        self.director = director
        self.category = category
        self.actors = actors
        self.coverlink = coverlink
        self.cover = cover
        self.links = links

    def __str__(self):
        '''
        print('番号:'.rjust(5) + self.code)
        print('标题:'.rjust(5) + self.title)
        print('发行:'.rjust(5) + self.issuedate)
        print('时长:'.rjust(5) + self.length)
        print('导演:'.rjust(5) + self.director)
        print('类别:'.rjust(5) + self.category)
        print('女优:'.rjust(5) + self.actors)
        print('预览:'.rjust(5) + self.coverlink)
        print('磁链:'.rjust(5) + self.links)
        '''
        return '番号:'.rjust(5) + self.code + '\n' + '标题:'.rjust(5) + self.title + '\n' + '发行:'.rjust(5) + self.issuedate + '\n' + '时长:'.rjust(5) + self.length + '\n' + '导演:'.rjust(5) + self.director + '\n' + '类别:'.rjust(5) + self.category + '\n' + '女优:'.rjust(5) + self.actors + '\n' + '预览:'.rjust(5) + self.coverlink + '\n' + '磁链:'.rjust(5) + self.links

    __repr__ = __str__


class avlinkinfo(object):
    '''
    code = ''
    title = ''
    time = ''
    hot = ''
    size = ''
    link = ''
    '''

    def __init__(self, code, title, time, hot, size, link):
        self.code = code
        self.title = title
        self.time = time
        self.hot = hot
        self.size = size
        self.link = link

    def __str__(self):
        return self.code + ' -- ' + self.title + ' -- ' + self.time + ' -- ' + self.hot + ' -- ' + self.size + ' -- ' + self.link

    __repr__ = __str__

######################################### DB START#########################################


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
        print('the [{}] is empty or equal None!'.format(sql))


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
        print('the [{}] is empty or equal None!'.format(sql))


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
        print('the [{}] is empty or equal None!'.format(sql))


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
        print('the [{}] is empty or equal None!'.format(sql))


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
            print('the [{}] equal None!'.format(data))
    else:
        print('the [{}] is empty or equal None!'.format(sql))


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
        print('the [{}] is empty or equal None!'.format(sql))


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
        print('the [{}] is empty or equal None!'.format(sql))

######################################### DB END#########################################


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
            except Exception as ex:
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
    headers = [('Host', 'img0.imgtn.bdimg.com'), ('Connection', 'keep-alive'), ('Cache-Control', 'max-age=0'), ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'), ('User-Agent', 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'), ('Accept-Encoding', '*'), ('Accept-Language', 'zh-CN,zh;q=0.8'), ('If-None-Match', '90101f995236651aa74454922de2ad74'), ('Referer', 'http://www.deviantart.com/whats-hot/'), ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

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
        pattern = re.compile(r'[A-Za-z]{2,5}-?\d{2,3}')
        if type == 'file':
            sfile = os.path.join(curDir(), keyword)
            chartype = charDetect(open(sfile, 'rb').read())
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
                p = re.match(r'([A-Za-z]+).*?(\d+)', code)
                code = str(p.group(1)) + '-' + str(p.group(2))
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
            coverlink = ''
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
                coverlink = content('a.bigImage').attr('href')
                cover = getHTML(coverlink, 5, 5, 0, proxy)
                # magnets = content('table#magnet-table')
                try:
                    links = avfilter(avlinkFetch(code, 'BT工厂', proxy)).link
                except Exception as ex:
                    print('#' * 32 + '  No magnet link!  Show info page.  ' + '#' * 32)
                    links = 'page:' + curl
            avs.append(av(code, title, issuedate, length, director, category, actors, coverlink, cover, links))
            print('#' * 100)
            print(avs[-1])
            # print(title + ' <-> ' + actors + ' <-> ' + coverlink + ' <-> ' + links)
        except Exception as ex:
            print('avinfoFetch:' + str(ex))
    return avs


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


def av2file(avs, dirpath):
    txtfs = None
    txtname = 'avinfos.txt'
    try:
        print('Saving AV Infomation to Files:')
        txtpath = os.path.join(dirpath, txtname)
        txtfs = open(txtpath, 'a', encoding='utf8')
        for cav in avs:
            print('*' * 100)
            print('Creating AV Information : ' + cav.title, end=' ...... ')
            txtfs.write('番号:'.rjust(5) + cav.code + '\n')
            txtfs.write('标题:'.rjust(5) + cav.title + '\n')
            txtfs.write('发行:'.rjust(5) + cav.issuedate + '\n')
            txtfs.write('时长:'.rjust(5) + cav.length + '\n')
            txtfs.write('导演:'.rjust(5) + cav.director + '\n')
            txtfs.write('类别:'.rjust(5) + cav.category + '\n')
            txtfs.write('女优:'.rjust(5) + cav.actors + '\n')
            txtfs.write('预览:'.rjust(5) + cav.coverlink + '\n')
            txtfs.write('磁链:'.rjust(5) + cav.links + '\n')
            txtfs.write('#' * 100 + '\n')
            ext = cav.coverlink.split('.')[-1] if '.' in cav.coverlink else 'jpg'
            imgname = cav.title + '.' + ext
            imgname = imgname.replace('<', '').replace('>', '').replace('/', '').replace('\\', '').replace('|', '').replace(':', '').replace('"', '').replace('*', '').replace('?', '')
            imgpath = os.path.join(dirpath, imgname)
            imgfs = open(imgpath, 'wb')
            # imgfs.write(getHTML(cav.coverlink, 5, 5, 0, proxy))
            imgfs.write(cav.cover)
            imgfs.close()
            print('READY')
        txtfs.close()
        print('COMPLETE')
    except Exception as ex:
        if txtfs is not None:
            txtfs.close()
        print('av2file:' + str(ex))
        print('FAILED')


def av2db(avs, dirpath):
    data = []
    dbname = 'avinfos.db'
    try:
        print('Saving AV Infomation to Database:')
        dbpath = os.path.join(dirpath, dbname)
        sql = '''CREATE TABLE  IF NOT EXISTS `av` (
              `code` varchar(20) NOT NULL,
              `title` varchar(300) NOT NULL,
              `issuedate` varchar(20) DEFAULT NULL,
              `length` varchar(20) DEFAULT NULL,
              `director` varchar(20) DEFAULT NULL,
              `category` varchar(300) DEFAULT NULL,
              `actors` varchar(300) DEFAULT NULL,
              `coverlink` varchar(300) DEFAULT NULL,
              `cover` BLOB DEFAULT NULL,
              `links` varchar(10000) DEFAULT NUll,
               PRIMARY KEY (`code`)
            )'''
        conn = get_conn(dbpath)
        create_table(conn, sql)

        sql = '''INSERT OR IGNORE INTO av values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'''
        for cav in avs:
            print('*' * 100)
            print('Creating AV Information : ' + cav.title, end=' ...... ')
            data.append((cav.code, cav.title, cav.issuedate, cav.length, cav.director, cav.category, cav.actors, cav.coverlink, sqlite3.Binary(cav.cover), cav.links))
            print('READY')
        save(conn, sql, data)
        print('COMPLETE')
    except Exception as ex:
        if 'UNIQUE constraint failed' in str(ex):
            print('av2db:' + str(ex))
            print('COMPLETE')
        else:
            print('av2db:' + str(ex))
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
            elif savetype.lower == 'both':
                av2file(avs, dirpath)
                av2db(avs, dirpath)

        except Exception as ex:
            print('avsave:' + str(ex))
    else:
        print('No AV Infomation')


def main(argv):
    codes = []
    stype = 'file'
    surl = ''
    sfile = ''
    tpath = ''
    sengine = ''
    sproxy = ''
    try:
        opts, args = getopt.getopt(argv, "hd:e:t:p:u:f:s:", ["dir=", "engine=", "type=", "proxy=", "url=", "file=", "code="])
    except getopt.GetoptError:
        print(
            '''Usage: avfetch.py [-d <targetpath>] [-e <engine>] [-t <savetype>] [-p <proxy>] [-u <url>] [-f <filename>] [-s <codes>] [<codes>]\n
            Example: avfetch.py -d D:/ -e javbus -t file -p socks5@127.0.0.1:1080 -u http://www.baidu.com -f a.txt -s ABP-563 SRS-064 SNIS-862'''
        )
        exit(2)

    if len(args) > 0:
        codes.extend(args)
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
        elif opt in ("-s", "--code"):
            codes.append(arg)
    try:
        if surl != '':
            codes.append(str(getHTML(surl, 5, 5, 0, proxy=sproxy)))
        if sfile != '' and os.path.isfile(sfile):
            chartype = charDetect(open(sfile, 'rb').read())
            for line in open(sfile, encoding=chartype):
                codes.append(line)
        if len(codes) > 0:
            avs = avinfoFetch(' '.join(codes), 'code', engine=sengine, proxy=sproxy)
            avsave(avs, stype, tpath)
    except Exception as ex:
        print('main:' + str(ex))


if __name__ == '__main__':
    main(sys.argv[1:])

# main(['-d', 'imgss', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-u', 'http://btgongchang.org/'])
main(['-d', 'C:/Users/xshrim/Desktop/imgs', '-e', 'javbus', '-t', 'db', '-s', 'IPZ-137', 'IPZ820 MDS-825 FSET-337 F-123 FS-1'])
# main(['-d', 'imgs', '-e', 'javbus', '-p', 'socks5@127.0.0.1:1080', '-f', 'a.txt'])
