import datetime
import getopt
import os
import random
import re
# import socket
import sqlite3
import sys
import time

# import chardet
import colorama
import pyperclip


# import webbrowser
# from urllib import parse, request
# from urllib.parse import quote

# import socks
# from pyquery import PyQuery


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

###############################################################
####            创建|删除表操作     START
###############################################################


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
    '''创建数据库表：student'''
    if sql is not None and sql != '':
        cu = get_cursor(conn)
        # print('执行sql:[{}]'.format(sql))
        cu.execute(sql)
        conn.commit()
        # print('创建数据库表成功!')
        close_all(conn, cu)
    else:
        print('the [{}] is empty or equal None!'.format(sql))


###############################################################
####            创建|删除表操作     END
###############################################################


def close_all(conn, cu):
    '''关闭数据库游标对象和数据库连接对象'''
    try:
        if cu is not None:
            cu.close()
    finally:
        if cu is not None:
            cu.close()

###############################################################
####            数据库操作CRUD     START
###############################################################


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
    res = []
    if sql is not None and sql != '':
        cu = get_cursor(conn)
        # print('执行sql:[{}]'.format(sql))
        cu.execute(sql)
        r = cu.fetchall()
        if len(r) > 0:
            for e in range(len(r)):
                res.append(r[e])
    else:
        print('the [{}] is empty or equal None!'.format(sql))
    return res


def fetchone(conn, sql, data):
    '''查询一条数据'''
    if sql is not None and sql != '':
        if data is not None:
            # Do this instead
            # d = (data,)
            cu = get_cursor(conn)
            # print('执行sql:[{}],参数:[{}]'.format(sql, data))
            cu.execute(sql, data)
            r = cu.fetchall()
            if len(r) > 0:
                for e in range(len(r)):
                    return r[e]
        else:
            print('the [{}] equal None!'.format(data))
    else:
        print('the [{}] is empty or equal None!'.format(sql))
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
###############################################################
####            数据库操作CRUD     END
###############################################################


def curDir():
    try:
        return os.path.split(os.path.realpath(__file__))[0]
    except Exception as ex:
        return os.getcwd()


'''
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
    headers = [('Host', 'img0.imgtn.bdimg.com'), ('Connection', 'close'), ('Cache-Control', 'max-age=0'), ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'), ('User-Agent', 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'), ('Accept-Encoding', '*'), ('Accept-Language', 'zh-CN,zh,en-US,en,*;q=0.8'), ('If-None-Match', '90101f995236651aa74454922de2ad74'), ('Referer', 'http://www.deviantart.com/whats-hot/'), ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

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
            print('getHTML:' + str(ex))
            if '403' in str(ex) or '404' in str(ex) or '11001'in str(ex):
                break
        i -= 1
    return contents
'''


def collect(level, dbfile=os.path.join(curDir(), 'orath.db')):
    '''
      `id` integer PRIMARY KEY autoincrement,   #自增主键
      `class` varchar(10) NOT NULL,             #类别(OCA, OCP, OCM)
      `level` varchar(10) NOT NULL,             #级别(051, 052, 053)
      `db` varchar(20) NOT NULL,                #数据库(oracle 11g, oracle 12c)
      `version` varchar(10) NOT NULL,           #版本(v8.02, v9.02)
      `qn` integer NOT NULL,                    #题号(1, 2, 3)
      `link` varchar(100),                      #链接
      `content` varchar(100000),                #题目内容
      `image` BLOB DEFAULT NULL,                #题目图片
      `options` varchar(100000),                #题目选项
      `parse` varchar(100000),                  #题目解析
      `reference` varchar(100000),              #参考内容
      `answer` varchar(100000),                 #答案
      `skill` integer DEFAULT 0,                #熟练度(答错次数)
      `star` integer DEFAULT 0,                 #星标
      `tmp1` varchar(100),                      #预留1
      `tmp2` varchar(100),                      #预留2
     '''

    create_table_sql = '''CREATE TABLE IF NOT EXISTS `orath` (
                          `id` integer PRIMARY KEY autoincrement,
                          `level` varchar(10) NOT NULL,
                          `db` varchar(20) NOT NULL,
                          `version` varchar(10) NOT NULL,
                          `qn` integer NOT NULL,
                          `link` varchar(100),
                          `content` varchar(100000),
                          `image` BLOB DEFAULT NULL,
                          `options` varchar(100000),
                          `parse` varchar(100000),
                          `reference` varchar(100000),
                          `answer` varchar(100000),
                          `skill` integer DEFAULT 0,
                          `star` integer DEFAULT 0,
                          `tmp1` varchar(100),
                          `tmp2` varchar(100)
                        )'''
    iconn = get_conn(dbfile)
    create_table(iconn, create_table_sql)

    idata = []
    isql = 'insert into `orath` values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'
    fetchTopic(level, idata)
    save(iconn, isql, idata)


