import sys
import random
import win32gui
from danmu import DanMuClient
from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *


class scrollTextLabel(QLabel):
    def __init__(self, txt, txtsize, txtcolor, rate, parent=None):
        super(scrollTextLabel, self).__init__(parent)
        self.setAlignment(Qt.AlignLeft)
        self.txt = txt
        self.txtsize = int(txtsize)
        self.txtcolor = txtcolor
        self.rate = int(rate)
        self.setFixedHeight(20)
        self.setFixedWidth(100)
        self.t = QTimer()
        self.font = QFont('微软雅黑, verdana', self.txtsize, QFont.Bold)
        self.t.timeout.connect(self.changeTxtPosition)
        self.update()

    def getText(self):
        return self.txt

    def setText(self, txt):
        if txt is not None and txt.strip() != '':
            self.txt = txt
            self.setFixedWidth(self.getTextPixel()[0])
            print(str(self.pos().x()) + ':' + str(self.pos().y()) + '->' + self.txt)
            self.t.start(self.rate)

    def getTextPixel(self):
        metrics = QFontMetrics(self.font)
        return (metrics.width(self.txt), metrics.height())

    def changeTxtPosition(self):
        if self.txt is not None and self.txt.strip() != '':
            if self.pos().x() <= 0:
                # self.hide()
                if len(self.txt) > 0:
                    self.txt = self.txt[1:]
                else:
                    self.t.stop()
                # self.move(self.desktop.width() / 2 - 50, self.pos().y())
            else:
                self.move(self.pos().x() - 5, self.pos().y())
        self.update()

    def paintEvent(self, event):
        painter = QPainter(self)
        painter.setFont(self.font)
        painter.setPen(QColor(self.txtcolor))
        self.textRect = painter.drawText(QRect(0, -7, self.width(), self.getTextPixel()[1] + 5), Qt.AlignLeft | Qt.AlignVCenter, self.txt)


class Worker(QThread):

    item_changed_signal = pyqtSignal(str)

    def __init__(self, url='', parent=None):
        super().__init__(parent)
        self.url = url

    def __del__(self):
        self.wait()

    def run(self):
        def pp(msg):
            self.item_changed_signal.emit(msg.encode(sys.stdin.encoding, 'ignore').decode(sys.stdin.encoding))
            # self.sleep(1)
            # self.item_changed_signal.emit(msg.encode(sys.stdin.encoding, 'ignore').decode(sys.stdin.encoding))
            # print(msg.encode(sys.stdin.encoding, 'ignore').decode(sys.stdin.encoding))

        dmc = DanMuClient(self.url)
        if not dmc.isValid():
            print('Url not valid')

        @dmc.danmu
        def danmu_fn(msg):
            pp(msg['Content'])
            # pp('[%s] %s' % (msg['NickName'], msg['Content']))

        @dmc.gift
        def gift_fn(msg):
            pass
            # pp('[%s] sent a gift!' % msg['NickName'])

        @dmc.other
        def other_fn(msg):
            pass
            # pp('Other message received')

        dmc.start(blockThread=True)
        # for i in range(self.total):
        #    self.item_changed_signal.emit(str(i))
        #    self.sleep(1)


