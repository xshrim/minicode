from pyquery import PyQuery as pq
from openpyxl import Workbook
from selenium import webdriver
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities

class product:
    platform = ''
    link = ''
    category = ''
    brand = ''
    model = ''
    price = ''
    msales = ''

    def __init__(self, platform, link, category, brand, model, price, msales):
        self.platform = platform
        self.link = link
        self.category = category
        self.brand = brand
        self.model = model
        self.price = price
        self.msales = msales

def gets(url, platform, category):
    pds = []
    hds = {'User-Agent':"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:45.0) Gecko/20100101 Firefox/45.0"}
    try:
        data = pq(url=url, headers=hds)
        if platform == '亚马逊':
            content = data('div#mainResults')
            if content is None or str(content).strip() == '':
                content = data('div#atfResults')
            items = content('ul li')
            for item in items.items():
                if item.attr('id') is not None:
                    #print(item.attr('id'))
                    view = item('div.s-item-container')
                    brandinfo = view.children('div').eq(1)
                    priceinfo = view.children('div').eq(2)
                    link = str(brandinfo.children('div').eq(0)('div')('a').attr('href')).strip()
                    model = str(brandinfo.children('div').eq(0)('div')('a').attr('title')).strip()
                    brand = str(brandinfo.children('div').eq(1)('span').eq(1).text()).strip()
                    price = str(priceinfo.children('div').eq(0)('a')('span').eq(1).text()).strip()
                    msales = '-'
                    print(link)
                    print('--------------------------------------------------')
                    print(brand)
                    print(model)
                    print(price)
                    print(msales)
                    print('**************************************************')
                    if model != '' and brand != '' and price != '' and model != '-':
                        pds.append(product(platform, link, category, brand, model, price, msales))
        if platform == '京东商城':
            content = data('div#plist')
            items = content('ul li')
            for item in items.items():
                if item.attr('class') is not None:
                    view = item('div[class="gl-i-wrap j-sku-item"]')
                    brandinfo = view('div.p-name')
                    priceinfo = view('div.p-price')
                    link = 'https:' + str(brandinfo('a').attr('href')).strip()
                    brand = '-'
                    model = str(brandinfo('a')('em').text()).strip()
                    price = '-'
                    msales = '-'
                    print(link)
                    print('--------------------------------------------------')
                    print(brand)
                    print(model)
                    print(price)
                    print(msales)
                    print('**************************************************')
                    if model != '' and brand != '' and price != '' and model != '-':
                        pds.append(product(platform, link, category, brand, model, price, msales))
        if platform == '天猫商城':
            content = data('div#J_ItemList')
            items = content('div.product')
            for item in items.items():
                if item.attr('class') is not None:
                    view = item('div.product-iWrap')
                    brandinfo = view('p.productTitle')
                    priceinfo = view('p.productPrice')
                    saleinfo = view('p.productStatus')
                    brand = '-'
                    link = 'https:' + str(brandinfo('a').attr('href')).strip()
                    model = str(brandinfo('a').attr('title')).strip()
                    price = str(priceinfo('em').attr('title')).strip()
                    msales = str(saleinfo.children('span').eq(0)('em').text()).strip().replace('笔', '')
                    print(link)
                    print('--------------------------------------------------')
                    print(brand)
                    print(model)
                    print(price)
                    print(msales)
                    print('**************************************************')
                    if model != '' and brand != '' and price != '' and model != '-':
                        pds.append(product(platform, link, category, brand, model, price, msales))
        if platform == '家乐福':
            dcap = dict(DesiredCapabilities.PHANTOMJS)
            dcap["phantomjs.page.settings.userAgent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:45.0) Gecko/20100101 Firefox/45.0"
            dcap["phantomjs.page.settings.loadImages"] = False
            driver = webdriver.PhantomJS(desired_capabilities=dcap)
            #driver = webdriver.PhantomJS()
            #driver.set_page_load_timeout(10)
            #driver.set_script_timeout(10)
            driver.get(url)
            #driver.find_element_by_class_name('right-botton')
            #btn.click()
            content = driver.find_element_by_id('pnl_prebuy')
            content = pq(content.get_attribute('innerHTML'))
            driver.quit()
            items = content('dl[class="list-h product-list-item"]')
            for item in items.items():
                view = item('dd')
                brandinfo = view('p.p-name')
                priceinfo = view('div.p-price-layer')
                link = 'http://www.carrefour.cn' + str(brandinfo('a').attr('href')).strip()
                brand = '-'
                model = str(brandinfo('a').text()).strip()
                price = str(priceinfo('span[class="p-price fl"]').text()).strip()
                msales = str(priceinfo('div.price').children('p').eq(0)('span').text()).strip()
                print(link)
                print('--------------------------------------------------')
                print(brand)
                print(model)
                print(price)
                print(msales)
                print('**************************************************')
                if model != '' and brand != '' and price != '' and model != '-':
                    pds.append(product(platform, link, category, brand, model, price, msales))
                
    except Exception as ex:
        print(str(ex))
    return pds

def export(products):
    try:
        filename = '价格统计表.xlsx'
        wb = Workbook(write_only=False)
        ws = wb.create_sheet(title='价格统计')
        ws.append(['品类', '品牌', '型号', '价格', '月销量', '平台', '链接'])
        for prod in products:
            ws.append([prod.category, prod.brand, prod.model, prod.price, prod.msales, prod.platform, prod.link])
        wb.remove_sheet(wb.active)
        wb.save(filename)
    except Exception as ex:
        print(str(ex))
        exit

def products(platform, category):
    objs = []
    for idx in range(1, 50):
        if platform == '亚马逊':
            if category == '手机':
                #手机
                purl = 'https://www.amazon.cn/s/ref=lp_665002051_pg_' + str(idx) + '?rh=n%3A2016116051%2Cn%3A%212016117051%2Cn%3A664978051%2Cn%3A665002051&page=' + str(idx) + '&bbn=664978051&ie=UTF8&qid=1490252939'
                #purl = 'https://www.amazon.cn/s/ref=lp_665002051_pg_' + str(idx) + '?rh=n%3A2016116051%2Cn%3A%212016117051%2Cn%3A664978051%2Cn%3A665002051&page=' + str(idx) + '&ie=UTF8&qid=1490257884&lo=communications'
            if category == '笔记本':
                #笔记本
                purl = 'https://www.amazon.cn/s/ref=lp_888483051_pg_' + str(idx) + '?rh=n%3A42689071%2Cn%3A%2142690071%2Cn%3A106200071%2Cn%3A888483051&page='+ str(idx) + '&ie=UTF8&qid=1490258252&spIA=B01N4QQ5SL,B01EJI48BU,B01N4QQ4EB'
            if category == '洗衣机':
                #洗衣机
                purl = 'https://www.amazon.cn/s/ref=lp_2121147051_pg_' + str(idx) + '?rh=n%3A80207071%2Cn%3A%2180208071%2Cn%3A2121147051&page=' + str(idx) + '&ie=UTF8&qid=1490259054&spIA=B01EFM8HVM,B01KXCR4KS,B01KXCVCOC'
        if platform == '京东商城':
            if category == '啤酒':
                #啤酒
                purl= 'https://list.jd.com/list.html?cat=12259,12260,9439&page=' + str(idx) + '&sort=sort_rank_asc&trans=1&JL=6_0_0#J_main'
        if platform == '天猫商城':
            if category == '啤酒':
                purl = 'https://list.tmall.com/search_product.htm?s=' + str((idx-1)*60) + '&q=%C6%A1%BE%C6&sort=s&style=g&from=.list.pc_1_searchbutton&spm=a220m.1000858.a2227oh.d100&type=pc#J_Filter'
            if category == '零食':
                purl = 'https://list.tmall.com/search_product.htm?s=' + str((idx-1)*60) + '&q=%C1%E3%CA%B3&sort=s&style=g&from=mallfp..pc_1_searchbutton&spm=875.7789098.a2227oh.d100&type=pc#J_Filter'
            if category == '进口零食':
                purl = 'https://list.tmall.com/search_product.htm?s=' + str((idx-1)*60) + '&q=%BD%F8%BF%DA%C1%E3%CA%B3&sort=s&style=g&from=mallfp..pc_1_searchbutton&spm=875.7789098.a2227oh.d100&type=pc#J_Filter'
        if platform == '家乐福':
            if category == '饼干糕点':
                purl = 'http://www.carrefour.cn/category?c=250020855&pageNum=' + str(idx)
            if category == '膨化食品':
                purl = 'http://www.carrefour.cn/category?c=250020856&pageNum=' + str(idx)
            if category == '休闲零食':
                purl = 'http://www.carrefour.cn/category?c=250020857&pageNum=' + str(idx)
            if category == '坚果炒货':
                purl = 'http://www.carrefour.cn/category?c=250021199&pageNum=' + str(idx)
            if category == '牛奶果汁':
                purl = 'http://www.carrefour.cn/category?c=250020867&pageNum=' + str(idx)
            if category == '饮料酒水':
                purl = 'http://www.carrefour.cn/category?c=250020798&pageNum=' + str(idx)
        res = gets(purl, platform, category)
        if res is not None and len(res) > 0:
            objs.extend(res)
    return objs

export(products('家乐福', '休闲零食'))

#苏宁易购
#手机
#http://list.suning.com/0-20006-1.html http://list.suning.com/0-20006-90.html

#京东
#啤酒
#https://list.jd.com/list.html?cat=12259,12260,9439&page=2&sort=sort_rank_asc&trans=1&JL=6_0_0#J_main

#家乐福
#酒
#http://www.carrefour.cn/category/?c=250020837-250020835&pageNum=2
        
    
