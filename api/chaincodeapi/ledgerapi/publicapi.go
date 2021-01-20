package ledgerapi

import (
	"encoding/json"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/api/context"
	"github.com/leyle/fabric-app-starter/chaincode/universal/ledgerstate"
)

func CreatePublicState(ctx *context.ApiContext, channel, chaincode string, form *ledgerstate.PublicStateForm) *model.CCApiResponse {
	ccResp := model.NewCCApiResponse()
	ccResp.DataId = form.Id

	contractRet := helper.GetContract(ctx, channel, chaincode)
	if contractRet.Err != nil {
		ccResp.Error = contractRet.Err
		return ccResp
	}
	defer contractRet.Close()
	contract := contractRet.Contract

	// dump form struct to string
	args, err := json.Marshal(form)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Msg("marshal public state to string failed")
		ccResp.Error = err
		return ccResp
	}

	ret, err := contract.SubmitTransaction(model.CCCreatePublicState, string(args))
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Send()
		ccResp.Error = err
		return ccResp
	}

	ctx.Logger().Info().Str("channel", channel).Str("chaincode", chaincode).Str("result", string(ret)).Msg("success")
	ccResp.CCRet = ret
	return ccResp
}
