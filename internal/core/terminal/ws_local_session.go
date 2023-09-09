package terminal

import (
	"TTPanel/internal/global"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"errors"
	"github.com/gorilla/websocket"
)

type LocalWsSession struct {
	slave  *LocalCommand
	wsConn *websocket.Conn

	writeMutex sync.Mutex
}

func NewLocalWsSession(cols, rows int, wsConn *websocket.Conn, slave *LocalCommand) (*LocalWsSession, error) {
	if err := slave.ResizeTerminal(cols, rows); err != nil {
		global.Log.Errorf("ssh pty change windows size failed, err: %v", err)
	}

	return &LocalWsSession{
		slave:  slave,
		wsConn: wsConn,
	}, nil
}

func (sws *LocalWsSession) Start(quitChan chan bool) {
	go sws.handleSlaveEvent(quitChan)
	go sws.receiveWsMsg(quitChan)
}

func (sws *LocalWsSession) handleSlaveEvent(exitCh chan bool) {
	defer setQuit(exitCh)
	defer global.Log.Debug("thread of handle slave event has exited now")

	buffer := make([]byte, 1024)
	for {
		select {
		case <-exitCh:
			return
		default:
			n, _ := sws.slave.Read(buffer)
			_ = sws.masterWrite(buffer[:n])
		}
	}
}

func (sws *LocalWsSession) masterWrite(data []byte) error {
	sws.writeMutex.Lock()
	defer sws.writeMutex.Unlock()
	wsData, err := json.Marshal(wsMsg{
		Type: wsMsgCmd,
		Data: base64.StdEncoding.EncodeToString(data),
	})
	if err != nil {
		return errors.New(fmt.Sprintf("LocalWsSession->masterWrite->Marshal->failed to encoding to json:%s", err.Error()))
	}
	err = sws.wsConn.WriteMessage(websocket.TextMessage, wsData)
	if err != nil {
		return errors.New(fmt.Sprintf("LocalWsSession->masterWrite->WriteMessage->WriteMessage->failed to write to master:%s", err.Error()))
	}
	return nil
}

func (sws *LocalWsSession) receiveWsMsg(exitCh chan bool) {
	wsConn := sws.wsConn
	defer setQuit(exitCh)
	defer global.Log.Debug("thread of receive ws msg has exited now")
	for {
		select {
		case <-exitCh:
			return
		default:
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				global.Log.Errorf("reading webSocket message failed, err: %v", err)
				return
			}
			msgObj := wsMsg{}
			_ = json.Unmarshal(wsData, &msgObj)
			switch msgObj.Type {
			case wsMsgResize:
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err := sws.slave.ResizeTerminal(msgObj.Cols, msgObj.Rows); err != nil {
						global.Log.Errorf("ssh pty change windows size failed, err: %v", err)
					}
				}
			case wsMsgCmd:
				decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.Data)
				if err != nil {
					global.Log.Errorf("websock cmd string base64 decoding failed, err: %v", err)
				}
				sws.sendWebsocketInputCommandToSshSessionStdinPipe(decodeBytes)
			case WsMsgHeartbeat:
				// 接收到心跳包后将心跳包原样返回，可以用于网络延迟检测等情况
				err = wsConn.WriteMessage(websocket.TextMessage, wsData)
				if err != nil {
					global.Log.Errorf("ssh sending heartbeat to webSocket failed, err: %v", err)
				}
			}
		}
	}
}

func (sws *LocalWsSession) sendWebsocketInputCommandToSshSessionStdinPipe(cmdBytes []byte) {
	if _, err := sws.slave.Write(cmdBytes); err != nil {
		global.Log.Errorf("ws cmd bytes write to ssh.stdin pipe failed, err: %v", err)
	}
}
