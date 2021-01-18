package helper

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/leyle/fabric-app-starter/api/context"
	"github.com/leyle/fabric-user-manager/jwtwrapper"
)

type AppContract struct {
	Err      error
	Contract *gateway.Contract
	Gw       *gateway.Gateway
}

func New() *AppContract {
	return &AppContract{}
}

func (a *AppContract) Close() {
	a.Gw.Close()
}

func GetContract(ctx *context.ApiContext, channel, chaincode string) *AppContract {
	ac := New()
	curUser := jwtwrapper.GetCurUser(ctx.C)
	gw, err := jwtwrapper.NewGateway(ctx.JWTCtx, curUser.UserName)
	if err != nil {
		ac.Err = err
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Msg("create fabric gateway failed")
		return ac
	}
	ac.Gw = gw

	network, err := gw.GetNetwork(channel)
	if err != nil {
		ac.Err = err
		ctx.Logger().Error().Err(err).Str("channel", channel).Str("chaincode", chaincode).Msg("create fabric network failed")
		return ac
	}

	contract := network.GetContract(chaincode)
	ac.Contract = contract
	ctx.Logger().Debug().Str("channel", channel).Str("chaincode", chaincode).Msg("get contract success")
	return ac
}
