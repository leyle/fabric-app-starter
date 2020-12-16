package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	. "github.com/leyle/ginbase/consolelog"
)

func main() {
	cc, err := contractapi.NewChaincode(&UniversalContract{})
	if err != nil {
		Logger.Errorf("", "NewChaincode Failed, %s", err.Error())
		return
	}

	err = cc.Start()
	if err != nil {
		Logger.Errorf("", "start chaincode failed, %s", err.Error())
		return
	}
}
