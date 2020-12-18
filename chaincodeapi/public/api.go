package public

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/authproxy"
	"github.com/leyle/fabric-app-starter/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/fabric-app-starter/fabricwallet"
	"github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/middleware"
)

// call public chaincode
type PublicForm struct {
	App    string `json:"-"`
	DataId string `json:"-"`

	Channel   string `json:"channel" binding:"required"`
	ChainCode string `json:"chaincode" binding:"required"`

	DataJson string `json:"dataJson" binding:"required"`
}

func CallPublicChainCode(ctx *context.ApiContext, c *gin.Context, form *PublicForm) *model.CCApiResponse {
	reqId := middleware.GetReqId(c)
	ccResp := &model.CCApiResponse{
		Error:  nil,
		DataId: form.DataId,
	}
	// get fabric gateway
	curUser := authproxy.GetCurUser(c)
	gw, err := fabricwallet.NewGateway(ctx, curUser.UserId)
	if err != nil {
		ccResp.Error = err
		return ccResp
	}
	defer gw.Close()
	// gw := authproxy.GetGateway(c)
	// if gw == nil {
	// 	err := errors.New("get fabric gateway failed")
	// 	consolelog.Logger.Errorf(reqId, "call public chaincode, %s", err.Error())
	// 	ccResp.Error = err
	// 	return ccResp
	// }

	// get network
	network, err := gw.GetNetwork(form.Channel)
	if err != nil {
		consolelog.Logger.Errorf(reqId, "call public chaincode, get channel failed, %s", err.Error())
		ccResp.Error = err
		return ccResp
	}

	// get contract
	contract := network.GetContract(form.ChainCode)

	// call chaincode's create method
	// chaincode's create method locates on chaincode/public/chaincode.go#Create
	// func (uc *UniversalContract) Create(ctx contractapi.TransactionContextInterface, app, dataId, data string) error
	ret, err := contract.SubmitTransaction(model.CCNameCreate, form.App, form.DataId, form.DataJson)
	if err != nil {
		consolelog.Logger.Errorf(reqId, "call pubic chaincode, Create failed, %s", err.Error())
		ccResp.Error = err
		return ccResp
	}
	consolelog.Logger.Infof(reqId, "call public chaincode, Create success, result is: %s", string(ret))
	ccResp.CCRet = ret

	return ccResp
}
