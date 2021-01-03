package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	cc, err := contractapi.NewChaincode(&UniversalContract{})
	if err != nil {
		fmt.Println("create new private chaincode failed", err.Error())
		return
	}

	err = cc.Start()
	if err != nil {
		fmt.Println("start private chaincode failed", err.Error())
		return
	}
}