'''
def fetchTopic(level, idata):
    if level == '1Z0-051':
        data = PyQuery(getHTML('http://blog.csdn.net/rlhua/article/details/17101765'))
        content = data('#article_content')
        table = content('table')
        items = table('tr:gt(0)')
        for item in items.items():
            qn = str(item('td:eq(0)').text()).strip()
            link = str(item('td:eq(1)').text()).strip()
            idata.append([None, '1Z0-051', 'Oracle 11g r2', 'v9.02', qn.strip(), link, None, None, None, None, None, None, None, None, None, None])

    if level == '1Z0-052':
        with open('C:/Users/xshrim/Desktop/1Z0-052V10.02.new.txt', 'r') as rf:
            tmpdata = []
            for line in rf:
                if line.strip() != '':
                    if 'Answer:' not in line:
                        tmpdata.append(line)
                    else:
                        answer = line.replace('Answer:', '').strip()
                        qn = tmpdata[0].split('.')[0]
                        if int(qn) < 213:
                            link = t052[qn]
                        else:
                            link = ''
                        idata.append([None, '1Z0-052', 'Oracle 11g r2', 'v10.02', qn.strip(), link, ''.join(tmpdata), None, None, None, None, answer, None, None, None, None])
                        tmpdata.clear()
    if level == '1Z0-053':
        with open('C:/Users/xshrim/Desktop/1Z0-053V14.02.new.txt', 'r') as rf:
            tmpdata = []
            for line in rf:
                if line.strip() != '':
                    if 'Answer:' not in line:
                        tmpdata.append(line)
                    else:
                        answer = line.replace('Answer:', '').strip()
                        qn = tmpdata[0].split('.')[0]
                        print(qn)
                        idata.append([None, '1Z0-053', 'Oracle 11g r2', 'v14.02', qn.strip(), None, ''.join(tmpdata), None, None, None, None, answer, None, None, None, None])
                        tmpdata.clear()
'''


def showTopic(level, qn, wb=False, dbfile=os.path.join(curDir(), 'orath.db')):
    iconn = get_conn(dbfile)
    isql = 'select * from `orath` where `level`=? and `qn`=?'
    res = fetchone(iconn, isql, (level, qn))
    if res is not None:
        level = res[1]
        db = res[2]
        version = res[3]
        qn = res[4]
        link = res[5]
        content = res[6]
        options = res[8]
        answer = res[11]
        star = res[13]
        print(colorama.Style.BRIGHT + colorama.Fore.GREEN + (str(db) + ' <--> ' + str(version) + ' <--> ' + str(level) + ' <--> ' + str(qn)).center(100, '='))
        print(colorama.Fore.CYAN + '题目:\n' + content)
        print(colorama.Fore.CYAN + '选项:\n' + options)
        print(colorama.Fore.MAGENTA + '答案:\n' + answer + colorama.Style.NORMAL + colorama.Fore.RESET)
        if wb:
            print(link)
            # webbrowser.open(link, new=0, autoraise=True)


