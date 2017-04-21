from os import path
import requests
import time
from pyquery import PyQuery
from urllib.request import urlretrieve
from openpyxl import Workbook
# sys.stdout = io.TextIOWrapper(sys.stdout.buffer,encoding='gbk')

hds = {'User-Agent': "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:45.0) Gecko/20100101 Firefox/45.0"}


class house:
    hsrc = ''
    htitle = ''
    htime = ''
    hpriceall = ''
    hpriceone = ''
    harea = ''
    htype = ''
    hlayout = ''
    hloc = ''
    hfitment = ''
    hforward = ''
    hestate = ''
    hschool = ''
    hyear = ''
    hfloor = ''
    hcontact = ''
    hphone = ''

    def __init__(self, hsrc, htitle, htime, hpriceall, hpriceone, htype, harea, hlayout, hloc, hestate, hfitment, hforward, hschool, hyear, hfloor, hcontact, hphone):
        self.hsrc = hsrc
        self.htitle = htitle
        self.htime = htime
        self.hpriceall = hpriceall
        self.hpriceone = hpriceone
        self.htype = htype
        self.harea = harea
        self.hlayout = hlayout
        self.hloc = hloc
        self.hfitment = hfitment
        self.hforward = hforward
        self.hestate = hestate
        self.hschool = hschool
        self.hyear = hyear
        self.hfloor = hfloor
        self.hcontact = hcontact
        self.hphone = hphone

    def display(self, dtype):
        if dtype == 'string':
            print('数据来源：' + self.hsrc)
            print('房屋链接：' + self.htitle)
            print('发布时间：' + self.htime)
            print('房屋总价：' + self.hpriceall)
            print('房屋单价：' + self.hpriceone)
            print('房屋户型：' + self.hlayout)
            print('房屋面积：' + self.harea)
            print('区域位置：' + self.hloc)
            print('房屋类型:' + self.htype)
            print('房屋朝向:' + self.hforward)
            print('装修程度：' + self.hfitment)
            print('建成年份：' + self.hyear)
            print('房屋楼层：' + self.hfloor)
            print('所在小区：' + self.hestate)
            print('是否校区：' + self.hschool)
            print('联系人:' + self.hcontact)
            print('联系电话:' + self.hphone)
        if dtype == 'excel':

            pass


def getURLs(src):
    hview = set()
    if src == '58同城':
        for i in range(1, 2):
            purl = 'http://hrb.58.com/ershoufang/pn' + str(i)
            time.sleep(1)
            data = PyQuery(url=purl, headers=hds)
            # print(data)
            links = data('p[class=bthead]')
            for link in links.items():
                print(link.children('a').attr('href'))
                hview.add(link.children('a').attr('href'))
                if len(hview) > 3:
                    break
    if src == '安居客':
        for i in range(1, 2):
            purl = 'http://heb.anjuke.com/sale/p' + str(i)
            time.sleep(1)
            data = PyQuery(url=purl, headers=hds)
            # print(data)
            links = data('div[class=house-title]')
            for link in links.items():
                print(link.children('a').attr('href'))
                hview.add(link.children('a').attr('href'))
                if len(hview) > 3:
                    break
    return hview

