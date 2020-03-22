package core

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// copy data from WebSocket to ssh server
// and copy data from ssh server to WebSocket

// write data to WebSocket
// the data comes from ssh server.
type wsBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

// implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

const (
	wsMsgCmd    = "cmd"
	wsMsgResize = "resize"
)

type wsMsg struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

// connect to ssh server using ssh session.
type SshConn struct {
	// calling Write() to write data into ssh server
	StdinPipe io.WriteCloser
	// Write() be called to receive data from ssh server
	ComboOutput *wsBufferWriter
	Session     *ssh.Session
}

//flushComboOutput flush ssh.session combine output into websocket response
func flushComboOutput(w *wsBufferWriter, wsConn *websocket.Conn, count int) error {
	if w.buffer.Len() != 0 {
		if err := wsConn.WriteMessage(websocket.TextMessage, w.buffer.Bytes()); err != nil {
			return err
		}
		w.buffer.Reset()
	} else if count == 0 {
		// 发送心跳包保持websocket连接, 防止代理层(如nginx)超时自动关闭连接
		// (另一种方式是设置nginx的proxy_read_timeout, proxy_connect_timeout, proxy_send_timeout自动断开超时时间)
		if err := wsConn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
			return err
		}
	}
	return nil
}

// setup ssh shell session
// set Session and StdinPipe here,
// and the Session.Stdout and Session.Sdterr are also set.
func NewSshConn(cols, rows int, sshClient *ssh.Client) (*SshConn, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	// we set stdin, then we can write data to ssh server via this stdin.
	// but, as for reading data from ssh server, we can set Session.Stdout and Session.Stderr
	// to receive data from ssh server, and write back to somewhere.
	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(wsBufferWriter)
	//ssh.stdout and stderr will write output into comboWriter
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	// Start remote shell
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	return &SshConn{StdinPipe: stdinP, ComboOutput: comboWriter, Session: sshSession}, nil
}

func (s *SshConn) Close() {
	if s.Session != nil {
		s.Session.Close()
	}

}

//ReceiveWsMsg  receive websocket msg do some handling then write into ssh.session.stdin
func (ssConn *SshConn) ReceiveWsMsg(wsConn *websocket.Conn, logBuff *bytes.Buffer, exitCh chan bool) {
	//tells other go routine quit
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			//read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				logrus.WithError(err).Error("reading webSocket message failed")
				return
			}
			//unmashal bytes into struct
			// 			msgObj := wsMsg{
			// 				Type: "cmd",
			// 				Cmd:  "",
			// 				Rows: 50,
			// 				Cols: 180,
			// 			}
			var msgObj wsMsg
			if err := json.Unmarshal(wsData, &msgObj); err != nil {
				logrus.WithError(err).WithField("wsData", string(wsData)).Error("unmarshal websocket message failed")
			}
			switch msgObj.Type {
			case wsMsgResize:
				//handle xterm.js size change
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err := ssConn.Session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						logrus.WithError(err).Error("ssh pty change windows size failed")
					}
				}
			case wsMsgCmd:
				//handle xterm.js stdin
				// bytes, err := base64.StdEncoding.DecodeString(msgObj.Cmd)
				// ////bytes := wsData
				// if err != nil {
				// 	logrus.WithError(err).Error("websock cmd string base64 decoding failed")
				// }
				bytes := []byte(msgObj.Cmd)
				if _, err := ssConn.StdinPipe.Write(bytes); err != nil {
					logrus.WithError(err).Error("ws cmd bytes write to ssh.stdin pipe failed")
				}
				//write input cmd to log buffer
				if _, err := logBuff.Write(bytes); err != nil {
					logrus.WithError(err).Error("write received cmd into log buffer failed")
				}
			}
		}
	}
}
func (ssConn *SshConn) SendComboOutput(wsConn *websocket.Conn, exitCh chan bool) {
	//tells other go routine quit
	defer setQuit(exitCh)

	//every 100ms write combine output bytes into websocket response
	count := 0
	tick := time.NewTicker(time.Millisecond * time.Duration(100))
	//for range time.Tick(120 * time.Millisecond){}
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if count < 100 {
				count += 1
			} else {
				count = 0
			}
			// write combine output bytes into websocket response
			// count归零表示定时器达到了10s, 每隔10s如果没有实际数据(TextMessage)发送, 则发送心跳包(PingMessage)维持连接
			if err := flushComboOutput(ssConn.ComboOutput, wsConn, count); err != nil {
				logrus.WithError(err).Error("ssh sending combo output to webSocket failed")
				return
			}
		case <-exitCh:
			return
		}
	}
}

func (ssConn *SshConn) SessionWait(quitChan chan bool) {
	if err := ssConn.Session.Wait(); err != nil {
		logrus.WithError(err).Error("ssh session wait failed")
		setQuit(quitChan)
	}
}

func setQuit(ch chan bool) {
	ch <- true
}
