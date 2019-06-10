# -*- encoding: utf-8 -*-
import json
import requests
import time
from pyquery import PyQuery as pq
from gevent import monkey
from gevent.pool import Pool
import ssl
import re
import f115api


monkey.patch_all()
proxy_pool = []
pool_num = 20
gpool = Pool(pool_num)
base_url = "https://www.seedmm.com/uncensored"
headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36',
        'Referer':'https://www.javbus.me/SDSI-026',
        'Upgrade-Insecure-Requests':'1',
    }


def crawl(page):
    if page != 1:
        url = base_url + "/page/" + str(page)
    else:
        url = base_url
    print url 
    resp = requests.get(url)
    response = pq(url, headers=headers)
    one_page_av = []
    for each in response('#waterfall a').items():
        one_page_av.append(detail_page(each.attr.href))
    # print one_page_av
    f115api.addLinktasks(one_page_av)




def detail_page(link):
    response = requests.get(link, headers=headers).text
    m = re.search('gid = (\d+)',response)
    gid = m.group(0).strip("gid =")
    url = 'https://www.seedmm.com/ajax/uncledatoolsbyajax.php'
    params = {
        'gid': gid,
        'uc': "1",
    }
    r = requests.get(url=url, params=params,headers=headers)
    r.encoding = 'utf-8'
    magnets = pq(r.text)("a[style='color:#333']")
    for magnet in magnets.items():
        return magnet.attr("href")

if __name__ == "__main__":
    ssl._create_default_https_context = ssl._create_unverified_context
    f115api.main()
    links = range(1, 5)
    for link in links:
        gpool.spawn(crawl, link)
    gpool.join()
