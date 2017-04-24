import subprocess
import sys
import chardet


def syscode():
    popen = subprocess.Popen('ping', stdout = subprocess.PIPE, stderr = subprocess.PIPE)
    charcode = chardet.detect(popen.stdout.read())
    # print(charcode)
    if float(charcode['confidence']) > 0.7:
        return charcode['encoding']
    else:
        return sys.getdefaultencoding()
    
scode = syscode()    
popen = subprocess.Popen(['ping', 'www.baidu.com', '-n', '2'], stdout = subprocess.PIPE)
while True:
    next_line = popen.stdout.readline().decode(scode)
    if next_line == '' and popen.poll() != None:
        break
    sys.stdout.write(next_line)
popen.wait()
    
print('aaa')
