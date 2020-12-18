package authproxy

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/fabric-app-starter/fabricwallet"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
)

const KeyGateway = "KEYGATEWAY"
const keyClaim = "KEYCLAIM"

type JWTClaim struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
	valid bool `json:"-"`
}

func ParseJWTToken(ctx *context.ApiContext, c *gin.Context) *JWTClaim {
	jwtKey := []byte(ctx.Cfg.JWT.Secret)
	claim := &JWTClaim{}
	token := c.Request.Header.Get("X-TOKEN")
	if token == "" {
		Logger.Errorf(middleware.GetReqId(c), "No token in request headers")
		claim.valid = false
		return claim
	}

	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		Logger.Errorf(middleware.GetReqId(c), "Parse token failed, %s", err.Error())
		claim.valid = false
		return claim
	}

	if !tkn.Valid {
		Logger.Errorf(middleware.GetReqId(c), "parse token, token invalid")
		claim.valid = false
		return claim
	}

	claim.valid = true
	return claim
}

func Auth(ctx *context.ApiContext, c *gin.Context) {
	// 1. check token is valid
	claim := ParseJWTToken(ctx, c)
	if !claim.valid {
		returnfun.Return401Json(c, "Invalid token")
		return
	}
	Logger.Infof(middleware.GetReqId(c), "token check ok, current user: id[%s]|userId[%s]|username[%s]|role[%s]", claim.Id, claim.UserId, claim.Username, claim.Role)

	// 2. check user has enrolled
	wallet, err := fabricwallet.NewWallet(ctx)
	if err != nil {
		Logger.Errorf(middleware.GetReqId(c), "Get fabric wallet failed, %s", err.Error())
		returnfun.Return401Json(c, "Get wallet info failed")
		return
	}
	userId := claim.UserId
	if !wallet.Exists(userId) {
		Logger.Errorf(middleware.GetReqId(c), "user[%s] doesn't have wallet credentials", userId)
		returnfun.Return403Json(c, "User doesn't have wallet credentials, enroll it first")
		return
	}

	// 3. get fabric gateway
	gw, err := fabricwallet.NewGateway(ctx, userId)
	if err != nil {
		returnfun.ReturnJson(c, 500, 500, "GetGateway failed", nil)
		return
	}
	defer gw.Close()

	// 4. save to current user's context
	SetUser(c, claim)
	// SetGateway(c, gw)

	c.Next()
}

func SetUser(c *gin.Context, claim *JWTClaim) {
	c.Set(keyClaim, claim)
}

func SetGateway(c *gin.Context, gw *gateway.Gateway) {
	c.Set(KeyGateway, gw)
}

func GetCurUser(c *gin.Context) *JWTClaim {
	claim, exist := c.Get(keyClaim)
	if !exist {
		return nil
	}
	result := claim.(*JWTClaim)
	return result
}

func GetGateway(c *gin.Context) *gateway.Gateway {
	gw, exist := c.Get(KeyGateway)
	if !exist {
		return nil
	}
	result := gw.(*gateway.Gateway)
	return result
}
