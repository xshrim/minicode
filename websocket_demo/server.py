import random, time
from websocket_server import WebsocketServer


# 当新的客户端连接时会提示
# Called for every client connecting (after handshake)
def new_client(client, server):
    print("New client connected and was given id %d" % client['id'])
    server.send_message_to_all("Hey all, a new client has joined us")
    while True:
        server.send_message(client, str(random.randint(0, 1000)))
        time.sleep(1)


# 当旧的客户端离开
# Called for every client disconnecting
def client_left(client, server):
    print("Client(%d) disconnected" % client['id'])


# 接收客户端的信息。
# Called when a client sends a message
def message_received(client, server, message):
    if len(message) > 200:
        message = message[:200] + '..'
    print("Client(%d) said: %s" % (client['id'], message))


PORT = 9001
server = WebsocketServer(PORT, "0.0.0.0")
server.set_fn_new_client(new_client)
server.set_fn_client_left(client_left)
server.set_fn_message_received(message_received)
server.run_forever()