class Window(QWidget):

    def __init__(self):
        super().__init__()
        self.initArgs()
        self.initUI()

    def initArgs(self):
        self.url = 'https://www.douyu.com/2124270'
        self.mode = 'danmu'
        self.fromv = 0
        self.tov = 300
        self.color = 'yellow'
        self.size = 15
        self.rate = 40
        self.items = []
        self.desktop = QApplication.desktop()
        if self.desktop.screenCount() > 1:
            self.screenWidth = self.desktop.width() / 2
        else:
            self.screenWidth = self.desktop.width()
        argv = sys.argv[1:]
        if argv is not None and len(argv) > 0:
            try:
                opts, args = getopt.getopt(argv, "hu:m:f:t:c:s:r:", ["url=", "mode=", "from=", "to=", "color=", "size=", "rate="])
            except getopt.GetoptError:
                print(
                    '''Usage: danmuText.py [-u <liveurl>] [-m <danmumode>] [-f <showfrom>] [-t <showto>] [-c <danmucolor>] [-s <danmusize>] [-r <danmurate>]\n
                    Example: danmuText.py -u https://www.douyu.com/748396 -m danmu -f 0 -t 400 -c green -s 15 -r 50'''
                )
                sys.exit(2)

            if len(args) > 0:
                texts.extend(args)
            for opt, arg in opts:
                if opt == '-h':
                    print(
                        '''Usage: danmuText.py [-u <liveurl>] [-m <danmumode>] [-f <showfrom>] [-t <showto>] [-c <danmucolor>] [-s <danmusize>] [-r <danmurate>]\n
                        Example: danmuText.py -u https://www.douyu.com/748396 -m danmu -f 0 -t 400 -c green -s 15 -r 50'''
                    )
                    sys.exit(2)
                if opt in ("-u", "--url"):
                    self.url = str(arg).replace('\'', '').replace('\"', '').strip()
                elif opt in ("-m", "--mode"):
                    self.mode = str(arg).replace('\'', '').replace('\"', '').strip()
                elif opt in ("-f", "--from"):
                    self.fromv = int(str(arg).replace('\'', '').replace('\"', '').strip())
                    if self.fromv < 0:
                        self.fromv = 0
                elif opt in ("-t", "--to"):
                    self.tov = int(str(arg).replace('\'', '').replace('\"', '').strip())
                    if self.tov > self.screenWidth:
                        self.tov = self.screenWidth
                    if self.fromv > self.tov:
                        self.fromv, self.tov = self.tov, self.fromv
                    if self.tov - self.fromv < 100:
                        self.tov = self.fromv + 100
                elif opt in ("-c", "--color"):
                    self.color = str(arg).replace('\'', '').replace('\"', '').strip()
                elif opt in ("-s", "--size"):
                    self.size = int(str(arg).replace('\'', '').replace('\"', '').strip())
                    if self.size < 5:
                        self.size = 5
                elif opt in ("-r", "--rate"):
                    self.rate = int(str(arg).replace('\'', '').replace('\"', '').strip())
                    if self.rate < 1:
                        self.rate = 1
                else:
                    pass

    def initUI(self):

        # self.w = QTimer()
        # self.w.timeout.connect(self.changeTxt)
        '''
        self.scrollLabel = scrollTextLabel()
        # self.scrollLabel.setAlignment(Qt.AlignRight)
        self.scrollLabel.setFixedWidth(self.screenWidth - 10)
        self.scrollLabel.setText('这')
        '''

        self.setAttribute(Qt.WA_TranslucentBackground, True)
        self.setWindowFlags(Qt.FramelessWindowHint)

        '''
        self.tl0 = scrollTextLabel('', self)
        # tl.setText('这是什么鬼')
        self.tl0.move(self.screenWidth - 50, 10)

        self.tl1 = scrollTextLabel('', self)
        # tl.setText('这是什么鬼')
        self.tl1.move(self.screenWidth - 50, 20)
        '''

        fslabel = scrollTextLabel('', self.size, self.color, self.rate, self)
        fslabel.move(self.screenWidth, 10)
        self.items.append(fslabel)
        loc = [10]
        while loc[-1] < self.tov - self.fromv - fslabel.getTextPixel()[1] - 30:
            loc.append(loc[-1] + fslabel.getTextPixel()[1] + 10)
        while len(self.items) < 100:
            self.items.append(scrollTextLabel('', self.size, self.color, self.rate, self))
            self.items[-1].move(self.screenWidth, random.choice(loc))

        self.setGeometry(0, self.fromv, self.screenWidth, self.tov - self.fromv)
        self.setWindowTitle('爱弹幕')

        self.show()
        self.thread = Worker(self.url)
        self.thread.item_changed_signal.connect(self.showDanmu)
        self.thread.start()
        # self.w.start(1000)

    def showDanmu(self, text):
        sitem = None
        eitem = [item for item in self.items if item.getText() == '']
        for titem in eitem:
            titem.move(self.screenWidth, titem.pos().y())
        if len(eitem) > 0:
            for h in sorted(set([item.pos().y() for item in eitem])):
                citem = random.choice([item for item in eitem if item.pos().y() == h])
                tmpitem = [item for item in self.items if item.pos().y() == citem.pos().y() and item.getText() != '']
                if len(tmpitem) > 0:
                    mitem = max(tmpitem, key=lambda x: x.pos().x())
                    if mitem.pos().x() + mitem.getTextPixel()[0] < self.screenWidth:
                        sitem = citem
                        break
                else:
                    sitem = citem
                    break
            '''
            while True:
                try:
                    citem = random.choice(eitem)
                    tmpitem = [item for item in self.items if item.pos().y() == citem.pos().y() and item.getText() != '']
                    if len(tmpitem) > 0:
                        mitem = max(tmpitem, key=lambda x: x.pos().x())
                        if mitem.pos().x() + mitem.getTextPixel()[0] < self.screenWidth:
                            break
                    else:
                        break
                except Exception as ex:
                    print(str(ex))
            '''
            if sitem is None:
                sitem = random.choice([item for item in self.items if item.getText() == ''])
            sitem.setText(text)
            '''
            if sitem is not None:
                sitem.setText(text)
            else:
                print('miss:' + text)
            '''

    '''
    def changeTxt(self):
        # self.addItem('hhh')
        eitem = [item for item in self.items if item.getText() == '']
        for titem in eitem:
            titem.move(self.screenWidth, titem.pos().y())
        if len(eitem) > 0:
            citem = random.choice(eitem)
            citem.setText('这是什么鬼')
    '''


