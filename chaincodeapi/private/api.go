package private

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/leyle/fabric-app-starter/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/context"
)

const (
	TransientActionAdd    = "CREATE"
	TransientActionDelete = "DELETE"
	TransientActionUpdate = "UPDATE"
)

type CreatePrivateForm struct {
	App    string `json:"-"`
	DataId string `json:"-"`

	Channel        string `json:"channel" binding:"required"`
	ChainCode      string `json:"chaincode" binding:"required"`
	CollectionName string `json:"collectionName" binding:"required"`

	DataJson string `json:"dataJson" binding:"required"`
}

type TransientInput struct {
	CollectionName string `json:"collectionName"`
	App            string `json:"app"`    // client app name
	DataId         string `json:"dataId"` // client data unique id
	Data           string `json:"data"`   // user input data, json encoded
}

func CallPrivateChainCode(ctx *context.ApiContext, form *CreatePrivateForm) *model.CCApiResponse {
	ccResp := model.NewCCApiResponse()
	ccResp.DataId = form.DataId

	contractRet := helper.GetContract(ctx, form.Channel, form.ChainCode)
	if contractRet.Err != nil {
		ccResp.Error = contractRet.Err
		return ccResp
	}
	defer contractRet.Close()
	contract := contractRet.Contract

	// generate private data
	key := generateTransientKey(form.CollectionName, TransientActionAdd)
	transient := &TransientInput{
		CollectionName: form.CollectionName,
		App:            form.App,
		DataId:         form.DataId,
		Data:           form.DataJson,
	}

	tByte, err := json.Marshal(&transient)
	if err != nil {
		ctx.Logger().Error().Err(err).Send()
		ccResp.Error = err
		return ccResp
	}

	tData := map[string][]byte{
		key: tByte,
	}
	withT := gateway.WithTransient(tData)

	txn, err := contract.CreateTransaction(model.CCNameCreate, withT)
	if err != nil {
		ccResp.Error = err
		return ccResp
	}

	ret, err := txn.Submit(key)
	if err != nil {
		ccResp.Error = err
		return ccResp
	}

	ccResp.CCRet = ret
	return ccResp
}

func CallPrivateChaincodeGetById(ctx *context.ApiContext, form *model.GetByIdForm) *model.CCApiResponse {
	ccResp := model.NewCCApiResponse()
	ccResp.DataId = form.DataId
	contractRet := helper.GetContract(ctx, form.Private.Channel, form.Private.ChainCode)
	if contractRet.Err != nil {
		ccResp.Error = contractRet.Err
		return ccResp
	}
	defer contractRet.Close()
	contract := contractRet.Contract

	ret, err := contract.EvaluateTransaction(model.CCNameGetById, form.Private.CollectionName, form.App, form.DataId)
	if err != nil {
		ccResp.Error = err
		return ccResp
	}

	ccResp.CCRet = ret
	return ccResp
}

func generateTransientKey(collection string, action string) string {
	return fmt.Sprintf("%s%s", collection, action)
}
