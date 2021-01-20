package chaincodeapi

import (
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/ledgerapi"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/api/context"
	"github.com/leyle/fabric-app-starter/chaincode/universal/ledgerstate"
	"github.com/leyle/go-api-starter/ginhelper"
)

type CreateStateForm struct {
	DataId  string                  `json:"dataId"`
	AppName string                  `json:"appName"`
	Public  *CreatePublicStateForm  `json:"public"`
	Private *CreatePrivateStateForm `json:"private"`
}

type CreatePublicStateForm struct {
	Channel    string `json:"channel" binding:"required"`
	Chaincode  string `json:"chaincode" binding:"required"`
	DataString string `json:"dataString" binding:"required"`
}

type CreatePrivateStateForm struct {
	Channel    string   `json:"channel" binding:"required"`
	Chaincode  string   `json:"chaincode" binding:"required"`
	DataString string   `json:"dataString" binding:"required"`
	ShareNames []string `json:"shareNames"`
}

func CreateStateHandler(ctx *context.ApiContext) {
	var form CreateStateForm
	err := ctx.C.BindJSON(&form)
	ginhelper.StopExec(err)

	// check if data id is valid uuid.
	if !helper.IsValidUUID(form.DataId) {
		ginhelper.ReturnErrJson(ctx.C, "Invalid dataId, it should be an uuid value")
		return
	}

	apiResp := model.NewApiResponse(form.AppName, form.DataId)

	publicReq := &ledgerstate.PublicStateForm{
		Id:      form.DataId,
		AppName: form.AppName,
		Data:    form.Public.DataString,
	}

	// 1. if has private data, we process private first
	if form.Private != nil {
		privateId := helper.GeneratePrivateDataId()
		privateMetaInfo := &ledgerstate.PrivateMetaInfo{
			SystemId:              privateId,
			Channel:               form.Private.Channel,
			Chaincode:             form.Private.Chaincode,
			CreatorCollectionName: ctx.Cfg.PrivateCfg.NoShare,
			ShareCollectionNames:  form.Private.ShareNames,
		}
		publicReq.PrivateMetaInfo = privateMetaInfo

		// 1.1. save private data into shared org's collection
		for _, name := range form.Private.ShareNames {
			ccResp := ledgerapi.CreatePrivateState(ctx, form.Private.Channel, form.Private.Chaincode, name, form.DataId, privateId, form.Private.DataString)
			if ccResp.Error != nil {
				apiResp.ErrMsg = ccResp.Error.Error()
				ginhelper.ReturnJson(ctx.C, 400, 400, apiResp.ErrMsg, apiResp)
				return
			}
		}
		// 1.2. save private data into self collection
		ccResp := ledgerapi.CreatePrivateState(ctx, form.Private.Channel, form.Private.Chaincode, ctx.Cfg.PrivateCfg.NoShare, form.DataId, privateId, form.Private.DataString)
		if ccResp.Error != nil {
			apiResp.ErrMsg = ccResp.Error.Error()
			ginhelper.ReturnJson(ctx.C, 400, 400, apiResp.ErrMsg, apiResp)
			return
		}

		ctx.Logger().Info().Msg("create private data success")
		apiResp.Success = model.ResponseSuccessPartial
	}

	// 2. save public data
	ccResp := ledgerapi.CreatePublicState(ctx, form.Public.Channel, form.Public.Chaincode, publicReq)
	if ccResp.Error != nil {
		apiResp.ErrMsg = ccResp.Error.Error()
		ginhelper.ReturnJson(ctx.C, 400, 400, apiResp.ErrMsg, apiResp)
		return
	}

	apiResp.Success = model.ResponseSuccessAll
	ginhelper.ReturnOKJson(ctx.C, apiResp)
	return
}
