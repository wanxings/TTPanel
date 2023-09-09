package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

type UserService struct{}

const LoginErrKey = "TTPanelUserLoginErr"
const MaxLoginErrTimes = 10

// DoLogin 用户认证
func (u *UserService) DoLogin(ctx *gin.Context, param *request.Login) (*model.User, error) {
	errLoginCountInt := 0
	// 检查登录错误次数
	if errLoginCount, ok := global.GoCache.Get(fmt.Sprintf("%s:%s", LoginErrKey, ctx.ClientIP())); ok {
		errLoginCountInt = errLoginCount.(int)
		if errLoginCountInt >= MaxLoginErrTimes {
			return nil, errors.New(helper.Message("user.TooManyLoginError"))
		}

	}

	user, err := (&model.User{Username: param.UserName}).Get(global.PanelDB)
	if err == nil && user.ID > 0 {
		// 对比密码是否正确
		if ValidPassword(user.Password, param.Password, user.Salt) {
			// 清空登录计数
			global.GoCache.Delete(fmt.Sprintf("%s:%s", LoginErrKey, ctx.ClientIP()))
			ctx.Set("UID", user.ID)
			ctx.Set("USERNAME", user.Username)
			user.LoginIp = ctx.ClientIP()
			user.LoginTime = time.Now().Unix()
			err = user.Update(global.PanelDB)
			if err != nil {
				return nil, err
			}

			return user, nil
		}
	}

	// 登录错误计数
	err = global.GoCache.Increment(fmt.Sprintf("%s:%s", LoginErrKey, ctx.ClientIP()), 1)
	if err != nil {
		global.Log.Errorf("DoLogin->global.GoCache.Increment Error:%s", err.Error())
		global.GoCache.Set(fmt.Sprintf("%s:%s", LoginErrKey, ctx.ClientIP()), 1, time.Minute*60)
	}
	return nil, errors.New(helper.MessageWithMap("user.UsernameOrPasswordError", map[string]any{"Count": MaxLoginErrTimes - (errLoginCountInt + 1)}))
}

func setupLogin(user *model.User, c *gin.Context) error {
	session := sessions.Default(c)
	session.Set("id", user.ID)
	session.Set("username", user.Username)
	if adminToken, ok := global.GoCache.Get("admin_token"); ok {
		session.Set("admin_token", adminToken)
	} else {
		adminToken = util.EncodeMD5(user.Username + user.Password)
		global.GoCache.Set("admin_token", adminToken, -1)
		session.Set("admin_token", adminToken)
	}
	c.Set("UID", user.ID)
	c.Set("USERNAME", user.Username)
	//session.Options(sessions.Options{MaxAge: global.Config.System.SessionExpire})
	err := session.Save()
	if err != nil {
		return err
	}
	return nil
}

// Info 用户信息
func (u *UserService) Info(id int64) (*model.User, error) {
	user, err := (&model.User{ID: id}).Get(global.PanelDB)
	if err != nil {
		return nil, err
	}

	if user.ID > 0 {
		return user, nil
	}

	return nil, errors.New(helper.Message("user.administratorDoesNotExist"))
}

// ValidPassword 检查密码是否一致
func ValidPassword(dbPassword, password, salt string) bool {
	global.Log.Debugf("ValidPassword->dbPassword:%s,password:%s,salt:%s,", dbPassword, password, salt)
	global.Log.Debugf("ValidPassword->md5Password:%s,", util.EncodeMD5(util.EncodeMD5(password)+salt))
	return strings.Compare(dbPassword, util.EncodeMD5(util.EncodeMD5(password)+salt)) == 0
}

// EncryptPasswordAndSalt 密码加密&生成salt
func EncryptPasswordAndSalt(password string) (string, string) {
	salt := uuid.Must(uuid.NewV4()).String()[:8]
	password = util.EncodeMD5(util.EncodeMD5(password) + salt)

	return password, salt
}

// CreateTemporaryUser 创建临时用户
func (u *UserService) CreateTemporaryUser(expireTime int, remark string) (*model.TemporaryUser, error) {
	token := uuid.Must(uuid.NewV4()).String()
	temporaryUser := &model.TemporaryUser{
		Token:      token,
		ExpireTime: expireTime,
		Remark:     remark,
	}
	temporaryUserC, err := temporaryUser.Create(global.PanelDB)
	if err != nil {
		return nil, err
	}
	return temporaryUserC, nil
}
