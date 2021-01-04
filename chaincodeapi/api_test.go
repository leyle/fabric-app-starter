package chaincodeapi

import (
	"context"
	"encoding/json"
	"github.com/leyle/fabric-app-starter/chaincodeapi/private"
	"github.com/leyle/fabric-app-starter/chaincodeapi/public"
	"github.com/leyle/go-api-starter/httpclient"
	"github.com/leyle/go-api-starter/logmiddleware"
	"github.com/leyle/go-api-starter/util"
	"testing"
	"time"
)

// suppose api server listening on localhost:9000
const (
	host             = "http://localhost:9000"
	publicChannel    = "one"
	publicChaincode  = "p2"
	privateChannel   = "ppp"
	privateChaincode = "ptwo"
	collectionName   = "org1private"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI1ZmYxZjk2MTI2MzEzM2JhZGEyNGJjODYiLCJ1c2VybmFtZSI6ImRldnRlc3QiLCJyb2xlIjoiY2xpZW50IiwiZXhwIjoxNjEwNTU4MDg5LCJpYXQiOjE2MDk2OTQwODl9.2HPAJnwpCJq6xFeRW9pABkx40IIcTSYayRKwSDPOM7E"

var headers = map[string]string{
	"X-TOKEN": token,
}

type PublicDataFormat struct {
	Id               string     `json:"id"`
	Type             string     `json:"type"`
	DataConsumerId   string     `json:"dataConsumerId"`
	DataConsumerName string     `json:"dataConsumerName"`
	DataProviderId   string     `json:"dataProviderId"`
	DataProviderName string     `json:"dataProviderName"`
	Scope            *DataScope `json:"scope"`
	Purpose          string     `json:"purpose"`
	ExpireTime       int64      `json:"expireTime"`
}

type PrivateDataFormat struct {
	DataOwnerId             string                     `json:"dataOwnerId"`
	DataOwnerIdType         string                     `json:"dataOwnerIdType"`
	DataOwnerName           string                     `json:"dataOwnerName"`
	PurposePrivate          string                     `json:"purposePrivate"`
	CustomerConsentFilelist []*CustomerConsentFileList `json:"customerConsentFileList"`
}

type DataScope struct {
	ScopeType          string `json:"scopeType"`
	DocumentType       string `json:"documentType"`
	UDR                string `json:"udr"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	QueryFrequencyType string `json:"queryFrequency"`
	AccessMode         string `json:"accessMode"`
}

type CustomerConsentFileList struct {
	Content    string `json:"content"`
	Type       string `json:"type"`
	FileFormat string `json:"fileFormat"`
}

func TestCreateHandler(t *testing.T) {
	createApi := host + "/api/chaincode/publicandprivate/create"
	dataId := util.GenerateDataId()

	scope := &DataScope{
		ScopeType:          "DataScopeTypeUDR",
		DocumentType:       "document type",
		UDR:                "udr",
		StartDate:          "2020-12-01",
		EndDate:            "2021-11-30",
		QueryFrequencyType: "DAILY",
		AccessMode:         "View",
	}
	reqPublic := &PublicDataFormat{
		Id:               dataId,
		Type:             "CustomerConsentTypeSignedDocRequired",
		DataConsumerId:   "FDI_PARTICIPANT_00002",
		DataConsumerName: "HKBEA",
		DataProviderId:   "FDI_PARTICIPANT_00001",
		DataProviderName: "TradeLink",
		Scope:            scope,
		Purpose:          "purpose value",
		ExpireTime:       time.Now().Unix() + 10*24*60*60,
	}

	reqJson, _ := json.Marshal(reqPublic)

	publicData := &public.CreatePublicForm{
		Channel:   publicChannel,
		ChainCode: publicChaincode,
		DataJson:  string(reqJson),
	}

	fileList := &CustomerConsentFileList{
		Content:    "file list content",
		Type:       "pdf",
		FileFormat: "pdf",
	}
	fileList2 := &CustomerConsentFileList{
		Content:    "file list content2",
		Type:       "png",
		FileFormat: "png",
	}

	reqPrivate := &PrivateDataFormat{
		DataOwnerId:             "DO_001",
		DataOwnerIdType:         "DO_ID_TYPE",
		DataOwnerName:           "DO_OWNER_NAME",
		PurposePrivate:          "private purpose value",
		CustomerConsentFilelist: []*CustomerConsentFileList{fileList, fileList2},
	}
	reqPrivateJson, _ := json.Marshal(reqPrivate)

	privateData := &private.CreatePrivateForm{
		Channel:        privateChannel,
		ChainCode:      privateChaincode,
		CollectionName: collectionName,
		DataJson:       string(reqPrivateJson),
	}

	reqData := &CreateForm{
		App:     "customerConsent",
		DataId:  util.GenerateDataId(),
		Public:  publicData,
		Private: privateData,
	}

	reqBody, _ := json.Marshal(reqData)

	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)
	ctx := context.Background()
	ctx = logger.WithContext(ctx)
	cReq := &httpclient.ClientRequest{
		Ctx:     ctx,
		Url:     createApi,
		Headers: headers,
		Body:    reqBody,
		Debug:   true,
	}

	resp := httpclient.Post(cReq)
	if resp.Err != nil {
		t.Fatal(resp.Err)
	}

	t.Log(resp.Code, string(resp.Body))
}
