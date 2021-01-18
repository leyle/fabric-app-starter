package public

import (
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/api/context"
)

// call public chaincode
type CreatePublicForm struct {
	App    string `json:"-"`
	DataId string `json:"-"`

	Channel   string `json:"channel" binding:"required"`
	ChainCode string `json:"chaincode" binding:"required"`

	DataJson string `json:"dataJson" binding:"required"`
}

func CallPublicChainCodeCreate(ctx *context.ApiContext, form *CreatePublicForm) *model.CCApiResponse {
	ccResp := model.NewCCApiResponse()
	ccResp.DataId = form.DataId

	contractRet := helper.GetContract(ctx, form.Channel, form.ChainCode)
	if contractRet.Err != nil {
		ccResp.Error = contractRet.Err
		return ccResp
	}
	defer contractRet.Close()
	contract := contractRet.Contract

	// call chaincode's create method
	// chaincode's create method locates on chaincode/public/chaincode.go#Create
	// func (uc *UniversalContract) Create(ctx contractapi.TransactionContextInterface, app, dataId, data string) error
	ret, err := contract.SubmitTransaction(model.CCNameCreate, form.App, form.DataId, form.DataJson)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", form.Channel).Str("chaincode", form.ChainCode).Msg("call public chaincode, submit failed")
		ccResp.Error = err
		return ccResp
	}
	ctx.Logger().Info().Str("channel", form.Channel).Str("chaincode", form.ChainCode).Str("result", string(ret)).Msg("success")
	ccResp.CCRet = ret

	return ccResp
}

func CallPublicChaincodeGetById(ctx *context.ApiContext, form *model.GetByIdForm) *model.CCApiResponse {
	ccResp := model.NewCCApiResponse()
	ccResp.DataId = form.DataId
	contractRet := helper.GetContract(ctx, form.Public.Channel, form.Public.ChainCode)
	if contractRet.Err != nil {
		ccResp.Error = contractRet.Err
		return ccResp
	}
	defer contractRet.Close()
	contract := contractRet.Contract

	// func (uc *UniversalContract) GetById(ctx contractapi.TransactionContextInterface, app, dataId string) (*StorageOut, error)
	ret, err := contract.EvaluateTransaction(model.CCNameGetById, form.App, form.DataId)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", form.Public.Channel).Str("chaincode", form.Public.ChainCode).Msg("call public chaincode, get by id failed")
		ccResp.Error = err
		return ccResp
	}

	ctx.Logger().Info().Str("channel", form.Public.Channel).Str("chaincode", form.Public.ChainCode).Str("result", string(ret)).Msg("success")
	ccResp.CCRet = ret

	return ccResp
}