def getInfos(site):
    abingo = 1
    bbingo = 1
    acount = 0
    bcount = 0
    houses = []
    if site == 'a' or site == 'A':
        for hurl in getURLs('58同城'):
            acount += 1
            try:
                print('* ' + str(abingo) + '/' + str(acount) + ' *')
                print('房屋链接:' + hurl)
                hsrc = '58同城'.strip()
                time.sleep(1)
                data = PyQuery(url=hurl, headers=hds)
                # print(hurl)
                htitle = data('div[class=bigtitle]').children('h1').text().replace('（', '(').replace('）', ')').strip()
                htime = data('li[class=time]').text().replace('（', '(').replace('）', ')').strip()
                hphone = data('#t_phone').html().replace(' ', '').replace('\n', '').replace('\t', '').strip()

                hinfoa = data('ul[class=suUl]')
                hprice = hinfoa('li').eq(0)('div[class=su_con]').text().replace('&nbsp;', '').replace(' ', '').replace('\n', '').replace('\t', '').replace('（', '(').replace('）', ')').strip()
                hlayout = hinfoa('li').eq(3)('div[class=su_con]').text().replace('&nbsp;', '').replace(' ', '').replace('\n', '').replace('\t', '').replace('（', '(').replace('）', ')').strip()
                if '位置：' in hlayout:
                    hlayout = hlayout[:hlayout.index('位置：')]
                hloc = hinfoa('li').eq(4)('a').eq(0).text().replace('（', '(').replace('）', ')').strip()
                hestate = hinfoa('li').eq(4)('a').eq(2).text().replace('（', '(').replace('）', ')').strip()
                # hreg = hinfoa('li').eq(4)('a').text().replace('（','(').replace('）',')').strip()

                hcontact = hinfoa('li').eq(6)('div[class=su_con]')('span')('a').text().replace(' ', '').strip()

                hinfob = data('ul[class=des_table]')
                htype = hinfob('li[class="des_tablerows clearfix"]').eq(0)('ul')('li').eq(1).text().replace('（', '(').replace('）', ')').strip()
                hfitment = hinfob('li[class="des_tablerows clearfix"]').eq(0)('ul')('li').eq(3).text().replace('（', '(').replace('）', ')').strip()
                htype = hinfob('li[class="des_tablerows clearfix"]').eq(1)('ul')('li').eq(1).text().replace('（', '(').replace('）', ')').strip()
                hyear = hinfob('li[class="des_tablerows clearfix"]').eq(2)('ul')('li').eq(1).text().replace('（', '(').replace('）', ')').strip()
                hfloor = hinfob('li[class="des_tablerows clearfix"]').eq(2)('ul')('li').eq(3).text().replace('（', '(').replace('）', ')').strip()
                hforward = hinfob('li[class="des_tablerows clearfix"]').eq(3)('ul')('li').eq(3).text().replace('（', '(').replace('）', ')').strip()
                hschool = '未知'

                if htitle is not None and htitle.strip() != '':
                    abingo += 1
                    htime = htime                                     # 时间
                    hpriceall = hprice[:hprice.index('万')].strip()     # 总价
                    hpriceone = hprice[hprice.index('(')+1:hprice.index('元')].strip()      # 单价
                    htype = htype                                 # 类型
                    harea = hlayout.split(' ')[-1]
                    harea = harea[:harea.index('㎡')].strip()     # 面积
                    hlayout = hlayout.split(' ')[0].strip()       # 户型
                    hloc = hloc                                   # 区域
                    hestate = hestate                             # 小区
                    hfitment = hfitment.replace(' ', '').replace('找装修', '').strip()   # 装修
                    hforward = hforward                           # 朝向
                    hyear = hyear                                 # 年份
                    hschool = hschool                             # 校区
                    hcontact = hcontact                           # 联系人
                    hphone = hphone                               # 联系电话

                    houseitem = house(hsrc, htitle, htime, hpriceall, hpriceone, htype, harea, hlayout, hloc, hestate, hfitment, hforward, hschool, hyear, hfloor, hcontact, hphone)
                    houseitem.display('string')
                    print('=================================================')
                    houses.append(houseitem)
            except Exception as ex:
                print(str(ex))
                # return []
    if site == 'b' or site == 'B':
        for hurl in getURLs('安居客'):
            bcount += 1
            try:
                print('* ' + str(bbingo) + '/' + str(bcount) + ' *')
                print('房屋链接:' + hurl)
                hsrc = '安居客'.strip()
                time.sleep(1)
                data = PyQuery(url=hurl, headers=hds)
                # print(data)

                htitle = data('div[class=wrapper]').children('h3').text().replace('（', '(').replace('）', ')').strip()
                htime = data('h4[class="block-title houseInfo-title"]').children('span').text().replace('（', '(').replace('）', ')').strip()

                hcontact = data('div[class=broker-wrap]')('p[class=broker-name]').text().replace(' ', '').replace('\n', '').replace('\t', '').strip()
                hphone = data('div[class=broker-wrap]')('p[class=broker-mobile]').text().replace(' ', '').replace('\n', '').replace('\t', '').replace('\"', '').replace('', '').strip()

                hprice = data('span[class="light info-tag"]').children('em').text().replace('（', '(').replace('）', ')').strip()

                hinfo = data('div[class="houseInfo-detail clearfix"]')
                if hinfo is None or hinfo.text().strip() == '':
                    hinfo = data('div[class="houseInfoV2-detail clearfix"]')

                hestate = hinfo('div[class="first-col detail-col"]').children('dl').eq(0).children('dd').children('a').eq(0).text().replace('（', '(').replace('）', ')').strip()
                hloc = hinfo('div[class="first-col detail-col"]').children('dl').eq(1).children('dd').children('p').children('a').eq(0).text().replace('（', '(').replace('）', ')').strip()
                hyear = hinfo('div[class="first-col detail-col"]').children('dl').eq(2).children('dd').text().replace('（', '(').replace('）', ')').strip()
                htype = hinfo('div[class="first-col detail-col"]').children('dl').eq(3).children('dd').text().replace('（', '(').replace('）', ')').strip()

                hlayout = hinfo('div[class="second-col detail-col"]').children('dl').eq(0).children('dd').text().replace('&nbsp;', '').replace(' ', '').replace('\n', '').replace('\t', '').replace('（', '(').replace('）', ')').strip()
                harea = hinfo('div[class="second-col detail-col"]').children('dl').eq(1).children('dd').text().replace('（', '(').replace('）', ')').strip()
                hforward = hinfo('div[class="second-col detail-col"]').children('dl').eq(2).children('dd').text().replace('（', '(').replace('）', ')').strip()
                hfloor = hinfo('div[class="second-col detail-col"]').children('dl').eq(3).children('dd').text().replace('&nbsp;', '').replace(' ', '').replace('\n', '').replace('\t', '').replace('（', '(').replace('）', ')').strip()

                hfitment = hinfo('div[class="third-col detail-col"]').children('dl').eq(0).children('dd').text().replace('（', '(').replace('）', ')').strip()
                hpriceone = hinfo('div[class="third-col detail-col"]').children('dl').eq(1).children('dd').text().replace('（', '(').replace('）', ')').strip()

                hschool = '未知'

                if htitle is not None and htitle.strip() != '':
                    bbingo += 1
                    htime = htime[htime.index('发布时间：') + 5:].strip()  # 时间
                    hpriceall = hprice                              # 总价
                    hpriceone = hpriceone[:hpriceone.index('元')].strip()  # 单价
                    htype = htype                                   # 类型
                    harea = harea.replace('平方米', '').strip()      # 面积
                    hlayout = hlayout                               # 户型
                    hloc = hloc                                     # 区域
                    hestate = hestate                               # 小区
                    hfitment = hfitment                             # 装修
                    hforward = hforward                             # 朝向
                    hyear = hyear.replace('年', '').strip()         # 年份
                    hfloor = hfloor                                 # 楼层
                    hschool = hschool                               # 校区
                    hcontact = hcontact                             # 联系人
                    hphone = hphone                                 # 联系电话

                    houseitem = house(hsrc, htitle, htime, hpriceall, hpriceone, htype, harea, hlayout, hloc, hestate, hfitment, hforward, hschool, hyear, hfloor, hcontact, hphone)
                    houseitem.display('string')
                    print('=================================================')
                    houses.append(houseitem)
            except Exception as ex:
                print(str(ex))
                # return []
    return houses


