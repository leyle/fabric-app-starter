package model

import "encoding/json"

// chaincode method names
// see chaincode/public/chaincode.go and chaincode/private/chaincode.go
const (
	CCCreatePublicState   = "CreatePublic"
	CCCreatePrivateState  = "CreatePrivate"
	CCGetPublicStateById  = "GetPublicStateById"
	CCGetPrivateStateById = "GetPrivateStateById"
)

const (
	TransientActionAdd    = "CREATE"
	TransientActionDelete = "DELETE"
	TransientActionUpdate = "UPDATE"
)

// call chaincode's wrap function response
type CCApiResponse struct {
	Error error

	// input dataId, return as original
	DataId string `json:"dataId"`

	// get/query response
	CCRet []byte `json:"ccRet"`
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
