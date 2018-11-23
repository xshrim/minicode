import xhtml2pdf.pisa as pisa
from pyquery import PyQuery

template = '''
<html>
    <style> 
        @font-face { 
            font-family: yahei; 
            src: url('/home/xshrim/lab/yahei.ttf'); 
        } 
        body{ 
            font-family: yahei; 
        } 
        @page { 
            margin: 1cm; 
            margin-bottom: 2.5cm; 
            font-family: yahei; 
            @frame content { 
                -pdf-frame-content: contentDiv; 
                margin-left: 1cm; 
                margin-right: 1cm; 
                height: 1cm; 
                font-family: yahei; 
            } 
            @frame footer { 
                -pdf-frame-content: footerDiv; 
                bottom: 2cm; 
                margin-left: 1cm; 
                margin-right: 1cm; 
                height: 1cm; 
                font-family: yahei; 
            } 
        } 
    </style> 
    <body> 
        <div id="headerDiv">
            <h3 id="title"><h3>
        </div>
        <div id="contentDiv"> 
        </div> 
    </body> 
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
    page("body")("#title").append(title)
    page("body")("#contentDiv").append(content)
    return page


# html = CreateHTML("https://studygolang.com/articles/12263", "#title", ".content.markdown-body")
html = CreateHTML("https://www.cnblogs.com/nerxious/archive/2012/12/21/2827303.html", ".postTitle", ".postBody")

title = html("body")("#title").text()
print(title)
content = html("body")("#contentDiv").html()
if len(PyQuery(content).children()) < 1:  # markdown
    content = "# " + title + "\n" + content
    with open(title + ".md", "w") as wf:
        wf.write(content)
else:  # html
    # pdf = pisa.CreatePDF(open('test.html','rb'),open('test.pdf','wb'))
    pdf = pisa.CreatePDF(html.outer_html(), open(title + '.pdf', 'wb'))
