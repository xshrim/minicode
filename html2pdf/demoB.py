import time, html2pdf
from pyquery import PyQuery
from urllib.parse import urljoin

catalog = []

print('Fetch Catalog'.center(100, '#'))

for i in range(1):
    url = "https://studygolang.com/subject/2"
    cpage = PyQuery(url)

    for item in cpage('.article-title').items():
        catalog.insert(0, (urljoin(url, item.attr('href')), item.text()))
        print(item.text())
    time.sleep(1)

print('Generate PDF'.center(100, '#'))
for cata in catalog:
    print(('Fetch Page:' + cata[1]).center(100, '-'), end='...')
    if html2pdf.CreatePDF(cata[0], '#title', '.content.markdown-body'):
        print('done')
    time.sleep(1)