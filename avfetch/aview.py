
import os
import sys
import sqlite3
from random import choice
from collections import OrderedDict
from PyQt5.QtWidgets import (QWidget, QLabel, QPushButton, QLineEdit, QTextEdit, QHBoxLayout, QVBoxLayout, QComboBox, QCheckBox, QRadioButton, QFileDialog, QApplication)
from PyQt5.QtGui import (QPixmap, QImage)


class av(object):
    '''
    code = ''
    title = ''
    issuedate = ''
    length = ''
    mosaic = ''
    director = ''
    manufacturer = ''
    publisher = ''
    series = ''
    category = ''
    actors = ''
    favor = ''
    coverlink = ''
    cover = b''
    link = ''
    '''

    def __init__(self, code, title, issuedate, length, mosaic, director, manufacturer, publisher, series, category, actors, favor, coverlink, cover, link):
        self.code = code
        self.title = title
        self.issuedate = issuedate
        self.length = length
        self.mosaic = mosaic
        self.director = director
        self.manufacturer = manufacturer
        self.publisher = publisher
        self.series = series
        self.category = category
        self.actors = actors
        self.favor = favor
        self.coverlink = coverlink
        if isinstance(cover, bytes):
            self.cover = cover
        else:
            self.cover = str(cover).encode()
        self.link = link

    def __str__(self):
        '''
        print('番号:'.rjust(5) + self.code)
        print('标题:'.rjust(5) + self.title)
        print('日期:'.rjust(5) + self.issuedate)
        print('时长:'.rjust(5) + self.length)
        print('修正:'.rjust(5) + self.mosaic)
        print('导演:'.rjust(5) + self.director)
        print('制作:'.rjust(5) + self.manufacturer)
        print('发行:'.rjust(5) + self.publisher)
        print('系列:'.rjust(5) + self.series)
        print('类别:'.rjust(5) + self.category)
        print('女优:'.rjust(5) + self.actors)
        print('收藏:'.rjust(5) + self.favor)
        print('预览:'.rjust(5) + self.coverlink)
        print('磁链:'.rjust(5) + self.link)
        '''
        return '番号:'.rjust(5) + self.code + '\n' + '标题:'.rjust(5) + self.title + '\n' + '日期:'.rjust(5) + self.issuedate + '\n' + '时长:'.rjust(5) + self.length + '\n' + '修正:'.rjust(5) + self.mosaic + '\n' + '导演:'.rjust(5) + self.director + '\n' + '制作:'.rjust(5) + self.manufacturer + '\n' + '发行:'.rjust(5) + self.publisher + '\n' + '系列:'.rjust(5) + self.series + '\n' + '类别:'.rjust(5) + self.category + '\n' + '女优:'.rjust(5) + self.actors + '\n' + '收藏:'.rjust(5) + self.favor + '\n' + '预览:'.rjust(5) + self.coverlink + '\n' + '磁链:'.rjust(5) + self.link

    __repr__ = __str__

######################################### DB START#########################################


def dict_factory(cursor, row):
    d = {}
    for idx, col in enumerate(cursor.description):
        d[col[0]] = row[idx]
    return d


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
        return r
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))
    return None


