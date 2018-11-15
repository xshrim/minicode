import os, time, pdfkit
from pyquery import PyQuery
from urllib.parse import urlparse

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
    'no-outline': None
}


def html2pdf(url, titletag, contenttag, opt):

    page = PyQuery('''<!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <style>
            @font-face {
                font-family: yahei;
                src: url('/home/xshrim/shell/yahei.ttf');
            }
            body{
                font-family: yahei;
            }
        </style>
    </head>
    <body>
    <body>
    </html>
    ''')

    data = PyQuery(url)
    # fname = os.path.basename(url).replace('.html', '.pdf')
    if titletag[0] == '#' or titletag[0] == '.':
        title = data(titletag).text().strip() + '.pdf'
    else:
        title = titletag + '.pdf'
    title = data(titletag).text().strip() + '.pdf'
    content = data(contenttag).html()
    page("body").append(content)
    # print(page.html())

    pdfkit.from_string(page.html(), title, options=opt)
    '''
    if os.path.exists(fname):
        os.rename(fname, title)
    '''


html2pdf('https://blog.csdn.net/qq_34598667/article/details/83959393', '.title-article', '.markdown_views.prism-atom-one-dark',
         options)
os._exit(0)

# html2pdf('https://www.cnblogs.com/CloudMan6/p/6693772.html', options)
catalog = []

print('Fetch Catalog'.center(100, '#'))

for i in range(2):
    url = "https://www.cnblogs.com/CloudMan6/tag/Docker/default.html?page=" + str(i + 1)
    cpage = PyQuery(url)

    for item in cpage('.PostList').items('a'):
        catalog.insert(0, (item.attr('href'), item.text()))
        print(item.text())
    time.sleep(1)

print('Generate PDF'.center(100, '#'))
for cata in catalog:
    print(('Fetch Page:' + cata[1]).center(100, '-'), end='...')
    html2pdf(cata[0], '#cb_post_title_url', '.post', options)
    print('done')
    time.sleep(1)
# pdfkit.from_url('http://www.baidu.com', 'out.pdf', options=options)
