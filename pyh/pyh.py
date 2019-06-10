# coding:utf-8
# @file: pyh.py
# @purpose: a HTML tag generator
# @author: Emmanuel Turlay <turlay@cern.ch>

__doc__ = """The pyh.py module is the core of the PyH package. PyH lets you
generate HTML tags from within your python code.
See http://code.google.com/p/pyh/ for documentation.
"""
__author__ = "Emmanuel Turlay <turlay@cern.ch>"
__version__ = '$Revision: 63 $'
__date__ = '$Date: 2010-05-21 03:09:03 +0200 (Fri, 21 May 2010) $'

from sys import _getframe, stdout, modules, version
nOpen = {}

# nl = '\n'
doctype = '<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">\n'
charset = '\n        <meta http-equiv="Content-Type" content="text/html;charset=utf-8" />\n'
# charset = dict(http-equiv="Content-Type", content="text/html;charset=utf-8")

# 添加style,colgroup标签,thead，取消了一个重复的script标签
tags = [
    'ul', '!--', '!DOCTYPE', 'a', 'abbr', 'acronym', 'address', 'applet', 'area', 'article', 'aside', 'audio', 'b', 'base',
    'basefont', 'bdi', 'bdo', 'big', 'blockquote', 'body', 'br', 'button', 'canvas', 'caption', 'center', 'cite', 'code', 'col',
    'colgroup', 'command', 'datalist', 'dd', 'del', 'details', 'dfn', 'dialog', 'dir', 'div', 'dl', 'dt', 'em', 'embed',
    'fieldset', 'figcaption', 'figure', 'font', 'footer', 'form', 'frame', 'frameset', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'head',
    'header', 'hr', 'html', 'i', 'iframe', 'img', 'input', 'ins', 'kbd', 'keygen', 'label', 'legend', 'li', 'link', 'main', 'map',
    'mark', 'menu', 'menuitem', 'meta', 'meter', 'nav', 'noframes', 'noscript', 'object', 'ol', 'optgroup', 'option', 'output',
    'p', 'param', 'pre', 'progress', 'q', 'rp', 'rt', 'ruby', 's', 'samp', 'script', 'section', 'select', 'small', 'source',
    'span', 'strike', 'strong', 'style', 'sub', 'summary', 'sup', 'table', 'tbody', 'td', 'textarea', 'tfoot', 'th', 'thead',
    'time', 'title', 'tr', 'track', 'tt', 'u', 'ul', 'var', 'video', 'wbr'
]

selfClose = [
    'meta', 'input', 'img', 'link', 'br', 'hr', 'base', 'col', 'frame', 'area', 'param', 'object', 'embed', 'keygen', 'source'
]


class Tag(list):
    tagname = ''

    def __init__(self, *arg, **kw):
        self.attributes = kw
        if self.tagname:
            name = self.tagname
            self.isSeq = False
        else:
            name = 'sequence'
            self.isSeq = True
        self.id = kw.get('id', name)
        # self.extend(arg)
        self.indentlvl = 0
        for a in arg:
            self.addObj(a)

    def __iadd__(self, obj):
        "重载累加操作符"
        if isinstance(obj, Tag) and obj.isSeq:
            for o in obj:
                self.addObj(o)
        else:
            self.addObj(obj)

        return self

    def addObj(self, obj):
        if not isinstance(obj, Tag):
            obj = str(obj)
            # obj=unicode(obj) python3 不支持该函数，所以在此将其注释
            pass
        id = self.setID(obj)
        setattr(self, id, obj)
        self.append(obj)

    def setID(self, obj):
        if isinstance(obj, Tag):
            id = obj.id
            n = len([t for t in self if isinstance(t, Tag) and t.id.startswith(id)])
        else:
            id = 'content'
            n = len([t for t in self if not isinstance(t, Tag)])
        if n:
            id = '%s_%03i' % (id, n)
        if isinstance(obj, Tag):
            obj.id = id
        return id

    def __add__(self, obj):
        if self.tagname:
            return Tag(self, obj)
        self.addObj(obj)
        return self

    def __lshift__(self, obj):
        "操作符<<重载"
        self += obj
        if isinstance(obj, Tag):
            return obj

    # 添加一个文件参数，可以输出代码段到指定文件，因此可以不必生成完整的页面

    def render(self, depth=0, file_name=None):
        result = ''

        if self.tagname:
            result = '    ' * depth + '<%s%s%s>' % (self.tagname, self.renderAtt(), self.selfClose() * ' /')

        if not self.selfClose():
            for c in self:
                if isinstance(c, Tag):
                    result += '\n' + c.render(depth + 1)
                else:
                    result += c
            if self.tagname:
                if result[-1] != '\n':
                    result += '</%s>' % self.tagname
                else:
                    result += '    ' * depth + '</%s>' % self.tagname
        result += '\n'

        if file_name:
            with open(file_name, "w") as f:
                f.write(result)
        else:
            return result

    def renderAtt(self):
        result = ''
        for n, v in self.attributes.items():
            if n != 'txt' and n != 'open':
                if n == 'cl':
                    n = 'class'
                result += ' %s="%s"' % (n, v)
        return result

    def selfClose(self):
        return self.tagname in selfClose


