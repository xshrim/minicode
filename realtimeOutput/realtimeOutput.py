import subprocess
import sys
import chardet

def charDetect(data):
    charsets = ['UTF-8', 'UTF-8-SIG', 'GBK', 'GB2312', 'GB18030', 'UTF-16', 'UTF-32', 'BIG5', 'LATIN-1', 'ASCII', 'SHIFT_JIS', 'EUC-CN', 'EUC-TW', 'EUC-JP', 'EUC-KR', 'HKSCS', 'KOREAN', 'KOI8-R', 'KOI8-U', 'UTF-7', 'ISO-8859-1', 'ISO-8859-1', 'ISO-8859-5', 'ISO-8859-6', 'ISO-8859-7', 'ISO-8859-11', 'ISO-8859-15', 'TIS-620']
    try:
        charinfo = chardet.detect(data)
        data.decode(charinfo['encoding'])
        return str(charinfo['encoding']).upper()
    except Exception as ex:
        for chartype in charsets:
            try:
                data.decode(chartype)
                return chartype
            except:
                continue
    return ''

def syscode():
    popen = subprocess.Popen('ping', stdout = subprocess.PIPE, stderr = subprocess.PIPE)
    return charDetect(popen.stdout.read())
    '''
    if float(charcode['confidence']) > 0.7:
        return charcode['encoding']
    else:
        return sys.getdefaultencoding()
    '''
    
scode = syscode() 
popen = subprocess.Popen(['ping', 'www.baidu.com', '-n', '2'], stdout = subprocess.PIPE)
while True:
    next_line = popen.stdout.readline().decode(scode)
    if next_line == '' and popen.poll() != None:
        break
    sys.stdout.write(next_line)
popen.wait()

