package terminal

import (
	"TTPanel/internal/global"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type safeBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *safeBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}
func (w *safeBuffer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}
func (w *safeBuffer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
}

const (
	wsMsgCmd       = "cmd"
	wsMsgResize    = "resize"
	WsMsgHeartbeat = "heartbeat"
)

type wsMsg struct {
	Type      string `json:"type"`
	Data      string `json:"data,omitempty"`
	Cols      int    `json:"cols,omitempty"`
	Rows      int    `json:"rows,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
}

type LogicSshWsSession struct {
	stdinPipe       io.WriteCloser
	comboOutput     *safeBuffer
	logBuff         *safeBuffer
	inputFilterBuff *safeBuffer
	session         *ssh.Session
	wsConn          *websocket.Conn
	Recorder        *Recorder
	isAdmin         bool
	IsFlagged       bool
}

func NewLogicSshWsSession(cols, rows int, isAdmin bool, sshClient *ssh.Client, wsConn *websocket.Conn, rec *Recorder) (*LogicSshWsSession, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(safeBuffer)
	logBuf := new(safeBuffer)
	inputBuf := new(safeBuffer)
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	if err := sshSession.Shell(); err != nil {
		return nil, err
	}
	lsws := &LogicSshWsSession{
		stdinPipe:       stdinP,
		comboOutput:     comboWriter,
		logBuff:         logBuf,
		inputFilterBuff: inputBuf,
		session:         sshSession,
		wsConn:          wsConn,
		isAdmin:         isAdmin,
		IsFlagged:       false,
	}
	if rec != nil {
		lsws.Recorder = rec
		lsws.Recorder.Lock()
		lsws.Recorder.WriteHeader(30, 150)
		lsws.Recorder.Unlock()
	}
	return lsws, nil
}

func (sws *LogicSshWsSession) Close() {
	if sws.session != nil {
		_ = sws.session.Close()
	}
	if sws.logBuff != nil {
		sws.logBuff = nil
	}
	if sws.comboOutput != nil {
		sws.comboOutput = nil
	}
}
func (sws *LogicSshWsSession) Start(quitChan chan bool, c *gin.Context) {
	go sws.receiveWsMsg(quitChan, c)
	go sws.sendComboOutput(quitChan)
}

func (sws *LogicSshWsSession) receiveWsMsg(exitCh chan bool, c *gin.Context) {
	wsConn := sws.wsConn
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				return
			}
			msgObj := wsMsg{}
			_ = json.Unmarshal(wsData, &msgObj)
			switch msgObj.Type {
			case wsMsgResize:
				global.Log.Debugf("ssh pty change windows size, cols: %d, rows: %d", msgObj.Cols, msgObj.Rows)
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err := sws.session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						global.Log.Errorf("ssh pty change windows size failed, err: %v", err)
					}
				}
			case wsMsgCmd:
				global.Log.Debugf("ssh receive cmd: %s", msgObj.Data)

				decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.Data)
				if err != nil {
					global.Log.Errorf("websock cmd string base64 decoding failed, err: %v", err)
				}
				//if sws.Recorder != nil {
				//	sws.Recorder.Lock()
				//	sws.Recorder.WriteData(OutPutType, string(decodeBytes))
				//	sws.Recorder.Unlock()
				//}

				sws.sendWebsocketInputCommandToSshSessionStdinPipe(decodeBytes)
			case WsMsgHeartbeat:
				global.Log.Debugf("ssh receive heartbeat")
				// 接收到心跳包后将心跳包原样返回，可以用于网络延迟检测等情况
				err = wsConn.WriteMessage(websocket.TextMessage, wsData)
				if err != nil {
					global.Log.Errorf("ssh sending heartbeat to webSocket failed, err: %v", err)
				}
			}
		}
	}
}

func (sws *LogicSshWsSession) sendWebsocketInputCommandToSshSessionStdinPipe(cmdBytes []byte) {
	if _, err := sws.stdinPipe.Write(cmdBytes); err != nil {
		global.Log.Errorf("ws cmd bytes write to ssh.stdin pipe failed, err: %v", err)
	}
}

func (sws *LogicSshWsSession) sendComboOutput(exitCh chan bool) {
	wsConn := sws.wsConn
	defer setQuit(exitCh)

	tick := time.NewTicker(time.Millisecond * time.Duration(60))
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if sws.comboOutput == nil {
				return
			}
			bs := sws.comboOutput.Bytes()
			if len(bs) > 0 {
				if sws.Recorder != nil {
					sws.Recorder.Lock()
					sws.Recorder.WriteData(OutPutType, string(bs))
					sws.Recorder.Unlock()
				}
				wsData, err := json.Marshal(wsMsg{
					Type: wsMsgCmd,
					Data: base64.StdEncoding.EncodeToString(bs),
				})
				if err != nil {
					global.Log.Errorf("encoding combo output to json failed, err: %v", err)
					continue
				}
				err = wsConn.WriteMessage(websocket.TextMessage, wsData)
				if err != nil {
					global.Log.Errorf("ssh sending combo output to webSocket failed, err: %v", err)
				}
				_, err = sws.logBuff.Write(bs)
				if err != nil {
					global.Log.Errorf("combo output to log buffer failed, err: %v", err)
				}
				sws.comboOutput.buffer.Reset()
			}
			if string(bs) == string([]byte{13, 10, 108, 111, 103, 111, 117, 116, 13, 10}) {
				sws.Close()
				return
			}

		case <-exitCh:
			return
		}
	}
}

func (sws *LogicSshWsSession) Wait(quitChan chan bool) {
	if err := sws.session.Wait(); err != nil {
		setQuit(quitChan)
	}
}

func setQuit(ch chan bool) {
	ch <- true
}
