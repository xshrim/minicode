import pymysql
from flask import Flask, jsonify, render_template, request
from flask_socketio import SocketIO, emit

app = Flask(__name__)
app.debug = True
app.config['SECRET_KEY'] = 'secret'
socketio = SocketIO(app)


@app.route('/fetch')
def fetch():
    symbol = request.args.get('symbol')
    db = pymysql.connect("localhost", "root", "root", "school")
    cursor = db.cursor()
    cursor.execute("select * from student where `name` = '" + symbol + "'")
    data = cursor.fetchone()
    if data is not None:
        return jsonify(sid=data[0], name=data[1], score=data[2])
    else:
        return jsonify(sid="null", name="null", score="null")


@app.route('/')
def index():
    return render_template('index.html')


@socketio.on('my event', namespace='/test')
def test_message(message):
    emit('my response', {'data': message['data']})


@socketio.on('my broadcast event', namespace='/test')
def test_broadcast_message(message):
    emit('my response', {'data': message['data']}, broadcast=True)


@socketio.on('connect', namespace='/test')
def test_connect():
    emit('my response', {'data': 'Connected'})


@socketio.on('disconnect', namespace='/test')
def test_disconnect():
    print('Client disconnected')


if __name__ == '__main__':
    # app.run(debug=True)
    socketio.run(app, host='0.0.0.0', port=5100)