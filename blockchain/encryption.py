import cv2
import time
import pickle

from Crypto.Cipher import AES


class Crypt(object):
    def __init__(self, key):
        self.mode = AES.MODE_CBC
        self.key = self.pad_key(key)

    def pad(self, text):
        while len(text) % 16 != 0:
            text += b'\0'
        return text

    def depad(self, text):
        while text[-1] == 0:
            text = text[:-1]
        return text

    def pad_key(self, key):
        key = bytes(key, encoding="utf8")
        while len(key) % 16 != 0:
            key += b'\0'
        return key

    def encrypt(self, text):
        texts = self.pad(text)
        aes = AES.new(self.key, self.mode, self.key)
        res = aes.encrypt(texts)
        return res

    def decrypt(self, text):
        texts = self.pad(text)
        aes = AES.new(self.key, self.mode, self.key)
        res = self.depad(aes.decrypt(texts))
        return res


def encrypt_frame(frame, key):

    print(time.asctime(time.localtime(time.time())))
    enc = Crypt(key).encrypt(frame)
    print(time.asctime(time.localtime(time.time())))

    with open("enc_vid.enc", 'ab') as fo:
        fo.write(enc)
        fo.write('\n'.encode())
        # fo.write("\n")  #Adding delimeter AFTER the encrypted frame


def decryptVideo(srcPath, tgtPath):
    cap = cv2.VideoCapture(srcPath)

    fps = cap.get(cv2.CAP_PROP_FPS)
    size = (int(cap.get(cv2.CAP_PROP_FRAME_WIDTH)), int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT)))

    while True:
        ret, frame = cap.read()

        if frame is None:
            # reach to the end of the video file
            break
        iframe = cv2.bitwise_not(frame)
        # iframe = cv2.bitwise_xor(iframe, frame)

        cv2.imshow("NOT", iframe)
        cv2.waitKey(int(30))
        continue
    cap.release()


def encryptVideo(srcPath, tgtPath):
    cap = cv2.VideoCapture(srcPath)

    fps = cap.get(cv2.CAP_PROP_FPS)
    size = (int(cap.get(cv2.CAP_PROP_FRAME_WIDTH)), int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT)))

    out = cv2.VideoWriter(tgtPath, cv2.VideoWriter.fourcc('I', '4', '2', '0'), fps, size)

    while True:
        ret, frame = cap.read()

        if frame is None:
            # reach to the end of the video file
            break
        iframe = cv2.bitwise_not(frame)
        # iframe = cv2.bitwise_xor(iframe, frame)
        out.write(iframe)

        # cv2.imshow("NOT", frame)
        # cv2.waitKey(int(fps))
        continue
    cap.release()


# encryptVideo('./1.mp4', '1.avi')

# decryptVideo('./1.avi', '1.avi')


def encryptMedia(srcPath, tgtPath):
    cap = cv2.VideoCapture(srcPath)

    fps = cap.get(cv2.CAP_PROP_FPS)
    size = (int(cap.get(cv2.CAP_PROP_FRAME_WIDTH)), int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT)))

    out = cv2.VideoWriter(tgtPath, cv2.VideoWriter.fourcc('I', '4', '2', '0'), fps, size)

    while True:
        ret, frame = cap.read()

        if frame is None:
            # reach to the end of the video file
            break

        # converting numpy array (frame) into str.
        # pickle preserves metadata like shape of array etc.
        iframe = pickle.dumps(frame, protocol=pickle.HIGHEST_PROTOCOL)
        # print(time.asctime(time.localtime(time.time())))
        iframe = Crypt("A").encrypt(iframe)
        # print(time.asctime(time.localtime(time.time())))
        cv2.waitKey(int(1000 / fps))
        out.write(iframe)
        print('*' * 50)
    cap.release()


def enMedia(srcPath, tgtPath):
    with open(srcPath, 'rb') as rf:
        data = rf.read()
        enc = Crypt("A").encrypt(data)
        with open(tgtPath, 'wb') as wf:
            wf.write(enc)
        print(enc)
        dec = Crypt("A").decrypt(enc)
        print(dec)


def deMedia(srcPath, tgtPath):
    with open(srcPath, 'rb') as rf:
        data = rf.read()
        dec = Crypt("A").decrypt(data)
        with open(tgtPath, 'wb') as wf:
            wf.write(dec)


# encryptMedia('./1.mp4', './1.avi')
enMedia('./1.mp4', './1.enc')
deMedia('./1.enc', './2.mp4')
# cleanup the camera and close any open windows
cv2.destroyAllWindows()