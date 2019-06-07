#!/usr/bin/env python
# coding=utf-8

import os
import re
import ssl
import sys
import time
import json
import requests
import threading
from PIL import Image


def showQrcode(QRImagePath):
    '''
    plt.ion()
    plt.imshow(plt.imread(QRImagePath))
    plt.axis('off')
    plt.show()
    plt.pause(10)
    plt.close()
    '''
    Image.open(QRImagePath).show()


def keepLogin(mySession, session_id):
    while True:
        url = 'http://im37.115.com/chat/r'
        params = {
            'VER': '2',
            'c': 'b0',
            's': session_id,
            '_t': str(int(time.time() * 1000)),
        }
        r = mySession.get(url=url, params=params)
        time.sleep(60)


def login():
    try:
        # 构建会话
        ssl._create_default_https_context = ssl._create_unverified_context
        headers = {'User-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2663.0 Safari/537.36'}
        mySession = requests.Session()
        mySession.headers.update(headers)

        # 获取登陆信息
        url = 'http://passport.115.com'
        params = {
            'ct': 'login',
            'ac': 'qrcode_token',
            'is_ssl': '1',
        }

        r = mySession.get(url=url, params=params)
        r.encoding = 'utf-8'
        data = r.text
        try:
            uid = json.loads(data)['uid']
            uidTime = json.loads(data)['time']
            sign = json.loads(data)['sign']
        except Exception as e:
            return None

        url = 'http://msg.115.com/proapi/anonymous.php'
        params = {
            'ac': 'signin',
            'user_id': uid,
            'sign': sign,
            'time': uidTime,
            '_': str(int(time.time() * 1000)),
        }
        r = mySession.get(url=url, params=params)
        session_id = json.loads(r.text)['session_id']
        if session_id is None:
            return None

        # 获取登陆二维码
        # global QrcodeUrl
        url = 'http://qrcode.115.com/api/qrcode.php'
        params = {
            'qrfrom': '1',
            'uid': uid,
            '_' + str(uidTime): '',
            '_t': str(int(time.time() * 1000)),
        }

        r = mySession.get(url=url, params=params)
        # QrcodeUrl = r.url
        r.encoding = 'utf-8'
        QRImagePath = os.path.join(os.getcwd(), 'qrcode.jpg')
        with open(QRImagePath, 'wb') as wbf:
            wbf.write(r.content)
        threading.Thread(target=showQrcode, args=(QRImagePath,)).start()

        # 扫码登陆
        print("使用115手机客户端扫码登录")
        time.sleep(1)

        while True:
            url = 'http://im37.115.com/chat/r'
            params = {
                'VER': '2',
                'c': 'b0',
                's': session_id,
                '_t': str(int(time.time() * 1000)),
            }
            r = mySession.get(url=url, params=params)
            try:
                status = json.loads(r.text)[0]['p'][0]['status']
                if status == 1001:
                    print("请点击登录")
                elif status == 1002:
                    print("登录成功")
                    break
                # else:
                    # print("使用115手机客户端扫码登录")
                    # return
            except Exception as e:
                pass
                # print("使用115手机客户端扫码登录")
                # print("超时，请重试")
                # sys.exit(0)

        # 触发登陆
        url = 'http://passport.115.com/'
        params = {
            'ct': 'login',
            'ac': 'qrcode',
            'key': uid,
            'v': 'android',
            'goto': 'http%3A%2F%2Fwww.J3n5en.com'
        }
        if mySession.get(url=url, params=params) is None:
            return None

        # 心跳保持持续登陆状态
        print('开启心跳线程')
        threading.Thread(target=keepLogin, args=(mySession, session_id,)).start()

        # 获取操作task所需信息
        url = 'http://115.com/'
        params = {
            'ct': 'offline',
            'ac': 'space',
            '_': str(int(time.time() * 1000)),
        }
        r = mySession.get(url=url, params=params)
        tsign = json.loads(r.text)['sign']
        ttime = json.loads(r.text)['time']
        if tsign is None or ttime is None:
            return None

        # 获取登陆用户信息
        url = 'http://passport.115.com/'
        params = {
            'ct': 'ajax',
            'ac': 'islogin',
            'is_ssl': '1',
            '_' + str(int(time.time() * 1000)): '',
        }
        uinfos = json.loads(mySession.get(url=url, params=params).text)
        uname = uinfos['data']['USER_NAME']
        userid = uinfos['data']['USER_ID']

        print("====================")
        print("用户ID：" + str(userid))
        print("用户名：" + str(uinfos['data']['USER_NAME']))
        if uinfos['data']['IS_VIP'] == 1:
            itype = '会员'
            print("身  份：" + "会员")
            url = 'http://115.com/web/lixian/?ct=lixian&ac=task_lists'
            data = {
                'page': '1',
                'uid': userid,
                'sign': tsign,
                'time': ttime,
            }
            quota = json.loads(mySession.post(url=url, data=data).text)['quota']
            total = json.loads(mySession.post(url=url, data=data).text)['total']
            print("月离线：" + str(quota) + "/" + str(total))
        else:
            itype = '非会员'
            quota = 0
            total = 0
            print("身  份：" + "非会员")
            print("月离线：" + str(quota) + "/" + str(total))
        print("===================")

        return {
            'sid': str(session_id),
            'session': mySession,
            'uname': str(uname),
            'uid': str(userid),
            'sign': str(tsign),
            'time': str(ttime),
            'itype': str(itype),
            'quota': str(quota),
            'total': str(total)
        }
    except Exception as e:
        return None


