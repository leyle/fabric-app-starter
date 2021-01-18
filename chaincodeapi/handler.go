package chaincodeapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/chaincodeapi/model"
	"github.com/leyle/fabric-app-starter/chaincodeapi/private"
	"github.com/leyle/fabric-app-starter/chaincodeapi/public"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/go-api-starter/couchdb"
	"github.com/leyle/go-api-starter/ginhelper"
)

// create handler
type CreateForm struct {
	// application name
	App string `json:"app" binding:"required"`

	// dataId, it should be unique in entire applications
	DataId string `json:"dataId" binding:"required"`

	// public data
	Public *public.CreatePublicForm `json:"public"`

	// private data
	Private *private.CreatePrivateForm `json:"private"`
}

func CreateHandler(ctx *context.ApiContext) {
	c := ctx.C
	var form CreateForm
	err := c.BindJSON(&form)
	ginhelper.StopExec(err)

	apiResp := &model.ApiResponse{
		Success: model.ResponseSuccessNone,
		App:     form.App,
		DataId:  form.DataId,
	}

	// public data call public chaincode
	// private data call private chaincode

	// because private data require more permission control
	// so we create it first, if failed, then whole create failed
	// if create private data success, then create public data

	// 1. check if we need to create private data
	privateForm := form.Private
	if privateForm != nil {
		privateForm.App = form.App
		privateForm.DataId = form.DataId
		privateResp := private.CallPrivateChainCode(ctx, privateForm)
		if privateResp.Error != nil {
			apiResp.ErrMsg = privateResp.Error.Error()
			ginhelper.ReturnJson(c, 400, 400, "", apiResp)
			return
		}
		ctx.Logger().Info().Msg("create private data success")
		apiResp.Success = model.ResponseSuccessPartial
	}

	// 2. check if we need to create public data
	publicForm := form.Public
	if publicForm != nil {
		publicForm.App = form.App
		publicForm.DataId = form.DataId
		resp := public.CallPublicChainCodeCreate(ctx, publicForm)
		if resp.Error != nil {
			apiResp.ErrMsg = resp.Error.Error()
			ginhelper.ReturnJson(c, 400, 400, "", apiResp)
			return
		}
		apiResp.Success = model.ResponseSuccessPartial
	}

	apiResp.Success = model.ResponseSuccessAll
	ginhelper.ReturnOKJson(c, apiResp)
}

func GetByIdHandler(ctx *context.ApiContext) {
	var form model.GetByIdForm
	err := ctx.C.BindJSON(&form)
	ginhelper.StopExec(err)

	apiResp := &model.ApiResponse{
		Success: model.ResponseSuccessNone,
		App:     form.App,
		DataId:  form.DataId,
	}

	if form.Public == nil && form.Private == nil {
		ctx.Logger().Error().Msg("getbyid, but no public and private query args")
		ginhelper.ReturnErrJson(ctx.C, "No public or private query args")
		return
	}

	if form.Public != nil {
		resp := public.CallPublicChaincodeGetById(ctx, &form)
		if resp.Error != nil {
			apiResp.ErrMsg = resp.Error.Error()
			ginhelper.ReturnJson(ctx.C, 400, 400, "", apiResp)
			return
		}
		apiResp.PublicCCResp = resp.CCRet
		apiResp.Success = model.ResponseSuccessPartial
	}
	if form.Private != nil {
		presp := private.CallPrivateChaincodeGetById(ctx, &form)
		if presp.Error != nil {
			apiResp.ErrMsg = presp.Error.Error()
			ginhelper.ReturnJson(ctx.C, 400, 400, "", apiResp)
			return
		}
		apiResp.PrivateCCResp = presp.CCRet
		apiResp.Success = model.ResponseSuccessPartial
	}

	apiResp.Success = model.ResponseSuccessAll
	ginhelper.ReturnOKJson(ctx.C, apiResp)
}

// search handler
type CCSearchForm struct {
	Channel   string `json:"channel" binding:"required"`
	Chaincode string `json:"chaincode" binding:"required"`
	// if collectionName is not null, we search private chaincode database
	CollectionName string `json:"collectionName"`

	Selector interface{} `json:"selector" binding:"required"`
	Sort     interface{} `json:"sort,omitempty"`
	Page     int         `json:"page"`
	Size     int         `json:"size"`
}

func SearchHandler(ctx *context.ApiContext) {
	var form CCSearchForm
	err := ctx.C.BindJSON(&form)
	ginhelper.StopExec(err)

	dbName := helper.GenerateCouchdbDatabaseName(form.Channel, form.Chaincode, form.CollectionName)
	ds := helper.GetCouchdbClient(ctx, dbName)

	limit := form.Size
	page := form.Page
	if page <= 1 {
		page = 1
	}
	skip := (page - 1) * form.Size

	searchReq := &couchdb.SearchRequest{
		Selector: form.Selector,
		Sort:     form.Sort,
		Limit:    limit,
		Skip:     skip,
	}

	resp, err := ds.Search(ctx.C.Request.Context(), searchReq, nil)
	if err != nil {
		ctx.Logger().Error().Err(err).Msg("search couchdb failed")
		ginhelper.ReturnErrJson(ctx.C, err.Error())
		return
	}

	docs := resp.Docs
	bookmark := resp.Bookmark
	ctx.Logger().Info().Str("bookmark", bookmark).Int("page", page).Int("size", form.Size).Send()

	retData := gin.H{
		"page": page,
		"size": form.Size,
		"data": docs,
	}

	ginhelper.ReturnOKJson(ctx.C, retData)
	return
}

type CreateStateForm struct {
	DataId  string                  `json:"dataId"`
	AppName string                  `json:"appName"`
	Public  *CreatePublicStateForm  `json:"public"`
	Private *CreatePrivateStateForm `json:"private"`
}

type CreatePublicStateForm struct {
	Channel   string `json:"channel" binding:"required"`
	Chaincode string `json:"chaincode" binding:"required"`
	DataJson  string `json:"dataJson" binding:"required"`
}

type CreatePrivateStateForm struct {
	Channel     string              `json:"channel" binding:"required"`
	Chaincode   string              `json:"chaincode" binding:"required"`
	Collections *CollectionNameForm `json:"collections" binding:"required"`
}

type CollectionNameForm struct {
	Self  string   `json:"self"`
	Share []string `json:"share"`
}

func CreateStateHandler(ctx *context.ApiContext) {
	var form CreateStateForm
	err := ctx.C.BindJSON(&form)
	ginhelper.StopExec(err)

	apiResp := model.NewApiResponse(form.AppName, form.DataId)

	// if has private data, create if first
	if form.Private != nil {
		fmt.Println(apiResp)
	}

}
