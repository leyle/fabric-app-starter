package jwtserver

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/authproxy"
	"github.com/leyle/fabric-app-starter/context"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/util"
	"time"
)

func CreateJWTToken(ctx *context.ApiContext, c *gin.Context, user *UserAccount) (string, error) {
	jwtKey := []byte(ctx.Cfg.JWT.Secret)
	expireTime := time.Now().Add(time.Duration(ctx.Cfg.JWT.Expire) * time.Hour)
	claim := &authproxy.JWTClaim{
		UserId:   user.Id,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  util.CurUnixTime(),
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		Logger.Errorf(middleware.GetReqId(c), "create jwt, sign failed, %s", err.Error())
		return "", err
	}
	return tokenStr, nil
}
