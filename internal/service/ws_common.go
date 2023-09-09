package service

import (
	"TTPanel/internal/global"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func wsHandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		global.Log.Errorf("handler ws faled:, err: %v", err)
		dt := time.Now().Add(time.Second)
		if wcErr := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); wcErr != nil {
			_ = ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
		return true
	}
	return false
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
