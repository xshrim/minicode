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

template = '''<!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <style>
            @font-face {
                font-family: yahei;
                src: url('/home/xshrim/lab/yahei.ttf');
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


def CreateHTML(url, titletag, contenttag):
    page = PyQuery(template)
    data = PyQuery(url)

    # fname = os.path.basename(url).replace('.html', '.pdf')
    if titletag[0] == '#' or titletag[0] == '.':
        title = data(titletag).text().strip()
    else:
        title = titletag

    content = data(contenttag).html()
    page("title").append(title)
    page("body").append(content)
    return page


def CreatePDF(url, titletag, contenttag):
    # html = CreateHTML("https://studygolang.com/articles/12263", "#title", ".content.markdown-body")
    # html = CreateHTML("https://www.cnblogs.com/nerxious/archive/2012/12/21/2827303.html", ".postTitle", ".postBody")

    try:
        html = CreateHTML(url, titletag, contenttag)
        title = html("title").text()
        content = html("body").html()
        if len(PyQuery(content).children()) < 2:  # markdown
            content = "# " + title + "\n" + content
            with open(title + ".md", "w") as wf:
                wf.write(content)
        else:  # html
            html("body").prepend("<h2>" + title + "</h2>")
            pdfkit.from_string(html.outer_html(), title + ".pdf", options=options)
        return True
    except Exception as ex:
        print(ex)
        return False
