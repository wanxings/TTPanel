package app

import (
	"TTPanel/internal/global"
	"TTPanel/internal/model"
	"TTPanel/pkg/util"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UID        int64  `json:"uid"`
	USERNAME   string `json:"username"`
	AdminToken string `json:"admin_token"`
	UseIP      string `json:"use_ip"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(global.Config.System.JwtSecret)
}

func GenerateToken(User *model.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(global.Config.System.JwtExpire) * time.Second)
	adminToken := util.EncodeMD5(User.Username + User.Password)
	global.GoCache.Set("admin_token", adminToken, -1)
	claims := Claims{
		UID:        User.ID,
		USERNAME:   User.Username,
		AdminToken: adminToken,
		UseIP:      User.LoginIp,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.Config.System.JwtIssuer + ":" + User.Salt,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