'''
def updateTopic(level, qn, dbfile=os.path.join(curDir(), 'orath.db')):
    try:
        iconn = get_conn(dbfile)
        isql = 'select * from `orath` where `level`=? and `qn`=?'
        res = fetchone(iconn, isql, (level, qn))
        if res is not None:
            link = res[5]
            data = PyQuery(getHTML(link))
            content = data('#article_content')
            image = content('img').attr('src')
            timage = getHTML(image)
            if not isinstance(timage, bytes):
                timage = b''
            fulldata = ''
            tcontent = ''
            tanswer = ''
            tparse = ''
            for item in content.children()('div, span').items():
                fulldata += item.text() + '\n'
            if 'Answer:' in fulldata:
                tcontent = fulldata.split('Answer:')[0]
                tcont = fulldata.split('Answer:')[1]
                if '答案解析' in tcont:
                    tanswer = tcont.split('答案解析')[0]
                    tparse = '\n'.join(str(tcont.split('答案解析')[1]).split('\n')[1:])
                else:
                    tanswer = tcont.split('\n')[0]
                    tparse = '\n'.join(tcont.split('\n')[1:])
            else:
                if '此题答案' in fulldata:
                    for line in fulldata.split('\n'):
                        if '此题答案' in line:
                            tcontent = fulldata.replace(line, '')
                            tparse = ''
                            tanswer = line.replace('此题答案', '').replace('选', '').replace('为', '')
            print(tcontent)
            print(tanswer)
            print(tparse)
            isql = 'update `orath` set `content`=?, `image`=?, `parse`=?, `answer`=? where `level`=? and `qn`=?'
            res = update(iconn, isql, [(tcontent, sqlite3.Binary(timage), tparse, tanswer, level, qn)])
            # print(nltk.clean_html(content))
    except Exception as ex:
        print('starTopic:' + str(ex))
'''


def starTopic(level, qn, star, dbfile=os.path.join(curDir(), 'orath.db')):
    try:
        iconn = get_conn(dbfile)
        isql = 'update `orath` set `star`=? where `level`=? and `qn`=?'
        res = update(iconn, isql, [(star.upper().strip(), level, qn)])
    except Exception as ex:
        print('starTopic:' + str(ex))


def getAnswer(level, keywords, dbfile=os.path.join(curDir(), 'orath.db')):
    res = None
    try:
        iconn = get_conn(dbfile)
        predc = ''
        predo = ''
        pred = ''
        for keyword in keywords.split(' '):
            if keyword.strip() != '':
                predc += '`content` like "%' + keyword.strip() + '%" and '
        predc = predc[:-5]
        for keyword in keywords.split(' '):
            if keyword.strip() != '':
                predo += '`options` like "%' + keyword.strip() + '%" and '
        predo = predo[:-5]
        pred = predc + ' or ' + predo
        if level.strip() == '':
            isql = 'select level, qn, content, options, answer from `orath` where ' + pred + ' order by qn'
        else:
            spred = ''
            level = level.replace('，', ',')
            for sl in level.split(','):
                if sl.strip() != '':
                    spred += '`level` = "' + sl.strip() + '" or '
            spred = spred[:-4]
            isql = 'select level, qn, content, options, answer from `orath` where (' + pred + ') and (' + spred + ') order by qn'
            print(isql)
        res = fetchall(iconn, isql)
    except Exception as ex:
        print('getAnswer:' + str(ex))
    return res


def refer(level, dbfile=os.path.join(curDir(), 'orath.db')):
    preclip = ''
    mode = 'manual'
    print(colorama.Style.BRIGHT + colorama.Back.RESET + colorama.Fore.WHITE + '#' * 100)
    print('1. Manual')
    print('2. Automatic')
    mode = input('Please select refer mode:')
    if mode == '1' or mode.lower() == 'manual':
        mode = 'manual'
    if mode == '2' or mode.lower() == 'automatic':
        mode = 'automatic'
    while True:
        try:
            if mode == 'manual':
                keyword = input(colorama.Style.BRIGHT + colorama.Back.RESET + colorama.Fore.GREEN + 'Please input the keyword: ')
            else:
                keyword = pyperclip.paste()
            if keyword is not None and keyword.strip() != '' and keyword != preclip:
                preclip = keyword
                print(colorama.Style.BRIGHT + colorama.Back.RED + colorama.Fore.YELLOW + keyword.strip().center(100, '='))
                for r in getAnswer(level, keyword, dbfile):
                    print(colorama.Style.BRIGHT + colorama.Back.RESET + colorama.Fore.YELLOW + (str(r[0]) + ' -> ' + str(r[1])).center(20, ' ').center(100, '#'))
                    # print('题号：' + str(r[0]) + ' -> ' + str(r[1]))
                    print(colorama.Style.BRIGHT + colorama.Back.RESET + colorama.Fore.CYAN + '题目：\n' + str(r[2]).strip())
                    print(colorama.Style.BRIGHT + colorama.Back.RESET + colorama.Fore.CYAN + '选项：\n' + str(r[3]).strip())
                    print(colorama.Style.BRIGHT + colorama.Back.RESET + colorama.Fore.MAGENTA + '答案：\n' + str(r[4]).strip())

                    # os.system('cmd /k echo ' + output)
                    '''
                    print('*' * 100)
                    print('题号：' + str(r[0]) + ' -> ' + str(r[1]))
                    print('题目：' + str(r[2]))
                    print('答案：' + str(r[3]))
                    '''
        except Exception as ex:
            print('refer' + str(ex))


