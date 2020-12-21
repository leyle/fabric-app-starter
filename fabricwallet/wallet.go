package fabricwallet

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/leyle/fabric-app-starter/context"
	. "github.com/leyle/ginbase/consolelog"
)

// process ca user register and enroll
// save ca credentials into file fabricwallet

const (
	FabricCAUserTypeUser    = "client"
	FabricCAUserTypeAdmin   = "admin"
	FabricCAUserTypePeer    = "peer"
	FabricCAUserTypeOrderer = "orderer"
)

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
	wallet, err := NewWallet(ctx)
	if err != nil {
		return nil, err
	}
	identity := gateway.WithIdentity(wallet, userId)
	gw, err := gateway.Connect(gwCfg, identity)
	if err != nil {
		Logger.Errorf("", "create fabric gateway failed, userId[%s], reason: %s", userId, err.Error())
		return nil, err
	}
	Logger.Debug("", "create fabric gateway object success")
	return gw, nil
}

func EnrollUser(ctx *context.ApiContext, userId, passwd string) error {
	// first check if wallet exist
	wallet, err := NewWallet(ctx)
	if err != nil {
		return err
	}
	if wallet.Exists(userId) {
		Logger.Infof("", "EnrollUser, userId[%s] has exists", userId)
		return nil
	}

	cfgFile := ctx.Cfg.Fabric.CCPath
	sdk, err := fabsdk.New(config.FromFile(cfgFile))
	if err != nil {
		Logger.Errorf("", "fabsdk.New() failed, %s", err.Error())
		return err
	}
	defer sdk.Close()

	curOrg := ctx.Cfg.Fabric.OrgName
	client, err := msp.New(sdk.Context(), msp.WithOrg(curOrg))
	if err != nil {
		Logger.Errorf("", "msp.New() failed, %s", err.Error())
		return err
	}

	eOpt := msp.WithSecret(passwd)
	err = client.Enroll(userId, eOpt)
	if err != nil {
		Logger.Errorf("", "client.Enroll failed, %s", err.Error())
		return err
	}

	si, err := client.GetSigningIdentity(userId)
	if err != nil {
		Logger.Errorf("", "client.GetSigningIdentity() failed, %s", err.Error())
		return err
	}

	publicKey := si.EnrollmentCertificate()
	privateKey, err := si.PrivateKey().Bytes()
	if err != nil {
		Logger.Errorf("", "get private key failed, %s", err.Error())
		return err
	}

	newIdentity := gateway.NewX509Identity(si.PublicVersion().Identifier().MSPID, string(publicKey), string(privateKey))
	err = wallet.Put(userId, newIdentity)
	if err != nil {
		Logger.Errorf("", "wallet.put(%s) failed, %s", userId, err.Error())
	}
	Logger.Infof("", "MSPID: %s, ID: %s", si.PublicVersion().Identifier().MSPID, si.PublicVersion().Identifier().ID)
	Logger.Infof("", "EnrollUser, save userId[%s] into wallet", userId)

	return nil
}

func RegisterCaUser(ctx *context.ApiContext, userId, passwd, userType string) error {
	// first check if wallet exist
	wallet, err := NewWallet(ctx)
	if err != nil {
		return err
	}
	if wallet.Exists(userId) {
		Logger.Infof("", "RegisterCaUser, userId[%s] has exists", userId)
		return nil
	}

	cfgFile := ctx.Cfg.Fabric.CCPath
	sdk, err := fabsdk.New(config.FromFile(cfgFile))
	if err != nil {
		Logger.Errorf("", "fabsdk.New() failed, %s", err.Error())
		return err
	}
	defer sdk.Close()

	curOrg := ctx.Cfg.Fabric.OrgName
	client, err := msp.New(sdk.Context(), msp.WithOrg(curOrg))
	if err != nil {
		Logger.Errorf("", "msp.New() failed, %s", err.Error())
		return err
	}

	regForm := &msp.RegistrationRequest{
		Name:           userId,
		Type:           userType,
		MaxEnrollments: -1,
		Secret:         passwd,
	}
	retPasswd, err := client.Register(regForm)
	if err != nil {
		Logger.Errorf("", "client.Register() failed, %s", err.Error())
		return err
	}
	Logger.Infof("", "reg new identity[%s][%s] success", userId, retPasswd)

	err = client.Enroll(userId, msp.WithSecret(retPasswd))
	if err != nil {
		Logger.Errorf("", "client.Enroll %s failed, %s", userId, err.Error())
		return err
	}

	si, err := client.GetSigningIdentity(userId)
	if err != nil {
		Logger.Errorf("", "client.GetSigningIdentity() failed, %s", err.Error())
		return err
	}

	publicKey := si.EnrollmentCertificate()
	privateKey, err := si.PrivateKey().Bytes()
	if err != nil {
		Logger.Errorf("", "get private key failed, %s", err.Error())
		return err
	}

	newIdentity := gateway.NewX509Identity(si.PublicVersion().Identifier().MSPID, string(publicKey), string(privateKey))
	err = wallet.Put(userId, newIdentity)
	if err != nil {
		Logger.Errorf("", "wallet.put(%s) failed, %s", userId, err.Error())
	}
	Logger.Infof("", "MSPID: %s, ID: %s", si.PublicVersion().Identifier().MSPID, si.PublicVersion().Identifier().ID)
	Logger.Infof("", "RegisterUser, save userId[%s] into wallet", userId)

	return nil
}
