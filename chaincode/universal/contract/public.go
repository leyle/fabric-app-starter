package contract

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/leyle/fabric-app-starter/chaincode/universal/helper"
	"github.com/leyle/fabric-app-starter/chaincode/universal/model"
	"time"
)

func (pc *Contract) CreatePublic(ctx contractapi.TransactionContextInterface, args string) error {
	// unmarshal args to public input form
	var form model.PublicStateForm
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
	publicState := &model.PublicState{
		Id:              form.Id,
		AppName:         form.AppName,
		CreatorMSPID:    creatorMSPId,
		CreatorId:       creatorId,
		Data:            []byte(form.Data),
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
