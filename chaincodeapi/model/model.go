package model

import "encoding/json"

// chaincode method names
// see chaincode/public/chaincode.go and chaincode/private/chaincode.go
const (
	CCNameCreate  = "Create"
	CCNameUpdate  = "Update"
	CCNameGetById = "GetById"
)

// call chaincode's wrap function response
type CCApiResponse struct {
	Error error

	// input dataId, return as original
	DataId string `json:"dataId"`

	// chaincode data row key
	CCRet []byte `json:"ccRet"`

	CCId string `json:"ccId"`
}

type ApiResponse struct {
	// return input data
	Success bool   `json:"success"`
	ErrMsg  string `json:"errMsg"`

	CCResp json.RawMessage `json:"response"`
	App    string          `json:"app"`
	DataId string          `json:"dataId"`
}

// error type todo
type CCError struct {
}

func NewCCApiResponse() *CCApiResponse {
	return &CCApiResponse{}
}

type GetByIdForm struct {
	// application name
	App string `json:"app" binding:"required"`

	// dataId, it should be unique in entire applications
	DataId string `json:"dataId" binding:"required"`

	// channel and chaincode
	Channel   string `json:"channel" binding:"required"`
	ChainCode string `json:"chaincode" binding:"required"`
}
