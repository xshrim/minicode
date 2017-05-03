#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re
import sys
import getopt
import platform
import subprocess

class cip:
    address = ''
    packetloss = ''
    mindelay = ''
    mindelay = ''
    avgdelay = ''
    output = ''
    status = ''
    def __init__(self, address, packetloss, mindelay, maxdelay, avgdelay, output, status):
        self.address = address
        self.packetloss = packetloss
        self.mindelay = mindelay
        self.maxdelay = maxdelay
        self.avgdelay = avgdelay
        self.output = output
        self.status = status


def checkIP(ip):
    pattern = re.compile(r"^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\." + "(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\." + "(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\." + "(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$")
    if ip is not None and pattern.match(ip) is not None:
        return True
    else:
        return False


def convertIP(ip):
    return '.'.join([x.rjust(3, '0') for x in ip.split('.')])
    '''
    ipa, ipb, ipc, ipd = ip.split('.')
    ipa = ipa.rjust(3, '0')
    ipb = ipb.rjust(3, '0')
    ipc = ipc.rjust(3, '0')
    ipd = ipd.rjust(3, '0')
    return '.'.join([ipa, ipb, ipc, ipd])
    '''


def showIP():
    try:
        syspf = platform.system().lower()
        if syspf == 'linux':
            cmd = 'ifconfig -a && route -n'
        if syspf == 'windows':
            cmd = 'ipconfig /all'
        p = subprocess.Popen(cmd, stdin = subprocess.PIPE, stdout = subprocess.PIPE, stderr = subprocess.PIPE, shell = True)
        output = p.stdout.read()
        return output.decode('gbk')
    except Exception as ex:
        return str(ex)


def getIPS(ipda, ipdb):
    ips = []
    starta, startb, startc, startd = ['', '', '', '']
    enda, endb, endc, endd = ['', '', '', '']
    ipda = convertIP(ipda)
    ipdb = convertIP(ipdb)
    if ipda > ipdb:
        ipda, ipdb = ipdb, ipda
    starta, startb, startc, startd = ipda.split('.')
    enda, endb, endc, endd = ipdb.split('.')
    for ipa in range(int(starta), int(enda) +1):
        startrange = 0
        endrange = 255
        if ipa == int(starta):
            startrange = int(startb)
        if ipa == int(enda):
            endrange = int(endb)
        for ipb in range(startrange, endrange + 1):
            startrange = 0
            endrange = 255
            if ipa == int(starta) and ipb == int(startb):
                startrange= int(startc)
            if ipa == int(enda) and ipb == int(endb):
                endrange = int(endc)
            for ipc in range(startrange, endrange + 1):
                for ipd in range(1, 255):
                    newip = '.'.join([str(ipa), str(ipb), str(ipc), str(ipd)])
                    if ipda <= convertIP(newip) <= ipdb:
                        #print(newip)
                        ips.append(newip)
    return ips