def addLinktask(link, infos):
    url = "http://115.com/web/lixian/?ct=lixian&ac=add_task_url"
    data = {
        'url': link,
        'uid': infos['uid'],
        'sign': infos['sign'],
        'time': infos['time']
    }
    linkinfo = json.loads(infos['session'].post(url, data=data).content)
    try:
        print(linkinfo['error_msg'])
    except Exception as e:
        print(linkinfo['name'])


def addLinktasks(linklist, infos):
    if len(linklist) > 15:
        for i in range(0, len(linklist), 15):
            newlist = linklist[i:i + 15]
            addLinktasks(newlist, infos)
    else:
        url = "http://115.com/web/lixian/?ct=lixian&ac=add_task_urls"
        data = {
            'uid': infos['uid'],
            'sign': infos['sign'],
            'time': infos['time']
        }
        for i in range(len(linklist)):
            data['url[' + str(i) + ']'] = linklist[i]
        linksinfo = json.loads(infos['session'].post(url, data=data).text)
        # print linksinfo['result']
        for linkinfo in linksinfo['result']:
            try:
                print(linkinfo['error_msg'])
            except Exception as e:
                print(linkinfo['name'])


def upload_torrent(filename, filedir, infos):
    # 获取种子上传信息
    url = 'http://115.com/'
    params = {
        'ct': 'lixian',
        'ac': 'get_id',
        'torrent': '1',
        '_': str(int(time.time() * 1000)),
    }
    cid = json.loads(infos['session'].post(url, params=params).text)['cid']
    req = infos['session'].get('http://115.com/?tab=offline&mode=wangpa').content
    reg = re.compile('upload\?(\S+?)"')
    ids = re.findall(reg, str(req))
    upload_url = ids[0]

    # 上传种子
    url = 'http://upload.115.com/upload?' + upload_url
    files = {
        'Filename': ('', 'torrent.torrent', ''),
        'target': ('', 'U_1_' + str(cid), ''),
        'Filedata': ('torrent.torrent', open(filedir, 'rb'), 'application/octet-stream'),
        'Upload': ('', 'Submit Query', ''),
    }
    # infos['session'].get('http://upload.115.com/crossdomain.xml')
    req = infos['session'].post(url=url, files=files)
    req = json.loads(req.content)
    if req['state'] is False:
        print("上传种子出错了1")
        return False
    data = {'file_id': req['data']['file_id']}
    post_url = 'http://115.com/lixian/?ct=lixian&ac=torrent'
    data = {
        'pickcode': req['data']['pick_code'],
        'sha1': req['data']['sha1'],
        'uid': infos['uid'],
        'sign': infos['sign'],
        'time': infos['time'],
    }
    resp = infos['session'].post(url=post_url, data=data)
    resp = json.loads(resp.content)
    if resp['state'] is False:
        print("上传种子出错2")
        return False
    wanted = None
    idx = 0
    for item in resp['torrent_filelist_web']:
        if item['wanted'] != -1:
            if wanted is None:
                wanted = str(idx)
            else:
                wanted = wanted + ',' + str(idx)
        idx += 1
    post_url = 'http://115.com/lixian/?ct=lixian&ac=add_task_bt'
    data = {
        'info_hash': resp['info_hash'],
        'wanted': wanted,
        'savepath': resp['torrent_name'].replace('\'', ''),
        'uid': infos['uid'],
        'sign': infos['sign'],
        'time': infos['time'],
    }
    resp = infos['session'].post(post_url, data).content
    ret = json.loads(str(resp))
    print(ret['name'])
    if 'error_msg' in ret:
        print(ret['error_msg'])


def add_many_bt(fpath, infos):
    for parent, dirnames, filenames in os.walk(fpath):
        for filename in filenames:
            filedir = os.path.join(parent, filename)
            # time.sleep(1)
            upload_torrent(filename, filedir, infos)
            # print open(qq,'rb')


def main():
    infos = login()
    if infos is not None:
        addLinktask("magnet:?xt=urn:btih:A63ABDC8972BA8746C6613549AC3BECB41AB2EC3&dn=MDS-807", infos)
        # addLinktasks([link]) 传入一个list
        # print tsign
        # print "fuck"
        # get_bt_upload_info()
        # upload_torrent()
        # add_many_bt()


if __name__ == '__main__':
    main()
    # print cid
    # print requests.get("http://j3n5en.com", proxies=proxies).text