def TagFactory(name):
    class f(Tag):
        tagname = name

    f.__name__ = name
    return f


thisModule = modules[__name__]

for t in tags:
    setattr(thisModule, t, TagFactory(t))


def ValidW3C():
    out = a(
        img(src='http://www.w3.org/Icons/valid-xhtml10', alt='Valid XHTML 1.0 Strict'),
        href='http://validator.w3.org/check?uri=referer')
    return out


class PyH(Tag):
    tagname = 'html'

    def __init__(self, name='MyPyHPage'):
        self += head()
        self += body()
        self.attributes = dict(xmlns='http://www.w3.org/1999/xhtml', lang='en')
        # self.head += meta(content='text/html', charset='utf-8')
        self.head += charset
        self.head += title(name)

    def __iadd__(self, obj):
        if isinstance(obj, head) or isinstance(obj, body):
            self.addObj(obj)
        elif isinstance(obj, meta) or isinstance(obj, link):
            self.head += obj
        else:
            self.body += obj
            id = self.setID(obj)
            setattr(self, id, obj)
        return self

    def addJS(self, *arg):
        for f in arg:
            self.head += script(type='text/javascript', src=f)

    def addCSS(self, *arg):
        for f in arg:
            self.head += link(rel='stylesheet', type='text/css', href=f)

    def addStyleSnippet(self, filename):
        "从某个文件中导入css代码段"
        with open(filename, "r") as f:
            txt = f.read()
            self.head += style(txt, type="text/css")

    def addScriptSnippet(self, filename):
        "从某个文件中导入js代码段"
        with open(filename, "r") as f:
            txt = f.read()
            self.head += script(txt, type="text/javascript")

    def printOut(self, file='', file_operation="w"):  #添加一个默认参数用来设置编码方式
        if file:
            f = open(file, file_operation)  #添加一个文件操作参数，实际使用过程中迭代page时不一定每次都要'w'覆盖，有时需要'a'
        else:
            f = stdout
        f.write(doctype)
        #        f.write(unicode(self.render()).encode(encodetype))
        f.write(self.render())
        f.flush()
        if file:
            f.close()


class TagCounter:
    _count = {}
    _lastOpen = []
    for t in tags:
        _count[t] = 0

    def __init__(self, name):
        self._name = name

    def open(self, tag):
        if isLegal(tag):
            self._count[tag] += 1
            self._lastOpen += [tag]

    def close(self, tag):
        if isLegal(tag) and self._lastOpen[-1] == tag:
            self._count[tag] -= 1
            self._lastOpen.pop()
        else:
            print('Cross tagging is wrong')

    def isAllowed(self, tag, open):
        if not open and self.isClosed(tag):
            print('TRYING TO CLOSE NON-OPEN TAG: %s' % tag)
            return False
        return True

    def isOpen(self, tag):
        if isLegal(tag): return self._count[tag]

    def isClosed(self, tag):
        if isLegal(tag): return not self._count[tag]


def isLegal(tag):
    if tag in tags: return True
    else:
        print('ILLEGAL TAG: %s' % tag)
        return False