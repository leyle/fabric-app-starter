package contract

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/leyle/fabric-app-starter/chaincode/universal/model"
	"time"
)

func (pc *Contract) CreatePrivate(ctx contractapi.TransactionContextInterface, transientKey string) error {
	var err error
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		err = fmt.Errorf("failed to get transient, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	// get private data from transient map
	tJson, ok := transientMap[transientKey]
	if !ok {
		err = fmt.Errorf("no key[%s] in transient map input", transientKey)
		fmt.Println(err.Error())
		return err
	}

	var inputData model.TransientForm
	err = json.Unmarshal(tJson, &inputData)
	if err != nil {
		err = fmt.Errorf("unmarshal input transient data failed, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	privateState := &model.PrivateState{
		Id:           inputData.Id,
		PublicDataId: inputData.PublicDataId,
		Data:         []byte(inputData.Data),
		CreatedAt:    time.Now().Unix(),
	}

	stateJson, err := json.Marshal(privateState)
	if err != nil {
		err = fmt.Errorf("marshal private storage in data failed, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	err = ctx.GetStub().PutPrivateData(inputData.CollectionName, privateState.Id, stateJson)
	if err != nil {
		err = fmt.Errorf("failed to put private data collection[%s], dataId[%s]", inputData.CollectionName, privateState.Id)
		fmt.Println(err.Error())
		return err
	}

	return nil
}
