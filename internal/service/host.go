package service

import (
	"TTPanel/internal/core/ssh"
	"TTPanel/internal/core/terminal"
	"TTPanel/internal/global"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"os"
	"path"
	"time"
)

type HostService struct {
}

// AddHostCategory 添加主机分类
func (s *HostService) AddHostCategory(name, remark string) error {
	//查询主机分类是否存在
	hostCategory, err := (&model.HostCategory{Name: name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostCategory.ID > 0 {
		return errors.New("host category already exists")
	}
	_, err = (&model.HostCategory{
		Name:   name,
		Remark: remark,
	}).Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// HostCategoryList 主机分类列表
func (s *HostService) HostCategoryList() ([]*model.HostCategory, int64, error) {
	return (&model.HostCategory{}).List(global.PanelDB, &model.ConditionsT{}, 0, 0)
}

// EditHostCategory 编辑主机分类
func (s *HostService) EditHostCategory(id int64, name, remark string) error {
	//查询主机分类是否存在
	hostCategory, err := (&model.HostCategory{ID: id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostCategory.ID == 0 {
		return errors.New("host category does not exist")
	}
	//更新主机分类
	hostCategory.Name = name
	hostCategory.Remark = remark
	err = hostCategory.Update(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// GetHostByID 根据ID获取主机
func (s *HostService) GetHostByID(id int64) (*model.Host, error) {
	return (&model.Host{ID: id}).Get(global.PanelDB)
}

// DeleteHostCategory 删除主机分类
func (s *HostService) DeleteHostCategory(id int64) error {
	if id == 1 {
		return errors.New("default host category cannot be deleted")
	}
	//查询主机分类是否存在
	hostCategory, err := (&model.HostCategory{ID: id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostCategory.ID == 0 {
		return errors.New("host category does not exist")
	}
	//删除主机分类
	err = hostCategory.Delete(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		return err
	}
	return nil
}

// AddHost 添加主机
func (s *HostService) AddHost(param *request.AddHostR) error {
	//查询主机分类是否存在
	hostGet, err := (&model.Host{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostGet.ID > 0 {
		return errors.New("host already exists")
	}
	hostGet.CID = param.CId
	hostGet.Name = param.Name
	hostGet.Remark = param.Remark
	hostGet.Address = param.Address
	hostGet.Port = param.Port
	hostGet.User = param.User
	hostGet.Password = param.Password
	hostGet.PrivateKey = param.PrivateKey
	_, err = hostGet.Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// HostList 主机列表
func (s *HostService) HostList(query string, cid int64, offset, limit int) ([]*model.Host, int64, error) {
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	if !util.StrIsEmpty(query) {
		query = "%" + query + "%"
		whereT["name LIKE ?"] = query
		whereOrT["address LIKE ?"] = query
		whereOrT["remark LIKE ?"] = query
	}
	if cid > 0 {
		whereT["cid"] = cid
	}
	list, total, err := (&model.Host{}).List(global.PanelDB, &whereT, &whereOrT, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// DeleteHost 删除主机
func (s *HostService) DeleteHost(id int64) (*model.Host, error) {
	//查询主机是否存在
	host, err := (&model.Host{ID: id}).Get(global.PanelDB)
	if err != nil {
		return host, err
	}
	if host.ID == 0 {
		return host, errors.New("host not exist")
	}
	//删除主机
	err = host.Delete(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		return host, err
	}
	return host, nil
}

// ShortcutCommandList 快捷命令列表
func (s *HostService) ShortcutCommandList(query string, offset, limit int) ([]*model.HostShortcutCommand, int64, error) {
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	if !util.StrIsEmpty(query) {
		query = "%" + query + "%"
		whereT["name LIKE ?"] = query
		whereOrT["description LIKE ?"] = query
		whereOrT["cmd LIKE ?"] = query
	}
	list, total, err := (&model.HostShortcutCommand{}).List(global.PanelDB, &whereT, &whereOrT, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// AddShortcutCommand 添加快捷命令
func (s *HostService) AddShortcutCommand(param *request.AddShortcutCommandR) error {
	//查询快捷命令是否存在
	hostShortcutCommand, err := (&model.HostShortcutCommand{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostShortcutCommand.ID > 0 {
		return errors.New("shortcut command name already exists")
	}
	hostShortcutCommand.Name = param.Name
	hostShortcutCommand.Description = param.Description
	hostShortcutCommand.Cmd = param.Cmd
	_, err = hostShortcutCommand.Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// EditShortcutCommand 编辑快捷命令
func (s *HostService) EditShortcutCommand(param *request.EditShortcutCommandR) error {
	//查询快捷命令是否存在
	hostShortcutCommand, err := (&model.HostShortcutCommand{ID: param.ID}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostShortcutCommand.ID == 0 {
		return errors.New("shortcut command does not exist")
	}
	//查询快捷命令名称是否存在
	hostShortcutCommandName, err := (&model.HostShortcutCommand{Name: param.Name}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if hostShortcutCommandName.ID > 0 && hostShortcutCommandName.ID != param.ID {
		return errors.New("shortcut command name already exists")
	}
	hostShortcutCommand.Name = param.Name
	hostShortcutCommand.Description = param.Description
	hostShortcutCommand.Cmd = param.Cmd
	err = hostShortcutCommand.Update(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// DeleteShortcutCommand 删除快捷命令
func (s *HostService) DeleteShortcutCommand(id int64) (*model.HostShortcutCommand, error) {
	//查询快捷命令是否存在
	hostShortcutCommand, err := (&model.HostShortcutCommand{ID: id}).Get(global.PanelDB)
	if err != nil {
		return hostShortcutCommand, err
	}
	if hostShortcutCommand.ID == 0 {
		return hostShortcutCommand, errors.New("shortcut command does not exist")
	}
	//删除快捷命令
	err = hostShortcutCommand.Delete(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		return hostShortcutCommand, err
	}
	return hostShortcutCommand, nil
}

// Terminal 主机终端
func (s *HostService) Terminal(c *gin.Context, param *request.TerminalR, host *model.Host) {

	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Log.Errorf("upGrader.Upgrade:gin context http handler failed, err: %v", err)
		return
	}
	defer func(wsConn *websocket.Conn) {
		_ = wsConn.Close()
	}(wsConn)
	c.Set("HOST_ID", host.ID)
	c.Set("HOST_NAME", host.Name)

	var recorder *terminal.Recorder
	castSavePath := global.Config.System.PanelPath + "/data/cast"
	_ = os.MkdirAll(castSavePath, 0766)
	fileName := path.Join(castSavePath, fmt.Sprintf("%s_%s_%s.cast", host.Address, host.User, time.Now().Format("20060102_150405")))
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0766)
	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"status": false, "msg": err.Error()})
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	recorder = terminal.NewRecorder(f)

	//连接主机
	var connInfo ssh.ConnInfo
	connInfo.Address = host.Address
	connInfo.Port = host.Port
	connInfo.User = host.User
	connInfo.Password = host.Password
	connInfo.PrivateKey = []byte(host.PrivateKey)

	client, err := connInfo.NewClient()
	if err != nil {
		if wsHandleError(wsConn, errors.New("failed to set up the connection. Please check the host information.ERROR:"+err.Error())) {
			return
		}
	}

	defer client.Close()

	ssConn, err := connInfo.NewSshConn(param.Cols, param.Rows)
	if err != nil {
		if wsHandleError(wsConn, err) {
			return
		}
	}
	defer ssConn.Close()

	sws, err := terminal.NewLogicSshWsSession(param.Cols, param.Rows, true, connInfo.Client, wsConn, recorder)
	if err != nil {
		if wsHandleError(wsConn, err) {
			return
		}
	}
	defer sws.Close()
	//fmt.Println(wsConn)
	quitChan := make(chan bool, 3)
	sws.Start(quitChan, c)
	go sws.Wait(quitChan)

	<-quitChan

	if err != nil {
		if wsHandleError(wsConn, err) {
			return
		}
	}

}
