
import os
import sys
import sqlite3
from PyQt5.QtWidgets import (QWidget, QLabel, QPushButton, QLineEdit, QTextEdit, QHBoxLayout, QVBoxLayout, QComboBox, QRadioButton, QFileDialog, QApplication)
from PyQt5.QtGui import QPixmap

######################################### DB START#########################################


def get_conn(path):
    conn = sqlite3.connect(path)
    if os.path.exists(path) and os.path.isfile(path):
        # print('硬盘上面:[{}]'.format(path))
        return conn
    else:
        conn = None
        # print('内存上面:[:memory:]')
        return sqlite3.connect(':memory:')


def get_cursor(conn):
    if conn is not None:
        return conn.cursor()
    else:
        return get_conn('').cursor()


def drop_table(conn, table):
    if table is not None and table != '':
        sql = 'DROP TABLE IF EXISTS ' + table
        # print('执行sql:[{}]'.format(sql))
        cu = get_cursor(conn)
        cu.execute(sql)
        conn.commit()
        # print('删除数据库表[{}]成功!'.format(table))
        close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def create_table(conn, sql):
    if sql is not None and sql != '':
        cu = get_cursor(conn)
        # print('执行sql:[{}]'.format(sql))
        cu.execute(sql)
        conn.commit()
        # print('创建数据库表成功!'
        close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def close_all(conn, cu):
    try:
        if cu is not None:
            cu.close()
    finally:
        if cu is not None:
            cu.close()


def save(conn, sql, data):
    if sql is not None and sql != '':
        if data is not None:
            cu = get_cursor(conn)
            for d in data:
                # print('执行sql:[{}],参数:[{}]'.format(sql, d))
                cu.execute(sql, d)
                conn.commit()
            close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def fetchall(conn, sql):
    if sql is not None and sql != '':
        cu = get_cursor(conn)
        # print('执行sql:[{}]'.format(sql))
        cu.execute(sql)
        r = cu.fetchall()
        if len(r) > 0:
            for e in range(len(r)):
                print(r[e])
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def fetchone(conn, sql, data):
    if sql is not None and sql != '':
        if data is not None:
            # Do this instead
            d = (data,)
            cu = get_cursor(conn)
            # print('执行sql:[{}],参数:[{}]'.format(sql, data))
            cu.execute(sql, d)
            r = cu.fetchall()
            if len(r) > 0:
                for e in range(len(r)):
                    print(r[e])
        else:
            logging.error('the [{}] equal None!'.format(data))
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def update(conn, sql, data):
    if sql is not None and sql != '':
        if data is not None:
            cu = get_cursor(conn)
            for d in data:
                # print('执行sql:[{}],参数:[{}]'.format(sql, d))
                cu.execute(sql, d)
                conn.commit()
            close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))


def delete(conn, sql, data):
    if sql is not None and sql != '':
        if data is not None:
            cu = get_cursor(conn)
            for d in data:
                # print('执行sql:[{}],参数:[{}]'.format(sql, d))
                cu.execute(sql, d)
                conn.commit()
            close_all(conn, cu)
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))

######################################### DB END#########################################


class Window(QWidget):

    def __init__(self):
        super().__init__()

        self.initUI()

    def initUI(self):

        self.preButton = QPushButton('上一个')
        self.nextButton = QPushButton('下一个')
        self.viewModeBox = QComboBox()
        self.viewModeBox.addItems(['顺序', '随机'])
        self.dbButton = QPushButton('载入数据库')
        self.dbButton.clicked.connect(self.showDialog)

        self.favorButton = QRadioButton('收藏')

        self.dbEdit = QLineEdit()

        self.infoEdit = QTextEdit()
        self.infoEdit.setReadOnly(True)
        self.infoEdit.setFixedSize(300, 300)
        self.infoEdit.append('adddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddaa')

        self.picLabel = QLabel()
        self.picLabel.setPixmap(QPixmap('D:\月儿.png').scaled(300, 300))

        self.hboxhead = QHBoxLayout()
        self.hboxhead.addWidget(self.dbEdit)
        self.hboxhead.addWidget(self.dbButton)
        self.hboxhead.addWidget(self.favorButton)
        self.hboxhead.addWidget(self.viewModeBox)
        self.hboxhead.addWidget(self.preButton)
        self.hboxhead.addWidget(self.nextButton)

        self.hboxinfo = QHBoxLayout()
        self.hboxinfo.addWidget(self.picLabel)
        self.hboxinfo.addWidget(self.infoEdit)

        self.vbox = QVBoxLayout()
        self.vbox.addLayout(self.hboxhead)
        self.vbox.addLayout(self.hboxinfo)

        self.setLayout(self.vbox)

        self.setGeometry(300, 300, 30, 350)
        self.setWindowTitle('浏览')
        self.show()

    def showDialog(self):
        fname = QFileDialog.getOpenFileName(self, '选择数据库', os.path.join(os.path.expanduser("~"), 'Desktop'))
        if fname[0]:
            self.dbEdit.setText(fname[0])


if __name__ == '__main__':

    app = QApplication(sys.argv)
    window = Window()
    sys.exit(app.exec_())
