import sys
import time
import random
from danmu import DanMuClient
from PyQt5.QtCore import *
from PyQt5.QtGui import *
from PyQt5.QtWidgets import *


class scrollTextLabel(QLabel):
    def __init__(self, s, parent=None):
        super(scrollTextLabel, self).__init__(parent)
        self.setAlignment(Qt.AlignLeft)
        self.txt = s
        self.setFixedHeight(30)
        self.setFixedWidth(100)
        self.t = QTimer()
        self.font = QFont('微软雅黑, verdana', 15, QFont.Bold)
        self.t.timeout.connect(self.changeTxtPosition)
        self.update()

    def getText(self):
        return self.txt

    def setText(self, s):
        if s is not None and s.strip() != '':
            self.txt = s
            if len(s) < 10:
                self.setFixedWidth(200)
            elif len(s) < 20:
                self.setFixedWidth(400)
            elif len(s) < 30:
                self.setFixedWidth(600)
            else:
                self.setFixedWidth(800)
            self.t.start(50)

    def changeTxtPosition(self):
        if self.txt is not None and self.txt.strip() != '':
            if self.pos().x() <= 0:
                # self.hide()
                self.txt = ''
                self.t.stop()
                # self.move(self.desktop.width() / 2 - 50, self.pos().y())
            else:
                self.move(self.pos().x() - 5, self.pos().y())
        self.update()

    def paintEvent(self, event):
        painter = QPainter(self)
        painter.setFont(self.font)
        painter.setPen(QColor('green'))
        self.textRect = painter.drawText(QRect(0, -7, self.width(), 30), Qt.AlignLeft | Qt.AlignVCenter, self.txt)


class Worker(QThread):

    item_changed_signal = pyqtSignal(str)

    def __init__(self, total=0, parent=None):
        super().__init__(parent)
        self.total = total

    def __del__(self):
        self.wait()

    def run(self):
        def pp(msg):
            self.item_changed_signal.emit(msg.encode(sys.stdin.encoding, 'ignore').decode(sys.stdin.encoding))
            self.sleep(1)
            # self.item_changed_signal.emit(msg.encode(sys.stdin.encoding, 'ignore').decode(sys.stdin.encoding))
            # print(msg.encode(sys.stdin.encoding, 'ignore').decode(sys.stdin.encoding))

        dmc = DanMuClient('https://www.douyu.com/493462')
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

        self.initUI()

    def initUI(self):
        self.desktop = QApplication.desktop()
        self.items = []
        # self.w = QTimer()
        # self.w.timeout.connect(self.changeTxt)
        '''
        self.scrollLabel = scrollTextLabel()
        # self.scrollLabel.setAlignment(Qt.AlignRight)
        self.scrollLabel.setFixedWidth(desktop.width() / 2 - 10)
        self.scrollLabel.setText('这')
        '''

        self.setAttribute(Qt.WA_TranslucentBackground, True)
        self.setWindowFlags(Qt.FramelessWindowHint)

        '''
        self.tl0 = scrollTextLabel('', self)
        # tl.setText('这是什么鬼')
        self.tl0.move(self.desktop.width() / 2 - 50, 10)

        self.tl1 = scrollTextLabel('', self)
        # tl.setText('这是什么鬼')
        self.tl1.move(self.desktop.width() / 2 - 50, 20)
        '''

        while len(self.items) < 100:
            self.items.append(scrollTextLabel('', self))
            self.items[-1].move(self.desktop.width() / 2 - 50, random.choice([10, 40, 70, 100, 130, 160, 190, 220, 250, 280, 310, 340, 370, 400, 430, 460, 490]))

        self.setGeometry(0, 0, self.desktop.width() / 2, 550)
        self.setWindowTitle('浏览')

        self.show()
        self.thread = Worker(10)
        self.thread.item_changed_signal.connect(self.showDanmu)
        self.thread.start()
        # self.w.start(1000)

    def showDanmu(self, text):
        # self.addItem('hhh')
        eitem = [item for item in self.items if item.getText() == '']
        for titem in eitem:
            titem.move(self.desktop.width() / 2 - 50, titem.pos().y())
        if len(eitem) > 0:
            citem = random.choice(eitem)
            citem.setText(text)

    def changeTxt(self):
        # self.addItem('hhh')
        eitem = [item for item in self.items if item.getText() == '']
        for titem in eitem:
            titem.move(self.desktop.width() / 2 - 50, titem.pos().y())
        if len(eitem) > 0:
            citem = random.choice(eitem)
            citem.setText('这是什么鬼')


if __name__ == '__main__':

    app = QApplication(sys.argv)
    window = Window()
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