def practise(level=None, scope=None, count=None, mode=None, sr=None, star=None, ktime=None, dbfile=os.path.join(curDir(), 'orath.db')):
    corlist = []
    wrolist = []
    topics = []
    iconn = get_conn(dbfile)
    isql = 'select level, qn, content, options, answer, star from `orath`'
    res = fetchall(iconn, isql)
    for r in res:
        topics.append((r[0], r[1], r[2], r[3], r[4], r[5]))

    print(colorama.Style.BRIGHT + colorama.Fore.WHITE + '#' * 100)
    if level is None or str(level).strip() == '':
        print('1. 1Z0-051 (OCA)')
        print('2. 1Z0-052 (OCP)')
        print('3. 1Z0-053 (OCP)')
        level = input('Please select level:')
    if level == '1' or level == '051' or str(level).strip().upper() == '1Z0-051':
        level = '1Z0-051'
    if level == '2' or level == '052' or str(level).strip().upper() == '1Z0-052':
        level = '1Z0-052'
    if level == '3' or level == '053' or str(level).strip().upper() == '1Z0-053':
        level = '1Z0-053'
    stopics = [topic for topic in topics if topic[0] == level]
    print(colorama.Fore.GREEN + str(level) + ' Total: ' + str(len(stopics)))
    print(colorama.Fore.WHITE + '-' * 100)

    if scope is None or str(scope).strip() == '0':
        scope = input('Please input question scope(-1 for All):')
    if scope == -1 or '-' not in str(scope):
        print(colorama.Fore.GREEN + 'Question Scope: All')
    else:
        try:
            qstart, qend = str(scope).split('-')
            stopics = stopics[int(qstart) - 1:int(qend)]
            print(colorama.Fore.GREEN + 'Question Scope: ' + str(scope))
        except Exception as ex:
            print(str(ex))
            print(colorama.Fore.GREEN + 'Question Scope: All')
    print(colorama.Fore.WHITE + '-' * 100)

    if count is None or str(count).strip() == '0':
        count = input('Please input question quantity(-1 for unlimited):')
    count = int(count)
    if count == -1:
        print(colorama.Fore.GREEN + 'Question Quantity: Unlimited')
    else:
        print(colorama.Fore.GREEN + 'Question Quantity: ' + str(count))
    print(colorama.Fore.WHITE + '-' * 100)

    if mode is None or str(mode).strip() == '':
        mode = input('Please select question mode(random or sequence):')
    if str(mode).strip() == '1' or str(mode).strip().lower() == 'r' or str(mode).strip().lower() == 'y' or str(mode).strip().lower() == 'random':
        mode = 'random'
    else:
        mode = 'sequence'
    print(colorama.Fore.GREEN + 'Mode: ' + mode)
    print(colorama.Fore.WHITE + '-' * 100)

    if sr is None or str(sr).strip() == '':
        sr = input('Show the question answer(yes or no):')
    if str(sr).strip() == '1' or str(sr).strip().lower() == 'y' or str(sr).strip().lower() == 'yes' or str(sr).strip().lower() == '是':
        sr = 'yes'
    else:
        sr = 'no'
    print(colorama.Fore.GREEN + 'Show Answer: ' + sr)
    print(colorama.Fore.WHITE + '-' * 100)

    if star is None or str(star).strip() == '':
        star = input('Only show star questions(yes or no):')
    if str(star).strip() == '1' or str(star).strip().lower() == 'y' or str(star).strip().lower() == 'yes' or str(star).strip().lower() == '是':
        star = 'yes'
        stopics = [topic for topic in stopics if topic[5] == 'S']
    else:
        star = 'no'
    print(colorama.Fore.GREEN + 'Only Show Star Questions: ' + star)
    print(colorama.Fore.WHITE + '-' * 100)

    if ktime is None or str(ktime).strip() == '':
        ktime = input('Please input time limit in minutes(-1 for unlimited):')
    try:
        ktime = int(ktime)
    except Exception as ex:
        ktime = -1
    if ktime == -1:
        print(colorama.Fore.GREEN + 'Time Limit: Unlimited')
    else:
        print(colorama.Fore.GREEN + 'Time Limit: ' + str(ktime) + ' Minutes')
    print(colorama.Fore.WHITE + '#' * 40 + ('SUM: ' + str(len(stopics))).center(20) + '#' * 40)

    idx = 1
    starttime = datetime.datetime.now()
    while True:
        if count != -1 and idx > count:
            break
        if mode == 'random':
            ctopic = random.choice(stopics)
        else:
            if idx <= len(stopics):
                ctopic = stopics[idx - 1]
            else:
                break
        if int(ktime) == -1:
            currenttime = datetime.datetime.now()
            print(colorama.Fore.GREEN + (starttime.strftime('%H:%M:%S').center(30, '=') + ('[' + str(idx) + ']').center(8) + '=> ' + str(ctopic[0]) + ' <--> ' + str(ctopic[1]) + '  ' + currenttime.strftime('%H:%M:%S').center(30, '=')).center(100, '='))
        else:
            currenttime = datetime.datetime.now()
            ctime = int((currenttime - starttime).seconds)
            mdelta = (int(ktime) * 60 - ctime) // 60
            sdelta = (int(ktime) * 60 - ctime) % 60
            if int(ktime) * 60 - ctime <= 0:
                break
            else:
                tdelta = str(mdelta).rjust(2, '0') + ':' + str(sdelta).rjust(2, '0')
                print(colorama.Fore.GREEN + (starttime.strftime('%H:%M:%S').center(30, '=') + ('[' + str(idx) + ']').center(8) + '=> ' + str(ctopic[0]) + ' <--> ' + str(ctopic[1]) + '  ' + (currenttime.strftime('%H:%M:%S') + ' (' + tdelta + ')').center(40, '=')).center(100, '='))
        print(colorama.Fore.CYAN + ctopic[2] + '\n')
        print(colorama.Fore.CYAN + ctopic[3])
        ipt = input(colorama.Fore.YELLOW + 'Please input your answer (q for quit | s for star | u for unstar) => ')
        if ipt.strip().lower() == 'q' or ipt.strip().lower() == 'quit' or ipt.strip().lower() == 'exit':
            break
        else:
            if ipt.lower().strip() == 's' or ipt.lower().strip() == 'u':
                starTopic(ctopic[0], ctopic[1], ipt, dbfile)
                ipt = input(colorama.Fore.YELLOW + 'Please input your answer (q for quit) => ')
            if ipt.strip().lower() == 'q' or ipt.strip().lower() == 'quit' or ipt.strip().lower() == 'exit':
                break
            elif ipt.strip().upper().replace(' ', '') == ctopic[4].strip().upper().replace(' ', ''):
                corlist.append((idx, ctopic[0], ctopic[1]))
                if sr == 'yes':
                    print(colorama.Fore.MAGENTA + 'Bingo!')
            else:
                wrolist.append((idx, ctopic[0], ctopic[1]))
                if sr == 'yes':
                    print(colorama.Fore.RED + 'Wrong: ' + ctopic[4])
        idx += 1
    print(colorama.Style.NORMAL + colorama.Fore.BLUE + 'ACCURACY: ' + str(len(corlist)) + '/' + str(idx - 1) + colorama.Style.NORMAL + colorama.Fore.RESET)
    if ipt is not None and (ipt.strip() == 'Q' or ipt.strip() == 'QUIT' or ipt.strip() == 'EXIT'):
        print(colorama.Style.NORMAL + colorama.Fore.BLUE + 'Wrong List: ' + colorama.Style.NORMAL + colorama.Fore.RESET)
        for item in wrolist:
            showTopic(item[1], item[2])
    return (corlist, wrolist)


