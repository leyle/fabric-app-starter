package chaincodeapi

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/chaincodeapi/private"
	"github.com/leyle/fabric-app-starter/chaincodeapi/public"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/returnfun"
)

// create handler
type CreateForm struct {
	// application name
	App string `json:"app" binding:"required"`

	// dataId, it should be unique in entire applications
	DataId string `json:"dataId" binding:"required"`

	// public data
	Public *public.PublicForm `json:"public"`

	// private data
	Private *private.PrivateForm `json:"private"`
}

func CreateHandler(ctx *context.ApiContext, c *gin.Context) {
	var form CreateForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)

	apiResp := model.ApiResponse{
		Result: false,
		App:    form.App,
		DataId: form.DataId,
	}

	// public data call public chaincode
	// private data call private chaincode

	// because private data require more permission control
	// so we create it first, if failed, then whole create failed
	// if create private data success, then create public data

	// 1. check if we need to create private data
	privateForm := form.Private
	if privateForm != nil {
		privateForm.App = form.App
		privateForm.DataId = form.DataId
		err = private.CallPrivateChainCode(ctx, c, privateForm)
	}

	// 2. check if we need to create public data
	publicForm := form.Public
	if publicForm != nil {
		publicForm.App = form.App
		publicForm.DataId = form.DataId
		resp := public.CallPublicChainCode(ctx, c, publicForm)
		if resp.Error != nil {
			apiResp.CCResp = resp.CCRet
			returnfun.ReturnJson(c, 400, 400, "", apiResp)
			return
		}
	}

	apiResp.Result = true
	returnfun.ReturnOKJson(c, apiResp)
}
