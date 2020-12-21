package model

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
	Result bool `json:"result"`

	CCResp string `json:"ccResp"`
	App    string `json:"app"`
	DataId string `json:"dataId"`
}

// error type todo
type CCError struct {
}
