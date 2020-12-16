package jwtserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leyle/fabric-app-starter/context"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/util"
	"net/http"
)

const (
	FabricCaApiRegisterAndEnroll = "/api/ca/register"
	FabricCaApiEnroll            = "/api/ca/enroll"
)

func getFullApi(ctx *context.ApiContext, path string) string {
	url := fmt.Sprintf("http://%s:%s%s", ctx.Cfg.Express.Host, ctx.Cfg.Express.Port, path)
	Logger.Debugf("", "Ca Url: %s", url)
	return url
}

func RegisterAndEnrollCaUser(ctx *context.ApiContext, userId, passwd string) error {
	url := getFullApi(ctx, FabricCaApiRegisterAndEnroll)

	type RegisterForm struct {
		UserId string `json:"userId"`
		Passwd string `json:"passwd"`
	}

	form := &RegisterForm{
		UserId: userId,
		Passwd: passwd,
	}
	data, _ := json.Marshal(form)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	resp, err := util.HttpPost(url, data, headers)
	if err != nil {
		Logger.Errorf("", "register and enroll ca user[%s] failed, %s", userId, err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {
		Logger.Errorf("", "register and enroll ca user[%s] failed, %d", userId, resp.StatusCode)
		return errors.New("register and enroll ca user failed")
	}

	return nil
}
