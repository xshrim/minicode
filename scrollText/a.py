#coding=utf-8

import sys
from PyQt5.QtWidgets import QApplication, QWidget, QPushButton
from PyQt5.QtCore import QCoreApplication

class Trans(QWidget):

    def __init__(self):
        super(Trans, self).__init__()
        self.initUI()
        

    def initUI(self):
        #self.setAttribute(QtCore.Qt.WA_NoSystemBackground, False)
        self.setAttribute(Qt.WA_TranslucentBackground, True)
        self.setWindowFlags(Qt.FramelessWindowHint)
        button = QPushButton('Close', self)
        button.clicked.connect(self.close)

if __name__ == '__main__':
    app = QApplication(sys.argv)
    trans = Trans()
    trans.show()
    sys.exit(app.exec_())