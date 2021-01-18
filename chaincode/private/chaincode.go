package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

var ErrNoIdData = errors.New("no id data")

type UniversalContract struct {
	contractapi.Contract
}

// input transient data
type TransientInput struct {
	CollectionName string `json:"collectionName"`
	App            string `json:"app"`    // client app name
	DataId         string `json:"dataId"` // client data unique id
	Data           string `json:"data"`   // user input data, json encoded
}

// data store in private data
type StorageIn struct {
	Id             string          `json:"id"` // system db id = channel + chaincode + dataId
	CollectionName string          `json:"collectionName"`
	App            string          `json:"app"`       // client app name
	DataId         string          `json:"dataId"`    // client data unique id
	Data           json.RawMessage `json:"data"`      // user input data, json encoded
	CreatedAt      int64           `json:"createdAt"` // system timestamp
}

// query from ledger, return to caller
type StorageOut struct {
	Id             string `json:"id"` // id = channel + chaincode + dataId
	CollectionName string `json:"collectionName"`
	App            string `json:"app"`
	DataId         string `json:"dataId"`
	Data           string `json:"data"`      // user input data, json encoded
	CreatedAt      int64  `json:"createdAt"` // system timestamp
}

// key is the transient map key
// a transient struct is map[string][]byte
func (uc *UniversalContract) Create(ctx contractapi.TransactionContextInterface, key string) error {
	tmap, err := ctx.GetStub().GetTransient()
	if err != nil {
		err := fmt.Errorf("failed to get transient, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	// get private data(transient data)
	transientJson, ok := tmap[key]
	if !ok {
		err = fmt.Errorf("no key[%s] in transient map input", key)
		fmt.Println(err.Error())
		return err
	}

	// unmarshal private data to TransientInput
	var inputData TransientInput
	err = json.Unmarshal(transientJson, &inputData)
	if err != nil {
		err = fmt.Errorf("unmarshal input transient data failed, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	// input data validate action move to chaincode api
	// so here we pass the validation

	// check if dataId value already exists
	dataId := getId(inputData.App, inputData.DataId)
	exist, err := uc.IsExist(ctx, inputData.CollectionName, dataId)
	if err != nil {
		return err
	}
	if exist {
		err = fmt.Errorf("create private data, but id[%s] in collection[%s] already exists", inputData.DataId, inputData.CollectionName)
		fmt.Println(err.Error())
		return err
	}

	// check submitter submits data from its own peer
	/*
		err = IsSubmitterFromSameOrg(ctx)
		if err != nil {
			return err
		}

	*/

	// save it
	storageIn := &StorageIn{
		Id:             dataId,
		CollectionName: inputData.CollectionName,
		App:            inputData.App,
		DataId:         inputData.DataId,
		Data:           []byte(inputData.Data),
		CreatedAt:      time.Now().Unix(),
	}

	storageJson, err := json.Marshal(storageIn)
	if err != nil {
		err = fmt.Errorf("marshal private storage in data failed, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	err = ctx.GetStub().PutPrivateData(storageIn.CollectionName, storageIn.Id, storageJson)
	if err != nil {
		err = fmt.Errorf("failed to put private data collection[%s], dataId[%s]", storageIn.CollectionName, storageIn.DataId)
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (uc *UniversalContract) GetById(ctx contractapi.TransactionContextInterface, collectionName, app, dataId string) (*StorageOut, error) {
	id := getId(app, dataId)
	dataJson, err := ctx.GetStub().GetPrivateData(collectionName, id)
	if err != nil {
		err = fmt.Errorf("get collection[%s] id[%s] private data failed, %s", collectionName, dataId, err.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	if dataJson == nil {
		err = fmt.Errorf("get collection[%s] id[%s] private data failed, %s", collectionName, dataId, ErrNoIdData.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	var storage *StorageIn
	err = json.Unmarshal(dataJson, &storage)
	if err != nil {
		err = fmt.Errorf("get collection[%s] id[%s] private data, unmarshal failed, %s", collectionName, dataId, err.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	storageOut := copyToStorageOut(storage)
	return storageOut, nil
}

func (uc *UniversalContract) IsExist(ctx contractapi.TransactionContextInterface, collectionName, id string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateData(collectionName, id)
	if err != nil {
		err = fmt.Errorf("get private data by id[%s] from collection[%s] failed, %s", id, collectionName, err.Error())
		fmt.Println(err.Error())
		return false, err
	}

	if data != nil {
		return true, nil
	}
	return false, nil
}

func GetSubmitterIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
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

func IsSubmitterFromSameOrg(ctx contractapi.TransactionContextInterface) error {
	submitterMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		err = fmt.Errorf("get submitter's mspid failed, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		err = fmt.Errorf("get peer's mspid failed, %s", err.Error())
		fmt.Println(err.Error())
		return err
	}

	if submitterMSPID != peerMSPID {
		err = fmt.Errorf("submitter's mspid[%s] is not equal peer's mspid[%s]", submitterMSPID, peerMSPID)
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func getId(app, dataId string) string {
	return fmt.Sprintf("%s|%s", app, dataId)
}

func copyToStorageOut(in *StorageIn) *StorageOut {
	if in == nil {
		return nil
	}
	out := &StorageOut{
		Id:             in.Id,
		CollectionName: in.CollectionName,
		App:            in.App,
		DataId:         in.DataId,
		Data:           string(in.Data),
		CreatedAt:      in.CreatedAt,
	}

	return out
}