def export(houses, site):
    try:
        if site == 'a' or site == 'A':
            filename = '哈尔滨58同城二手房信息.xlsx'
        if site == 'b' or site == 'B':
            filename = '哈尔滨安居客二手房信息.xlsx'
        wb = Workbook(write_only=False)
        ws = wb.create_sheet(title='二手房统计')
        ws.append(['区域位置', '面积', '总价(万)', '单价(元/㎡)', '发布时间', '建成年份', '房屋户型', '所在小区', '房屋类型', '是否校区', '装修程度', '房屋楼层', '房屋朝向', '联系人', '联系电话', '数据来源'])
        for house in houses:
            ws.append([house.hloc, house.harea, house.hpriceall, house.hpriceone, house.htime, house.hyear, house.hlayout, house.hestate, house.htype, house.hschool, house.hfitment, house.hfloor, house.hforward, house.hcontact, house.hphone, house.hsrc])
        wb.remove_sheet(wb.active)
        wb.save(filename)
    except Exception as ex:
        print(str(ex))
        exit


while True:
    srcsite = input('Please choose the propery listings:\nA.58同城\nB.安居客\nYour choice:')
    if srcsite == 'a' or srcsite == 'A' or srcsite == 'A' or srcsite == 'B':
        houses = getInfos(srcsite)
        export(houses, srcsite)
        break
    elif srcsite == 'ab' or srcsite == 'aB' or srcsite == 'Ab' or srcsite == 'AB':
        houses = getInfos('a')
        export(houses, 'a')
        houses = getInfos('b')
        export(houses, 'b')
        break
    else:
        continue
