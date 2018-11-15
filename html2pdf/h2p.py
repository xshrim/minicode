import xhtml2pdf.pisa as pisa
from pyquery import PyQuery

html = '''
<html>
    <style> 
        @font-face { 
            font-family: yahei; 
            src: url('/home/xshrim/shell/yahei.ttf'); 
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
        <div id="contentDiv"> 
            This is a 中文 on page #<pdf:pagenumber> 
        </div> 
    </body> 
</html>
'''
# pdf = pisa.CreatePDF(open('test.html','rb'),open('test.pdf','wb'))
pdf = pisa.CreatePDF(html, open('test.pdf', 'wb'))