def fetchone(conn, sql, data):
    if sql is not None and sql != '':
        if data is not None:
            # Do this instead
            d = (data,)
            cu = get_cursor(conn)
            # print('执行sql:[{}],参数:[{}]'.format(sql, data))
            cu.execute(sql, d)
            r = cu.fetchall()
            return r
        else:
            logging.error('the [{}] equal None!'.format(data))
    else:
        logging.error('the [{}] is empty or equal None!'.format(sql))
    return None


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

        self.offset = 0
        # self.count = 0
        # self.codes = []
        self.avs = []
        self.conn = None

        self.preButton = QPushButton('上一个')
        self.preButton.clicked.connect(self.itemTurn)
        self.nextButton = QPushButton('下一个')
        self.nextButton.clicked.connect(self.itemTurn)
        self.dbButton = QPushButton('载入数据库')
        self.dbButton.setFixedWidth(100)
        self.dbButton.clicked.connect(self.showDialog)
        self.favorButton = QRadioButton('收藏')
        self.favorButton.setToolTip('收藏/取消收藏当前项')
        self.favorButton.clicked.connect(self.flagItem)
        self.searchButton = QPushButton('搜索')
        self.searchButton.setFixedWidth(100)
        self.searchButton.clicked.connect(self.searchItem)

        self.favorCheck = QCheckBox(self)
        self.favorCheck.setToolTip('是否只浏览已收藏项')

        self.viewModeBox = QComboBox()
        self.viewModeBox.addItems(['顺序', '随机'])
        self.viewModeBox.setToolTip('浏览模式')
        self.viewModeBox.setCurrentIndex(1)
        self.searchBox = QComboBox()
        self.searchBox.currentTextChanged.connect(self.showSelected)
        # self.searchBox.setFixedWidth(300)

        self.dbEdit = QLineEdit()
        self.dbEdit.setFixedWidth(700)
        self.dbEdit.setReadOnly(True)
        self.dbEdit.setAcceptDrops(True)
        self.searchEdit = QLineEdit()
        self.searchEdit.setFixedWidth(700)

        self.infoEdit = QTextEdit()
        self.infoEdit.setReadOnly(True)
        # self.infoEdit.setFixedWidth(300)
        # self.infoEdit.append('addddddddddddaa')

        self.pixmap = QPixmap()
        self.picLabel = QLabel()
        self.picLabel.setPixmap(self.pixmap)
        self.picLabel.mousePressEvent = self.saveImg
        # self.picLabel.setPixmap(self.pixmap.scaled(700, 400))

        self.hboxhead = QHBoxLayout()
        self.hboxhead.addWidget(self.dbEdit)
        self.hboxhead.addWidget(self.dbButton)
        self.hboxhead.addWidget(self.favorButton)
        self.hboxhead.addWidget(self.favorCheck)
        self.hboxhead.addWidget(self.viewModeBox)
        self.hboxhead.addWidget(self.preButton)
        self.hboxhead.addWidget(self.nextButton)

        self.hboxinfo = QHBoxLayout()
        self.hboxinfo.addWidget(self.picLabel)
        self.hboxinfo.addWidget(self.infoEdit)

        self.hboxtail = QHBoxLayout()
        self.hboxtail.addWidget(self.searchEdit)
        self.hboxtail.addWidget(self.searchButton)
        self.hboxtail.addWidget(self.searchBox)

        self.vbox = QVBoxLayout()
        self.vbox.addLayout(self.hboxhead)
        self.vbox.addLayout(self.hboxinfo)
        self.vbox.addLayout(self.hboxtail)

        self.setLayout(self.vbox)

        self.setGeometry(300, 300, 1000, 600)
        self.setWindowTitle('浏览')
        self.show()

    def saveImg(self, event):
        try:
            if event.button() == 2 and self.picLabel.toolTip() is not None and str(self.picLabel.toolTip()).strip() != '':
                fname = QFileDialog.getSaveFileName(self, self.tr('选择图片保存位置'), self.tr(str(self.picLabel.toolTip()).strip() + '.jpg'))[0]
                self.pixmap.toImage().save(fname)
        except Exception as ex:
            print(str(ex))

    def showInfo(self, cav):
        if cav is not None:
            self.infoEdit.setText('')
            self.infoEdit.append('番号:'.center(5) + cav.code)
            self.infoEdit.append('标题:'.center(5) + cav.title)
            self.infoEdit.append('日期:'.center(5) + cav.issuedate)
            self.infoEdit.append('时长:'.center(5) + cav.length)
            self.infoEdit.append('修正:'.center(5) + cav.mosaic)
            self.infoEdit.append('导演:'.center(5) + cav.director)
            self.infoEdit.append('制作:'.center(5) + cav.manufacturer)
            self.infoEdit.append('发行:'.center(5) + cav.publisher)
            self.infoEdit.append('系列:'.center(5) + cav.series)
            self.infoEdit.append('类别:'.center(5) + cav.category)
            self.infoEdit.append('女优:'.center(5) + cav.actors)
            self.infoEdit.append('收藏:'.center(5) + cav.favor)
            self.infoEdit.append('预览:'.center(5) + cav.coverlink)
            self.infoEdit.append('磁链:'.center(5) + cav.link)
            self.pixmap.loadFromData(cav.cover)
            self.picLabel.setPixmap(self.pixmap)
            # self.picLabel.setPixmap(self.pixmap.scaled(700, 400))
            self.setWindowTitle('浏览 - ' + '第' + str(self.avs.index((cav.code, cav.favor)) + 1) + '/' + str(len(self.avs)) + '条')
            # self.setGeometry(300, 300, 1000, 600)
            self.setFixedHeight(600)
            self.picLabel.setToolTip(cav.code + ' ' + cav.title)
            self.offset = self.avs.index((cav.code, cav.favor))
            if int(cav.favor) == 0:
                self.favorButton.setChecked(False)
            else:
                self.favorButton.setChecked(True)

    def getreqCode(self):
        code = self.avs[self.offset][0]
        if self.viewModeBox.currentText() == '随机':
            if self.favorCheck.isChecked():
                if len([item[0] for item in self.avs if int(item[1]) > 0]) == 0:
                    if self.picLabel.toolTip() is not None and str(self.picLabel.toolTip()).strip() != '':
                        code = str(self.picLabel.toolTip()).strip()
                else:
                    code = choice([item[0] for item in self.avs if int(item[1]) > 0])
            else:
                code = choice([item[0] for item in self.avs])
        else:
            if self.picLabel.toolTip() is not None and str(self.picLabel.toolTip()).strip() != '':
                preindex = [item[0] for item in self.avs].index(str(self.picLabel.toolTip()))
                if self.favorCheck.isChecked():
                    if len([item[0] for item in self.avs if int(item[1]) > 0]) == 0:
                        code = str(self.picLabel.toolTip()).strip()
                    else:
                        if (len(self.avs) + self.offset - preindex) % len(self.avs) == 1:
                            for i in range(self.offset, len(self.avs) + self.offset):
                                if int(self.avs[i - len(self.avs)][1]) > 0:
                                    code = self.avs[i - len(self.avs)][0]
                                    break
                        else:
                            for i in range(self.offset, -len(self.avs) + 1, -1):
                                if int(self.avs[i][1]) > 0:
                                    code = self.avs[i][0]
                                    break
        return code

    def fetchInfo(self):
        cav = None
        sql = 'SELECT * FROM av WHERE code = ?'
        for item in fetchone(self.conn, sql, self.getreqCode()):
            cav = av(item['code'], item['title'], item['issuedate'], item['length'], item['mosaic'], item['director'], item['manufacturer'], item['publisher'], item['series'], item['category'], item['actors'], item['favor'], item['coverlink'], item['cover'], item['link'])
        return cav

    def itemTurn(self):
        try:
            sender = self.sender()
            if sender.text() == '上一个':
                self.offset -= 1
            if sender.text() == '下一个':
                self.offset += 1
            self.offset = (len(self.avs) + self.offset) % len(self.avs)
            self.showInfo(self.fetchInfo())
        except Exception as ex:
            print('itemTurn:' + str(ex))

    def showDialog(self):
        fname = QFileDialog.getOpenFileName(self, self.tr('选择数据库'), os.path.join(os.path.expanduser("~"), 'Desktop'), self.tr('DBFile(*.db)'))
        if fname[0]:
            self.dbEdit.setText(fname[0])
            try:
                self.offset = 0
                self.conn = get_conn(fname[0])
                sql = 'SELECT code, favor FROM av order by rowid'
                res = fetchall(self.conn, sql)
                self.avs = [(str(item[0]), str(item[1])) for item in res]
                self.conn.row_factory = dict_factory
                self.showInfo(self.fetchInfo())
            except Exception as ex:
                print('showDialog:' + str(ex))

    def flagItem(self):
        try:
            code = str(self.picLabel.toolTip()).split(' ')[0].strip()
            infoText = self.infoEdit.toPlainText()
            if self.favorButton.isChecked():
                favorState = 1
            else:
                favorState = 0
            sql = 'UPDATE av SET favor = ? WHERE code = ?'
            update(self.conn, sql, [(favorState, code)])
            for i in range(0, len(self.avs)):
                if self.avs[i][0] == code:
                    if favorState == 1:
                        self.avs[i] = (code, '1')
                        infoText = infoText.replace('收藏: 0', '收藏: 1')
                    else:
                        self.avs[i] = (code, '0')
                        infoText = infoText.replace('收藏: 1', '收藏: 0')
                    break
            self.infoEdit.setText(infoText)
        except Exception as ex:
            print('flagItem:' + str(ex))

    def searchItem(self):
        try:
            if self.conn is not None:
                keywords = str(self.searchEdit.text()).strip()
                sql = 'SELECT code, title FROM av where code like "%' + keywords + '%" or title like "%' + keywords + '%" or issuedate like "%' + keywords + '%" or length like "%' + keywords + '%" or mosaic like "%' + keywords + '%" or director like "%' + keywords + '%" or manufacturer like "%' + keywords + '%" or publisher like "%' + keywords + '%" or series like "%' + keywords + '%" or category like "%' + keywords + '%" or actors like "%' + keywords + '%"'
                res = fetchall(self.conn, sql)
                self.searchBox.clear()
                for item in res:
                    itemData = item['code'] + ' ' + item['title']
                    self.searchBox.addItem(itemData)
        except Exception as ex:
            print('searchItem:' + str(ex))

    def showSelected(self):
        try:
            selectedData = str(self.searchBox.currentText()).strip()
            code = selectedData.split(' ')[0]
            sql = 'SELECT * FROM av where code = ?'
            res = fetchone(self.conn, sql, code)
            if res is not None and len(res) > 0:
                item = res[0]
                cav = av(item['code'], item['title'], item['issuedate'], item['length'], item['mosaic'], item['director'], item['manufacturer'], item['publisher'], item['series'], item['category'], item['actors'], item['favor'], item['coverlink'], item['cover'], item['link'])
                self.showInfo(cav)
        except Exception as ex:
            print('showSelected:' + str(ex))


if __name__ == '__main__':

    app = QApplication(sys.argv)
    window = Window()
    sys.exit(app.exec_())
