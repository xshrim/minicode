from PyQt5.QtGui import *
from PyQt5.QtCore import *
from PyQt5.QtWidgets import *

class scrollTextLabel(QLabel):
    def __init__(self, s, parent=None):
        super(scrollTextLabel, self).__init__(parent)
        self.txt = s
        self.newX = 10      
        self.t = QTimer()
        self.font = QFont('微软雅黑, verdana', 10)
        self.t.timeout.connect(self.changeTxtPosition)

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
        #painter.setPen(QColor('Red'));
        painter.setPen(QColor('transparent'));

        #以透明色绘制文字，来取得绘制后的文字宽度
        self.textRect = painter.drawText(QRect(0, -7, self.width(), 25), Qt.AlignHCenter | Qt.AlignVCenter, self.txt)

        #if self.textRect.width() != self.width():
        if self.width() != 1:

            #如果绘制文本宽度大于控件显示宽度，准备滚动：
            painter.setPen(QColor('Green'))
            painter.drawText(QRect(self.newX, -7, self.textRect.width(), 25), Qt.AlignLeft | Qt.AlignVCenter, self.txt)
            #每150ms毫秒滚动一次
            self.t.start(100)
        else:

            #如果绘制文本宽度小于控件宽度，不需要滚动，文本居中对齐
            painter.setPen(QColor('Green'));
            self.textRect = painter.drawText(QRect(0, -7, self.width(), 25), Qt.AlignHCenter | Qt.AlignVCenter, self.txt)
            self.t.stop()

class Window(QWidget):

    def __init__(self):
        super().__init__()

        self.initUI()

    def initUI(self):
        
        self.scrollLabel = scrollTextLabel('这到底')
        self.scrollLabel.setFixedWidth(700)
        
        self.gbox = QVBoxLayout()
        self.gbox.addWidget(self.scrollLabel)
        self.setLayout(self.gbox)

        self.setGeometry(300, 300, 800,100)
        self.setWindowTitle('浏览')
        self.show()
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