def windowEnumerationHandler(hwnd, top_windows):
    top_windows.append((hwnd, win32gui.GetWindowText(hwnd)))


if __name__ == '__main__':

    app = QApplication(sys.argv)
    window = Window()
    top_windows = []
    win32gui.EnumWindows(windowEnumerationHandler, top_windows)
    for i in top_windows:
        if "爱弹幕" in i[1].lower():
            win32gui.ShowWindow(i[0], 5)
            win32gui.SetForegroundWindow(i[0])
            win32gui.SetWindowPos(i[0], -1, 0, 0, 0, 0, 3)
            break
    sys.exit(app.exec_())
'''
from PyQt5.QtGui import *
from PyQt5.QtCore import *
try:
    _fromUtf8 = QString.fromUtf8
except AttributeError:
    def _fromUtf8(s):
        return s

class scrollTextLabel(QLabel):
    def __init__(self, parent=None):
        super(scrollTextLabel, self).__init__(parent)
        self.txt = QString()
        self.newX = 10
        self.t = QTimer()
        self.font = QFont(_fromUtf8('微软雅黑, verdana'), 8)
        self.connect(self.t, SIGNAL("timeout()"), self.changeTxtPosition)

    def changeTxtPosition(self):
        if not self.parent().isVisible():
            # 如果parent不可见，则停止滚动，复位
            self.t.stop()
            self.newX = 10
            return
        if self.textRect.width() + self.newX > 0:
        #每次向前滚动5像素
            self.newX -= 5
        else:
            self.newX = self.width()
        self.update()

    #用drawText来绘制文字，不再需要setText，重写
    def setText(self, s):
        self.txt = s

        #滚动起始位置设置为10,留下视觉缓冲
        #以免出现 “没注意到第一个字是什么” 的情况
        self.newX = 10
        self.update()

    def paintEvent(self, event):
        painter = QPainter(self)
        painter.setFont(self.font)
        #设置透明颜色
        painter.setPen(QColor('transparent'));

        #以透明色绘制文字，来取得绘制后的文字宽度
        self.textRect = painter.drawText(QRect(0, -7, self.width(), 25), Qt.AlignHCenter | Qt.AlignVCenter, self.txt)

        if self.textRect.width() > self.width():
            #如果绘制文本宽度大于控件显示宽度，准备滚动：
            painter.setPen(QColor(255, 255, 255, 255))
            painter.drawText(QRect(self.newX, -7, self.textRect.width(), 25), Qt.AlignLeft | Qt.AlignVCenter, self.txt)
            #每150ms毫秒滚动一次
            self.t.start(150)
        else:
            #如果绘制文本宽度小于控件宽度，不需要滚动，文本居中对齐
            painter.setPen(QColor(255, 255, 255, 255));
            self.textRect = painter.drawText(QRect(0, -7, self.width(), 25), Qt.AlignHCenter | Qt.AlignVCenter, self.txt)
            self.t.stop()
'''
