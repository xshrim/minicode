#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''
   #文件名：pdfmerge.py
   本脚本用来合并pdf文件，输出的pdf文件按输入的pdf文件名生成书签
   使用示例如下：
   python pdfmerge.py -p "D:\pdf-files" -o "merged-out.pdf" -b True'

   示例说明：
   要合并的pdf文件所在的路径： D:\pdf-files
   合并后的pdf文件的输出文件名：merged-out.pdf
   是否从pdf文件中导入书签的值：True
'''
import os, sys
from PyPDF2 import PdfFileReader, PdfFileWriter, PdfFileMerger


def getfilenames(filepath='', filelist_out=[], file_ext='all'):
    # 遍历filepath下的所有文件，包括子目录下的文件
    for fpath, dirs, fs in os.walk(filepath):
        for f in fs:
            fi_d = os.path.join(fpath, f)
            if file_ext == 'all':
                filelist_out.append(fi_d)
            elif os.path.splitext(fi_d)[1] == file_ext:
                filelist_out.append(fi_d)
            else:
                pass
    return filelist_out


def mergefiles(path, output_filename, import_bookmarks=False):
    # 遍历目录下的所有pdf将其合并输出到一个pdf文件中，输出的pdf文件默认带书签，书签名为之前的文件名
    # 默认情况下原始文件的书签不会导入，使用import_bookmarks=True可以将原文件所带的书签也导入到输出的pdf文件中
    merger = PdfFileMerger()
    filelist = getfilenames(filepath=path, file_ext='.pdf')
    filelist = sorted(filelist)
    if len(filelist) == 0:
        print("当前目录及子目录下不存在pdf文件")
        sys.exit()
    for filename in filelist:
        with open(filename, 'rb') as rf:
            file_rd = PdfFileReader(rf)
            short_filename = os.path.basename(os.path.splitext(filename)[0])
            '''
            bm = short_filename.split('-')[0].strip()
            bm = bm.split(']')[1].strip()
            short_filename = bm
            '''
            if file_rd.isEncrypted == True:
                print('不支持的加密文件：%s' % (filename))
                continue
            merger.append(file_rd, bookmark=short_filename, import_bookmarks=import_bookmarks)
            print('合并文件：%s' % (filename))

    out_filename = os.path.join(os.path.abspath(path), output_filename)
    merger.write(out_filename)
    print('合并后的输出文件：%s' % (out_filename))
    merger.close()


mergefiles("./", "/home/xshrim/merge.pdf", True)
