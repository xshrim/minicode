# -*- coding: utf-8 -*-
import sys
from PyQt5 import QtCore
from PyQt5 import QtGui
from PyQt5 import QtWidgets


class LcdTime(QtWidgets.QFrame):
    def __init__(self, parent=None):
        super(LcdTime, self).__init__(parent)

        self.hour = QtWidgets.QLCDNumber(8, self)
        self.hour.setGeometry(10, 10, 200, 70)
        self.hour.setSegmentStyle(QtWidgets.QLCDNumber.Flat)
        self.display()

        self.timer = QtCore.QTimer()
        self.timer.timeout.connect(self.display)
        #self.connect(self.timer, QtCore.SIGNAL('timeout()'), self.display)
        self.timer.start(1000)

        self.build_tray()
        self.resize(220, 100)
        self.central()

        # 边框透明
        self.hour.setFrameShape(QtWidgets.QFrame.NoFrame)
        self.setWindowFlags(QtCore.Qt.FramelessWindowHint | QtCore.Qt.SubWindow | QtCore.Qt.WindowStaysOnTopHint)
        # 透明处理，移动需要拖动数字
        self.setAttribute(QtCore.Qt.WA_TranslucentBackground, True)
        self.setMouseTracking(True)

    def mousePressEvent(self, event):
        if event.button() == QtCore.Qt.LeftButton:
            self.dragPosition = event.globalPos() - self.frameGeometry().topLeft()
            event.accept()

    def mouseMoveEvent(self, event):
        if event.buttons() == QtCore.Qt.LeftButton:
            self.move(event.globalPos() - self.dragPosition)
            event.accept()

    def build_tray(self):
        self.trayIcon = QtWidgets.QSystemTrayIcon(self)
        self.trayIcon.setIcon(QtGui.QIcon('resource/logo.png'))
        self.trayIcon.show()
        self.trayIcon.setToolTip('时钟 -LiKui')
        self.trayIcon.activated.connect(self.trayClick)

        menu = QtWidgets.QMenu()
        normalAction = menu.addAction('正常显示')
        miniAction = menu.addAction('最小化托盘')
        exitAction = menu.addAction('退出')
        normalAction.triggered.connect(self.showNormal)
        exitAction.triggered.connect(self.exit)
        miniAction.triggered.connect(self.showMinimized)

        self.trayIcon.setContextMenu(menu)

    def exit(self):
        # 不设置Visible为False，退出后TrayIcon不会刷新
        self.trayIcon.setVisible(False)
        sys.exit(0)

    def trayClick(self, reason):
        if reason == QtWidgets.QSystemTrayIcon.DoubleClick:
            self.showNormal()
            self.repaint()

    def display(self):
        current = QtCore.QTime.currentTime()
        self.hour.display(current.toString('HH:mm:ss'))

    def showNormal(self):
        super(LcdTime, self).showNormal()
        self.repaint()

    def central(self):
        screen = QtWidgets.QDesktopWidget().screenGeometry()
        size = self.geometry()
        self.move(screen.width() - size.width(), 0)


app = QtWidgets.QApplication(sys.argv)
lcd = LcdTime()
lcd.show()
sys.exit(app.exec_())