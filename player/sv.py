#!/usr/bin/env python3
#-*- coding:utf-8 -*-
import cv2
import threading
import win32gui, win32con


class Producer(threading.Thread):
    """docstring for ClassName"""

    def __init__(self, str_rtsp):
        super(Producer, self).__init__()
        self.str_rtsp = str_rtsp
        self.play = True
        #通过cv2中的类获取视频流操作对象cap
        self.cap = cv2.VideoCapture(self.str_rtsp)
        #调用cv2方法获取cap的视频帧（帧：每秒多少张图片）
        fps = self.cap.get(cv2.CAP_PROP_FPS)
        print(fps)
        #获取cap视频流的每帧大小
        size = (int(self.cap.get(cv2.CAP_PROP_FRAME_WIDTH)), int(self.cap.get(cv2.CAP_PROP_FRAME_HEIGHT)))
        print(size)
        #定义编码格式mpge-4
        fourcc = cv2.VideoWriter_fourcc('M', 'P', '4', '2')
        #定义视频文件输入对象
        self.outVideo = cv2.VideoWriter('saveDir.avi', fourcc, fps, size)
        cv2.namedWindow("cap video", 0)

    def run(self):
        print('in producer')
        while True:
            ret, image = self.cap.read()
            if (ret == True):
                if win32gui.FindWindow(None, 'cap video'):
                    cv2.imshow('cap video', image)
                    self.outVideo.write(image)
                else:
                    self.outVideo.release()
                    self.cap.release()
                    cv2.destroyAllWindows()
                    break
            cv2.waitKey(1)
            if cv2.waitKey(1) & 0xFF == ord('q'):
                self.outVideo.release()
                self.cap.release()
                cv2.destroyAllWindows()
                break
                # continue


if __name__ == '__main__':
    print('run program')
    rtsp_str = 'rtmp://live.hkstv.hk.lxdns.com/live/hks'
    producer = Producer(rtsp_str)
    producer.start()
