import sys, html2pdf

titleTag = ""
contentTag = ""
catalog = []

# url = "https://blog.csdn.net/xxb249/article/details/80790632"
# html2pdf.GenPDF(url)

if len(sys.argv) < 2 and len(catalog) < 1:
    print("no url, exit")
    sys.exit(1)
else:
    for i in range(1, len(sys.argv)):
        catalog.append(str(sys.argv[i]).strip())

for url in catalog:
    print(('Generate PDF For ' + url).center(150, '#'))
    print(('Fetch Page:' + url).center(150, '-'))
    if "studygolang.com/" in url:
        titleTag = "#title"
        contentTag = ".content.markdown-body"
    elif "blog.csdn.net" in url or "blog.csdn.com" in url:
        titleTag = ".title-article"
        contentTag = "#article_content"
    elif "cnblogs.com" in url:
        titleTag = "#cb_post_title_url"
        contentTag = ".post"
    elif "jianshu.com" in url:
        titleTag = ".title"
        contentTag = ".show-content"

    if html2pdf.CreatePDF(url, titleTag, contentTag):
        print(('Generate PDF For ' + url + " < - > Done").center(150, '#'))
    else:
        print(('Generate PDF For ' + url + " < - > Fail").center(150, '#'))
