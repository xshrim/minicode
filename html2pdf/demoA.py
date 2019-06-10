import time, html2pdf
from pyquery import PyQuery

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
    html2pdf.CreatePDF(cata[0], '#cb_post_title_url', '.post')
    print('done')
    time.sleep(1)
