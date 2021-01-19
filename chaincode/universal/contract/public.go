package contract

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/leyle/fabric-app-starter/chaincode/universal/helper"
	"github.com/leyle/fabric-app-starter/chaincode/universal/ledgerstate"
	"time"
)

func (pc *Contract) CreatePublic(ctx contractapi.TransactionContextInterface, args string) error {
	// unmarshal args to public input form
	var form ledgerstate.PublicStateForm
	err := json.Unmarshal([]byte(args), &form)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// check if data id has existed
	exists, err := pc.IsExists(ctx, form.Id)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("create public state, but data id[%s] has already existed", form.Id)
		fmt.Println(err.Error())
		return err
	}

	creatorMSPId, _ := helper.GetClientMSPID(ctx)
	creatorId, _ := helper.GetClientID(ctx)

	// save into public state ledger
	publicState := &ledgerstate.PublicState{
		Id:              form.Id,
		AppName:         form.AppName,
		CreatorMSPID:    creatorMSPId,
		CreatorId:       creatorId,
		DataJson:        []byte(form.Data),
		PrivateMetaInfo: form.PrivateMetaInfo,
		CreatedAt:       time.Now().Unix(),
	}
	publicState.UpdatedAt = publicState.CreatedAt

	stateJson, err := json.Marshal(publicState)
	if err != nil {
		fmt.Println("marshal public data failed,", err.Error())
		return err
	}

	err = ctx.GetStub().PutState(form.Id, stateJson)
	return err
}

func (pc *Contract) IsExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		fmt.Println("get public state by id failed", err.Error())
		return false, err
	}
	if state == nil {
		return false, nil
	}
	return true, nil
}

func (pc *Contract) GetPublicStateById(ctx contractapi.TransactionContextInterface, id string) (*ledgerstate.PublicState, error) {
	state, err := ctx.GetStub().GetState(id)
	if err != nil {
		fmt.Println("get public state by id failed,", id, err.Error())
		return nil, err
	}
	if state == nil {
		fmt.Println("get public state by id, data doesn't exists", id)
		return nil, ErrNoIdData
	}

	var publicState ledgerstate.PublicState
	err = json.Unmarshal(state, &publicState)
	if err != nil {
		fmt.Println("get public state by id, unmarshal failed,", id, err.Error())
		return nil, err
	}

	reqData := publicState.DataJson
	publicState.DataString = string(reqData)
	publicState.DataJson = nil

	return &publicState, nil
}
