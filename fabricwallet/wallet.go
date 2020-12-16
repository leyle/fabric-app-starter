package fabricwallet

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/leyle/fabric-app-starter/context"
	. "github.com/leyle/ginbase/consolelog"
)

// process ca user register and enroll
// save ca credentials into file fabricwallet

func NewWallet(ctx *context.ApiContext) (*gateway.Wallet, error) {
	wallet, err := gateway.NewFileSystemWallet(ctx.Cfg.Fabric.WalletPath)
	if err != nil {
		Logger.Errorf("", "create file system fabricwallet failed, %s", err.Error())
		return nil, err
	}
	return wallet, nil
}

func NewGateway(ctx *context.ApiContext, userId string) (*gateway.Gateway, error) {
	gwCfg := gateway.WithConfig(config.FromFile(ctx.Cfg.Fabric.CCPath))
	identity := gateway.WithIdentity(ctx.Wallet, userId)
	gw, err := gateway.Connect(gwCfg, identity)
	if err != nil {
		Logger.Errorf("", "create fabric gateway failed, userId[%s], reason: %s", userId, err.Error())
		return nil, err
	}
	return gw, nil
}
