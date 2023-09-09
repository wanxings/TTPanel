package api

import (
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/helper/web_ssh"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net/http"
	"strconv"
	"time"
)

type WebSSHApi struct{}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsSSH
// @Tags      System
// @Summary   wsSSH
// @Router    /system/GetNetWork [get]
func (s *WebSSHApi) WsSSH(c *gin.Context) {
	response := app.NewResponse(c)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer func(wsConn *websocket.Conn) {
		_ = wsConn.Close()
	}(wsConn)

	//cIp := c.ClientIP()

	cols, err := strconv.Atoi(c.DefaultQuery("cols", "80"))
	if wsHandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "40"))
	if wsHandleError(wsConn, err) {
		return
	}

	HostInfo, err := ServiceGroupApp.HostServiceApp.GetHostByID(int64(id))
	if err != nil {
		return
	}

	client, err := ServiceGroupApp.WebSSHServiceApp.NewSshClient(HostInfo)
	if wsHandleError(wsConn, err) {
		return
	}
	defer func(client *ssh.Client) {
		_ = client.Close()
	}(client)
	//startTime := time.Now()
	//ssConn, err := util.NewSshConn(cols, rows, client)
	//if wsHandleError(wsConn, err) {
	//	return
	//}
	//defer ssConn.Close()

	sws, err := web_ssh.NewLogicSshWsSession(cols, rows, true, client, wsConn)
	if wsHandleError(wsConn, err) {
		return
	}
	defer sws.Close()

	quitChan := make(chan bool, 3)
	sws.Start(quitChan)
	go sws.Wait(quitChan)

	<-quitChan
	//保存日志

	//write logs
	//xtermLog := model.SshLog{
	//	StartedAt: startTime,
	//	UserId:    userM.Id,
	//	Log:       sws.LogString(),
	//	MachineId: idx,
	//	ClientIp:  cIp,
	//}
	//err = xtermLog.Create()
	if wsHandleError(wsConn, err) {
		return
	}

}
func wsHandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		logrus.WithError(err).Error("handler ws ERROR:")
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logrus.WithError(err).Error("websocket writes control message failed:")
		}
		return true
	}
	return false
}
