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

// can return both public and private data
const (
	ResponseSuccessAll     = "all"
	ResponseSuccessPartial = "partial"
	ResponseSuccessNone    = "none"
)

type ApiResponse struct {
	// return input data
	Success string `json:"success"`
	ErrMsg  string `json:"errMsg"`

	PublicCCResp  json.RawMessage `json:"public"`
	PrivateCCResp json.RawMessage `json:"private"`
	App           string          `json:"app"`
	DataId        string          `json:"dataId"`
}

func NewApiResponse(app, id string) *ApiResponse {
	return &ApiResponse{
		Success: ResponseSuccessNone,
		App:     app,
		DataId:  id,
	}
}

func NewCCApiResponse() *CCApiResponse {
	return &CCApiResponse{}
}

type GetByIdForm struct {
	// application name
	App string `json:"app" binding:"required"`

	// dataId, it should be unique in entire applications
	DataId string `json:"dataId" binding:"required"`

	Public  *GetByIdPublicForm  `json:"public"`
	Private *GetByIdPrivateForm `json:"private"`
}

type GetByIdPublicForm struct {
	Channel   string `json:"channel" binding:"required"`
	ChainCode string `json:"chaincode" binding:"required"`
}

type GetByIdPrivateForm struct {
	Channel        string `json:"channel" binding:"required"`
	ChainCode      string `json:"chaincode" binding:"required"`
	CollectionName string `json:"collectionName"`
}
