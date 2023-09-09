package api

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/app"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/helper/errcode"
	"TTPanel/internal/model/request"
	"github.com/gin-gonic/gin"
)

type UserApi struct{}

// Login 用户登录
func (u *UserApi) Login(c *gin.Context) {
	param := request.Login{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	user, err := ServiceGroupApp.UserServiceApp.DoLogin(c, &param)

	if err != nil {
		global.Log.Errorf("service.DoLogin err: %v", err)
		go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByUserLogin, err.Error())
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	token, err := app.GenerateToken(user)
	if err != nil {
		global.Log.Errorf("app.GenerateToken err: %v", err)
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	err = ServiceGroupApp.MonitorServiceApp.LoginPanelEvent(c.ClientIP())
	if err != nil {
		global.Log.Errorf("Login->ServiceGroupApp.MonitorServiceApp.LoginPanelEvent err: %v", err)
	}
	go ServiceGroupApp.LogAuditServiceApp.WriteOperationLog(c, constant.OperationLogTypeByUserLogin, helper.Message("tip.LoginSuccess"))
	response.ToResponse(gin.H{
		"token": token,
	})
}

// Info 获取用户基本信息
func (u *UserApi) Info(c *gin.Context) {
	param := request.Login{}
	response := app.NewResponse(c)

	if user, exists := c.Get("USERNAME"); exists {
		param.UserName = user.(string)
	}

	user, err := ServiceGroupApp.UserServiceApp.Info(1)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     "admin",
	})

}

// Logout 用户登出
func (u *UserApi) Logout(c *gin.Context) {
	response := app.NewResponse(c)
	//global.GoCache.Set(fmt.Sprintf("block_jwt_token:%s", LoginErrKey, ctx.ClientIP()), 1, time.Minute*60)
	//session := sessions.Default(c)
	//session.Clear()
	//_ = session.Save()
	response.ToResponseMsg(helper.Message("tips.LogoutSuccess"))
}

//func (u *UserApi) GetCaptcha(c *gin.Context) {
//	caps := captcha.New()
//	if err := caps.SetFont("./config/comic.ttf"); err != nil {
//		panic(err.Error())
//	}
//
//	caps.SetSize(160, 64)
//	caps.SetDisturbance(captcha.MEDIUM)
//	caps.SetFrontColor(color.RGBA{A: 255})
//	caps.SetBkgColor(color.RGBA{R: 218, G: 240, B: 228, A: 255})
//	img, password := caps.Create(4, captcha.NUM)
//	emptyBuff := bytes.NewBuffer(nil)
//	_ = png.Encode(emptyBuff, img)
//
//	key := util.EncodeMD5(uuid.Must(uuid.NewV4()).String())
//
//	// 五分钟有效期
//	global.GoCache.Set("LoginCaptcha:"+key, password, time.Minute*1)
//
//	response := app.NewResponse(c)
//	response.ToResponse(gin.H{
//		"id":   key,
//		"b64s": "data:image/png;base64," + base64.StdEncoding.EncodeToString(emptyBuff.Bytes()),
//	})
//}

//func userFrom(c *gin.Context) (*model.User, bool) {
//	if u, exists := c.Get("USER"); exists {
//		user, ok := u.(*model.User)
//		return user, ok
//	}
//	return nil, false
//}

// CreateTemporaryUser 创建临时用户
func (u *UserApi) CreateTemporaryUser(c *gin.Context) {
	param := request.CreateTemporaryUserR{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Log.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	temporaryUser, err := ServiceGroupApp.UserServiceApp.CreateTemporaryUser(param.ExpireTime, param.Remark)
	if err != nil {
		global.Log.Errorf("service.CreateTemporaryUser err: %v", err)
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponse(temporaryUser)
}
