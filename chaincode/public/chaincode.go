package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/util"
)

// universal public chaincode
// use couchdb as world state database
// implement CRUD

var ErrNoIdData = errors.New("no id data")

type UniversalContract struct {
	contractapi.Contract
}

// save in ledger
type StorageIn struct {
	Id        string          `json:"id"`        // system db id = channel + chaincode + dataId
	App       string          `json:"app"`       // client app name
	DataId    string          `json:"dataId"`    // client data unique id
	Data      json.RawMessage `json:"data"`      // user input data, json encoded
	CreatedAt int64           `json:"createdAt"` // system timestamp
}

// query from ledger, return to caller
type StorageOut struct {
	Id        string `json:"id"` // system db id = channel + chaincode + dataId
	App       string `json:"app"`
	DataId    string `json:"dataId"`
	Data      string `json:"data"`      // user input data, json encoded
	CreatedAt int64  `json:"createdAt"` // system timestamp
}

// create
func (uc *UniversalContract) Create(ctx contractapi.TransactionContextInterface, app, dataId, data string) error {
	// check if data has exist
	id := getId(app, dataId)
	exist, err := uc.IsExist(ctx, id)
	if err != nil {
		return err
	}
	if exist {
		emsg := fmt.Errorf("data[%s] already exist", id)
		Logger.Error("", emsg.Error())
		return emsg
	}

	storage := &StorageIn{
		Id:        id,
		App:       app,
		DataId:    dataId,
		Data:      []byte(data),
		CreatedAt: util.CurUnixTime(),
	}

	storageJson, err := json.Marshal(storage)
	if err != nil {
		Logger.Errorf("", "marshar storage failed, %s", err.Error())
		return err
	}

	err = ctx.GetStub().PutState(id, storageJson)
	if err != nil {
		Logger.Errorf("", "save data[%s] into ledger failed, %s", id, err.Error())
	}
	return err
}

func (uc *UniversalContract) Update(ctx contractapi.TransactionContextInterface, app, dataId, data string) error {
	id := getId(app, dataId)
	exist, err := uc.IsExist(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		emsg := fmt.Errorf("data[%s] doesn't exist", id)
		Logger.Error("", emsg.Error())
		return emsg
	}

	storage := &StorageIn{
		Id:        id,
		App:       app,
		DataId:    dataId,
		Data:      []byte(data),
		CreatedAt: util.CurUnixTime(),
	}

	storageJson, err := json.Marshal(storage)
	if err != nil {
		Logger.Errorf("", "marshar storage failed, %s", err.Error())
		return err
	}

	err = ctx.GetStub().PutState(id, storageJson)
	if err != nil {
		Logger.Errorf("", "update data[%s] into ledger failed, %s", id, err.Error())
	}
	return err
}

func (uc *UniversalContract) DeleteById(ctx contractapi.TransactionContextInterface, app, dataId string) error {
	id := getId(app, dataId)
	exist, err := uc.IsExist(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		emsg := fmt.Errorf("data[%s] doesn't exist", id)
		Logger.Error("", emsg.Error())
		return emsg
	}

	err = ctx.GetStub().DelState(id)
	if err != nil {
		Logger.Errorf("", "delete data[%s] from ledger failed, %s", id, err.Error())
	}
	return err
}

// get by id
func (uc *UniversalContract) GetById(ctx contractapi.TransactionContextInterface, app, dataId string) (*StorageOut, error) {
	id := getId(app, dataId)
	dataJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		Logger.Errorf("", "GetById[%s] failed, %s", id, err.Error())
		return nil, err
	}
	if dataJson == nil {
		Logger.Errorf("", "GetById[%s], data not exist", id)
		return nil, ErrNoIdData
	}
	var storage StorageIn
	err = json.Unmarshal(dataJson, &storage)
	if err != nil {
		Logger.Errorf("", "GetById[%s], unmarshal failed, %s", id, err.Error())
		return nil, err
	}

	storageOut := copyToStorageOut(&storage)

	return storageOut, nil
}

// get data by start and end key
func (uc *UniversalContract) QueryByRange(ctx contractapi.TransactionContextInterface, start, end string) ([]*StorageOut, error) {
	return nil, nil
}

// rich query
func (uc *UniversalContract) Search(ctx contractapi.TransactionContextInterface, filter string, size int, bookmark string) ([]*StorageOut, error) {
	return nil, nil
}

func (uc *UniversalContract) IsExist(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	jdata, err := ctx.GetStub().GetState(id)
	if err != nil {
		Logger.Errorf("", "Check id[%s] exist failed, %s", id, err.Error())
		return false, err
	}

	if jdata == nil {
		return false, nil
	}
	return true, nil
}
