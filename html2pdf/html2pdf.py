import os, re, time, zlib, socks, socket, pdfkit, chardet
from pyquery import PyQuery
from urllib import request
from urllib.parse import quote

options = {
    'page-size': 'Letter',
    'margin-top': '0.75in',
    'margin-right': '0.75in',
    'margin-bottom': '0.75in',
    'margin-left': '0.75in',
    'encoding': "utf-8",  # 支持中文
    # 'custom-header': [('Accept-Encoding', 'gzip')],
    'cookie': [
        ('cookie-name1', 'cookie-value1'),
        ('cookie-name2', 'cookie-value2'),
    ],
    'quiet': '',
}

template = '''<!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <style>
            @font-face {
                font-family: yahei;
                src: url('./yahei.ttf');
            }
            body{
                font-family: yahei;
            }
        </style>
        <title>
        </title>
    </head>
    <body>
    <body>
    </html>
    '''


def charDetect(data):
    charsets = [
        'UTF-8', 'UTF-8-SIG', 'GBK', 'GB2312', 'GB18030', 'BIG5', 'SHIFT_JIS', 'EUC-CN', 'EUC-TW', 'EUC-JP', 'EUC-KR', 'ASCII',
        'HKSCS', 'KOREAN', 'UTF-7', 'TIS-620', 'LATIN-1', 'KOI8-R', 'KOI8-U', 'ISO-8859-5', 'ISO-8859-6', 'ISO-8859-7',
        'ISO-8859-11', 'ISO-8859-15', 'UTF-16', 'UTF-32'
    ]
    try:
        charinfo = chardet.detect(data)
        data.decode(charinfo['encoding'])
        return str(charinfo['encoding']).upper()
    except Exception as ex:
        print('charDetect error:' + str(ex))
        for chartype in charsets:
            try:
                data.decode(chartype)
                return chartype
            except Exception as ex:
                print('charDetect error:' + str(ex))
                continue
    return ''


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
    headers = [('Host', 'img0.imgtn.bdimg.com'), ('Connection', 'close'), ('Cache-Control', 'max-age=0'),
               ('Accept', 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'),
               ('User-Agent',
                'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36'),
               ('Accept-Encoding', '*'), ('Accept-Language', 'en-US,zh-CN,zh;q=0.8'),
               ('If-None-Match', '90101f995236651aa74454922de2ad74'), ('Referer', 'http://www.deviantart.com/whats-hot/'),
               ('If-Modified-Since', 'Thu, 01 Jan 1970 00:00:00 GMT')]

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
            print('getHTML error:' + str(ex))
            if '403' in str(ex) or '404' in str(ex) or '502' in str(ex) or '11001' in str(ex):
                break
        i -= 1
    return contents


def CreateHTML(url, titletag, contenttag):
    try:
        page = PyQuery(template)
        data = PyQuery(
            getHTML(url, 5, 3, 1).replace('data-original-src="//',
                                          'src="http://').replace('href="//', 'href="http://').replace('src="//', 'src="http://'))
        data(".image-container-fill").remove()

        # data = PyQuery(getHTML(url, 5, 3, 1).replace('href="//', 'href="http://').replace('src="//', 'src="http://'))
        # fname = os.path.basename(url).replace('.html', '.pdf')
        if titletag[0] == '#' or titletag[0] == '.':
            title = data(titletag + ":first").text().strip()
        else:
            title = titletag
        # print("Article Title: " + title)

        page("title").append(title)
        '''
        for script in data.items("script"):
            page("head").append(script)
        for style in data.items("style"):
            page("head").append(style)
        for link in data.items("link"):
            page("head").append(link)
        '''

        for style in data.items("style"):
            page("head").append(style)
        for link in data.items("link"):
            page("head").append(link)

        # for script in data.items("script"):
        #     page("head").append(script)

        content = data(contenttag).html()

        page("body").append(content)
        '''
        with open("temp.html", "w") as wf:
            wf.write(page.outer_html())
        '''
        return page
    except Exception as ex:
        print("CreateHTML error:" + str(ex))
        return None


def CreatePDF(url, titletag, contenttag, idx=None):
    # html = CreateHTML("https://studygolang.com/articles/12263", "#title", ".content.markdown-body")
    # html = CreateHTML("https://www.cnblogs.com/nerxious/archive/2012/12/21/2827303.html", ".postTitle", ".postBody")

    filepath = ""

    try:
        html = CreateHTML(url, titletag, contenttag)
        if html is None:
            return False
        title = html("title").text().replace("/", "-")
        content = html("body").html()

        if idx is not None:
            title = str(idx) + "." + title

        if len(PyQuery(content).children()) < 2:  # markdown
            filepath = title + ".md"
            content = "# " + title + "\n" + content
            with open(filepath, "w") as wf:
                wf.write(content)
        else:  # html
            filepath = title + ".pdf"
            html("body").prepend("<h2>" + title + "</h2>")
            pdfkit.from_string(html.outer_html(), title + ".pdf", options=options)
        return True
    except Exception as ex:
        print("CreatePDF error:" + str(ex))
        if os.path.isfile(filepath):
            return True
        else:
            return False


def GenPDF(url):
    pdfkit.from_url(url, "test.pdf", options=options)
    if os.path.isfile("test.pdf"):
        return True
    else:
        return False