def main(argv):
    type = 'practise'
    dbfile = ''
    level = None
    qn = '1'
    scope = None
    count = None
    mode = None
    sr = None
    star = None
    ktime = None
    if argv is not None and len(argv) > 0:
        try:
            opts, args = getopt.getopt(argv, "hf:t:l:n:p:c:m:r:s:k:", ["file=", "type=", "level=", "number=", "scope=", "count=", "mode=", "sr=", "star=", "ktime="])
        except getopt.GetoptError:
            print(
                '''Usage: avfetch.py [-f <dbfile>] [-t <type>] [-l <level>] [-n <number>] [-p <scope>] [-c <count>] [-m <mode>] [-r <showres>] [-s <star>] [-k <ktime>]\n
                Example: orath.py -f C:/orath.db -t show -l 1Z0-051 -n 10 -p 10-110 -c 100 -m random -r yes -s yes -k 120'''
            )
            sys.exit(2)

        if len(args) > 0:
            texts.extend(args)
        for opt, arg in opts:
            if opt == '-h':
                print(
                    '''Usage: avfetch.py [-f <dbfile>] [-t <type>] [-l <level>] [-n <number>] [-p <scope>] [-c <count>] [-m <mode>] [-r <showres>] [-s <star>] [-k <ktime>]\n
                    Example: orath.py -f C:/orath.db -t show -l 1Z0-051 -n 10 -p 10-110 -c 100 -m random -r yes -s yes -k 120'''
                )
                sys.exit()
            elif opt in ("-t", "--type"):
                type = arg
            elif opt in ("-f", "--file"):
                dbfile = arg
            elif opt in ("-l", "--level"):
                level = arg
            elif opt in ("-n", "--number"):
                qn = arg
            elif opt in ("-p", "--scope"):
                scope = arg
            elif opt in ("-c", "--count"):
                count = arg
            elif opt in ("-m", "--mode"):
                mode = arg
            elif opt in ("-r", "--sr"):
                sr = arg
            elif opt in ("-s", "--star"):
                star = arg
            elif opt in ("-k", "--ktime"):
                ktime = arg
            else:
                pass
        try:
            if type == 'fetch':
                if dbfile is None or dbfile.strip() == '':
                    collect(level)
                else:
                    collect(level, dbfile)
            if type == 'show':
                if dbfile is None or dbfile.strip() == '':
                    showTopic(level, qn, False)
                else:
                    showTopic(level, qn, False, dbfile)
            if type == 'practise':
                if dbfile is None or dbfile.strip() == '':
                    res = practise(level, scope, count, mode, sr, star, ktime)
                else:
                    res = practise(level, scope, count, mode, sr, star, ktime, dbfile)
            if type == 'refer':
                if dbfile is None or dbfile.strip() == '':
                    refer(level)
                else:
                    refer(level, dbfile)
        except Exception as ex:
            print('main:' + str(ex))


