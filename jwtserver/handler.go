package jwtserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/authproxy"
	"github.com/leyle/fabric-app-starter/context"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
	"github.com/leyle/ginbase/util"
	"strings"
)

// when login in
// first check username and passwd
// then check if had exported ca credentials
// if not, try to enroll and export it
// if user not register to ca, return failed, in this situation, maybe the system has something wrong
type LoginCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(ctx *context.ApiContext, c *gin.Context) {
	var form LoginCredentials
	err := c.BindJSON(&form)
	middleware.StopExec(err)

	form.Username = strings.TrimSpace(form.Username)
	form.Password = strings.TrimSpace(form.Password)

	// check user and passwd from sqlite db
	ds := context.NewDs(ctx.DbFile)
	defer ds.Close()

	var account UserAccount
	tx := ds.Db.Where("username = ?", form.Username).First(&account)
	if tx.RowsAffected == 0 {
		Logger.Warnf(middleware.GetReqId(c), "Login by[%s], but no record", form.Username)
		returnfun.Return401Json(c, "username doesn't exist")
		return
	}

	passHash := util.GenerateHashPasswd(account.Salt, form.Password)
	if passHash != account.PassHash {
		Logger.Warnf(middleware.GetReqId(c), "Login by[%s], wrong passwd", form.Username)
		returnfun.Return401Json(c, "invalid password")
		return
	}

	// check if use has register for ca, todo

	// generate jwt token
	token, err := CreateJWTToken(ctx, c, &account)
	middleware.StopExec(err)

	returnfun.ReturnOKJson(c, gin.H{"token": token})
	return
}

type CreateUserForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func CreateUserHandler(ctx *context.ApiContext, c *gin.Context) {
	// check caller's permission
	// it should be admin user
	claim := authproxy.GetCurUser(c)
	if claim == nil {
		returnfun.Return401Json(c, "Get Current User's info failed")
		return
	}
	if claim.Role != UserRoleAdmin {
		emsg := fmt.Sprintf("Current User[%s][%s] cannot create user", claim.UserId, claim.Role)
		Logger.Error(middleware.GetReqId(c), emsg)
		returnfun.Return403Json(c, emsg)
		return
	}

	var form CreateUserForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)

	form.Username = strings.TrimSpace(form.Username)
	form.Password = strings.TrimSpace(form.Password)
	form.Role = strings.TrimSpace(form.Role)

	ds := context.NewDs(ctx.DbFile)
	defer ds.Close()

	// register user has two step
	// 1. create useraccount
	// 2. register/enroll ca

	// 1. check if username repeat
	dbua, err := GetUserByUsername(ds, form.Username)
	middleware.StopExec(err)
	if dbua != nil {
		returnfun.ReturnErrJson(c, "username has exist")
		return
	}

	// 2. save user
	salt := util.GenerateDataId()
	passHash := util.GenerateHashPasswd(salt, form.Password)
	ua := &UserAccount{
		Id:        util.GenerateDataId(),
		Username:  form.Username,
		Role:      form.Role,
		Salt:      salt,
		PassHash:  passHash,
		Valid:     UserStatusValid,
		CreatedAt: util.CurUnixTime(),
	}
	ua.UpdatedAt = ua.CreatedAt
	err = SaveUserAccount(ds, ua)
	middleware.StopExec(err)

	// call caapi register/enroll ca
	err = RegisterAndEnrollCaUser(ctx, ua.Id, ua.PassHash)
	if err != nil {
		returnfun.ReturnErrJson(c, err.Error())
		return
	}

	returnfun.ReturnOKJson(c, "")
	return
}
