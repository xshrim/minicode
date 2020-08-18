import time, html2pdf
from pyquery import PyQuery

catalog = []

print('Fetch Catalog'.center(100, '#'))

for i in range(1):
    url = "https://blog.csdn.net/wangshubo1989/article/category/6729509/" + str(i + 1)
    cpage = PyQuery(url)

    for item in cpage('.article-list').items('h4'):
        catalog.append((item('a').attr('href'), item('a').text().strip("原 ")))
        print(item.text().strip("原 "))
    time.sleep(1)

#catalog.append(("https://blog.csdn.net/wangshubo1989/article/details/79270450", "设计模式"))

print('Generate PDF'.center(100, '#'))
for idx, cata in enumerate(reversed(catalog)):
    print(('Fetch Page:' + cata[1]).center(100, '-'), end='...')
    html2pdf.CreatePDF(cata[0], '.title-article', '.blog-content-box', idx)
    print('done')
    time.sleep(1)
