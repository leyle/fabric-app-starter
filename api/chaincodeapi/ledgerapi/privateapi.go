package ledgerapi

import (
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/api/context"
	"github.com/leyle/fabric-app-starter/chaincode/universal/ledgerstate"
)

func CreatePrivateState(ctx *context.ApiContext, channel, chaincode, collectionName, publicId, privateId, data string) *model.CCApiResponse {
	ccResp := model.NewCCApiResponse()
	ccResp.DataId = privateId

	contractRet := helper.GetContract(ctx, channel, chaincode)
	if contractRet.Err != nil {
		ccResp.Error = contractRet.Err
		return ccResp
	}
	defer contractRet.Close()
	contract := contractRet.Contract

	// transient key
	key := helper.GenerateTransientKey()
	transient := &ledgerstate.TransientForm{
		Id:             privateId,
		PublicDataId:   publicId,
		CollectionName: collectionName,
		Data:           data,
	}

	tBytes, err := json.Marshal(&transient)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Send()
		ccResp.Error = err
		return ccResp
	}

	tData := map[string][]byte{
		key: tBytes,
	}
	withT := gateway.WithTransient(tData)

	txn, err := contract.CreateTransaction(model.CCCreatePrivateState, withT)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Send()
		ccResp.Error = err
		return ccResp
	}

	ret, err := txn.Submit(key)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Send()
		ccResp.Error = err
		return ccResp
	}
	ccResp.CCRet = ret
	ctx.Logger().Info().Str("channel", channel).Str("chaincode", chaincode).Str("result", string(ret)).Msg("success")
	return ccResp
}
