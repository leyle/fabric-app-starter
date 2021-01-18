package helper

import (
	"encoding/base64"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func GetClientID(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		err = fmt.Errorf("get client identitiy id failed, %s", err.Error())
		fmt.Println(err.Error())
		return "", err
	}

	decodeId, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		err = fmt.Errorf("get client identity id, decode b64 data failed, %s", err.Error())
		fmt.Println(err.Error())
		return "", err
	}

	return string(decodeId), nil
}

func GetClientMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		err = fmt.Errorf("get submitter's mspid failed, %s", err.Error())
		fmt.Println(err.Error())
		return "", err
	}
	return mspID, nil
}

func GetPeerMSPID() (string, error) {
	mspID, err := shim.GetMSPID()
	if err != nil {
		err = fmt.Errorf("get peer's mspid failed, %s", err.Error())
		fmt.Println(err.Error())
		return "", err
	}
	return mspID, nil
}
