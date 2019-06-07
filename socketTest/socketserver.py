import socket
import time
import threading

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  #创建socket (AF_INET:IPv4, AF_INET6:IPv6) (SOCK_STREAM:面向流的TCP协议)

s.bind(('127.0.0.1', 10021))                           #绑定本机IP和任意端口(>1024)

s.listen(1)                                            #监听，等待连接的最大数目为1

print('Server is running...')                          
                               
def TCP(sock, addr):                                   #TCP服务器端处理逻辑
    
    print('Accept new connection from %s:%s.' %addr)   #接受新的连接请求

    while True:
        data = sock.recv(1024)                         #接受其数据
        time.sleep(1)                                  #延迟
        if not data or data.decode() == 'quit':        #如果数据为空或者'quit'，则退出
            break
        sock.send(data.decode('utf-8').upper().encode())  #发送变成大写后的数据,需先解码,再按utf-8编码,  encode()其实就是encode('utf-8')

    sock.close()                                       #关闭连接
    print('Connection from %s:%s closed.' %addr)       

while True:
    
    sock, addr = s.accept()                            #接收一个新连接
    TCP(sock, addr)                                    #处理连接