def getPING(domain, times=2, timeout=200):
    try:
        ipaddr = ''
        lost = ''
        minimum = ''
        maximum = ''
        average = ''
        syspf = platform.system().lower()
        if syspf == 'linux':
            cmd = 'ping -c ' + str(times) + ' -w ' + str(timeout) + ' ' + str(domain)
        if syspf == 'windows':
            cmd = 'ping.exe -n ' + str(times) + ' -w '+ str(timeout) + ' ' + str(domain)
        p = subprocess.Popen(cmd, stdin = subprocess.PIPE, stdout = subprocess.PIPE, stderr = subprocess.PIPE, shell = True)
        output = p.stdout.read()
        if syspf == 'windows':
            output = output.decode('gbk')
            regIP = r'\d+\.\d+\.\d+\.\d+'   ## Pinging www.a.shifen.com [115.239.211.112] with 32 bytes of data
            regLost = r'(.*?)(\d+%)(.*)' ## Packets: Sent = 4, Received = 4, Lost = 0 (0% loss)   数据包: 已发送 = 4，已接收 = 4，丢失 = 0 (0% 丢失)，
            regMinimum = r'Minimum = \d+ms|最短 = \d+ms'## Minimum = 37ms, Maximum = 38ms, Average = 37ms   最短 = 37ms，最长 = 77ms，平均 = 48ms
            regMaximum = r'Maximum = \d+ms|最长 = \d+ms'
            regAverage = r'Average = \d+ms|平均 = \d+ms'
            ipaddr = re.search(regIP, output)
            lost = re.search(regLost, output)
            minimum = re.search(regMinimum, output)
            maximum = re.search(regMaximum, output)
            average = re.search(regAverage, output)
            if ipaddr is not None:
                ipaddr = ipaddr.group()
            if lost is not None:
                lost = lost.group(2)
            if minimum is not None:
                minimum = re.search(r'(.*?)(\d+)(.*)', minimum.group()).group(2)
                #minimum = ''.join(filter(lambda x:x.isdigit(),minimum.group()))
            if maximum is not None:
                maximum = re.search(r'(.*?)(\d+)(.*)', maximum.group()).group(2)
                #maximum = ''.join(filter(lambda x:x.isdigit(),maximum.group()))
            if average is not None:
                average = re.search(r'(.*?)(\d+)(.*)', average.group()).group(2)
                #average = ''.join(filter(lambda x:x.isdigit(),average.group()))
        if syspf == 'linux':
            output = output.decode('gbk')
            regIP = r'\d+\.\d+\.\d+\.\d+'
            regLost = r'(.*?)(\d+%)(.*)'
            regmum = r'(min/avg/max/mdev = )(.*)ms'
            #print(output)
            ipaddr = re.search(regIP, output)
            lost = re.search(regLost, output)
            mum = re.search(regmum, output)
            if ipaddr is not None:
                ipaddr = ipaddr.group()
            if lost is not None:
                lost = lost.group(2)
            if mum is not None:
                mums = str(mum.group(2)).replace(' ', '')
                minimum, maximum, average = mums.split('/')[:3]
        #print(average.group())
        return [str(ipaddr), str(lost), str(minimum), str(maximum), str(average), str(output)]
    except Exception as ex:
        #print(ex)
        return [str(domain), '100%', 'None', 'None', 'None', str(ex)]


def testIP(iplist, times=2, timeout=200, rtop=True):
    ips = []
    if isinstance(iplist, (tuple, list)):
        for ipaddr in iplist:
            if True or checkIP(ipaddr):
                if rtop:
                    print('Testing ' + ipaddr, end=' ... ')
                res = getPING(ipaddr, times, timeout)
                #print(res)
                if res[1] == 'None':
                    if rtop:
                        print('[NotFound]')
                    ips.append(cip(res[0], res[1], res[2], res[3], res[4], res[5], 'NotFound'))
                elif res[1] == '100%':
                    if rtop:
                        print('[Unreachable]')
                    ips.append(cip(res[0], res[1], res[2], res[3], res[4], res[5], 'Unreachable'))
                else:
                    if rtop:
                        print('[Reachable]')
                    ips.append(cip(res[0], res[1], res[2], res[3], res[4], res[5], 'Reachable'))
            else:
                print('The ip address format is illegal!')
    else:
        print('The argument is illegal!')
    return ips


def testIPS(startip, endip, times=2, timeout=200, rtop=True):
    ips = []
    if checkIP(startip) and checkIP(endip):
        for ipaddr in getIPS(startip, endip):
            if rtop:
                print('Testing ' + ipaddr, end=' ... ')
            res = getPING(ipaddr, times, timeout)
            #print(res)
            if res[1] == 'None':
                if rtop:
                    print('[NotFound]')
                ips.append(cip(res[0], res[1], res[2], res[3], res[4], res[5], 'NotFound'))
            elif res[1] == '100%':
                if rtop:
                    print('[Unreachable]')
                ips.append(cip(res[0], res[1], res[2], res[3], res[4], res[5], 'Unreachable'))
            else:
                if rtop:
                    print('[Reachable]')
                ips.append(cip(res[0], res[1], res[2], res[3], res[4], res[5], 'Reachable'))
    else:
        print('The ip address format is illegal!')
    return ips


def main(argv):
    startip = ''
    endip = ''
    if len(argv) == 0:
        print(showIP())
    else:
        try:
            opts, args = getopt.getopt(argv, "hf:t:", ["from=", "to="])
        except getopt.GetoptError:
            print ('Usage  : ipp.py [-f <startip>] [-t <endip>] [<ipaddr>]')
            exit(2)
        if len(args) > 0:
            testIP(args)
        for opt, arg in opts:
            if opt == '-h':
                print ('Usage  : ipp.py [-f <startip>] [-t <endip>] [<ipaddr>]')
                exit()
            elif opt in ("-f", "--from"):
                startip = arg
            elif opt in ("-t", "--to"):
                endip = arg
        if startip != '' and endip != '':
            testIPS(startip, endip)


if __name__ == '__main__':
    main(sys.argv[1:])