if __name__ == "__main__":
    main(sys.argv[1:])


# main(['-t', 'practise', '-l', '1Z0-051', '-n', '15', '-f', os.path.join(curDir(), 'orath.db')])
# main(['-t', 'practise', '-l', '1Z0-052', '-c', '2', '-s', 'y', '-m', 'random'])
# main(['-t', 'practise', '-l', '1Z0-053', '-p', '3-50', '-c', '-1', '-m', 'r', '-r', 'y', '-s', 'n', '-k', '10'])
# main(['-t', 'show', '-l', '1Z0-052', '-n', '12', '-f', os.path.join(curDir(), 'orath.db')])

'''
iconn = get_conn('D:/Git/minicode/orath/orath.db')
isql = 'select level, qn, content from `orath`'

res = fetchall(iconn, isql)
for r in res:
    level = r[0]
    qn = r[1]
    content = r[2]
    pattern = re.compile(r'(.*)(A\..*)', re.S)
    m = pattern.match(content)
    rcontent = m.group(1).strip()
    options = m.group(2).strip()
    # print(options)
    print(str(r[0]) + ':' + str(r[1]))
    isql = 'update `orath` set `content` = ?, `options` = ? where `level` = ? and `qn` = ?'
    update(get_conn('D:/Git/minicode/orath/orath.db'), isql, [(rcontent, options, level, qn)])
    # break
'''


'''
for i in range(63, 64):
    updateTopic('1Z0-051', int(i), os.path.join(curDir(), 'orath.db'))
'''

'''

'''


'''
with open('C:/Users/xshrim/Desktop/aaaaaa.txt', 'w') as wf:
    iconn = get_conn(os.path.join(curDir(), 'orath.db'))
    isql = 'select * from `orath` where `level`="1Z0-052"'
    res = fetchall(iconn, isql)
    for r in res:
        print(str(r[4]) + ':' + str(r[5]))
        wf.write(str(r[4]) + ':' + str(r[5]) + '\n')
'''
