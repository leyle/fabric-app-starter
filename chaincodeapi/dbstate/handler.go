package dbstate

import (
	"github.com/leyle/fabric-app-starter/chaincodeapi/helper"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/go-api-starter/ginhelper"
	"strings"
)

// create couchdb index
type CreateCouchdbIndexForm struct {
	Channel   string `json:"channel" binding:"required"`
	Chaincode string `json:"chaincode" binding:"required"`

	// if collectionName has value, this is a private data
	CollectionName string `json:"collectionName"`

	// index key
	Fields []string `json:"fields" binding:"required"`
}

func CreateCouchdbIndexHandler(ctx *context.ApiContext) {
	var form CreateCouchdbIndexForm
	err := ctx.C.BindJSON(&form)
	ginhelper.StopExec(err)

	for idx, v := range form.Fields {
		form.Fields[idx] = strings.TrimSpace(v)
	}

	dbName := helper.GenerateCouchdbDatabaseName(form.Channel, form.Chaincode, form.CollectionName)
	ds := helper.GetCouchdbClient(ctx, dbName)

	keys := generateIndexKey(form.Fields)
	err = ds.CreateIndex(ctx.C.Request.Context(), keys)
	if err != nil {
		ctx.Logger().Error().Err(err).Str("channel", form.Channel).Str("chaincode", form.Chaincode).Msg("create index failed")
		ginhelper.ReturnErrJson(ctx.C, err.Error())
		return
	}

	ginhelper.ReturnOKJson(ctx.C, "")
	return
}

func generateIndexKey(fields []string) []string {
	var keys []string
	for _, val := range fields {
		val = "data." + val
		keys = append(keys, val)
	}
	keys = append(keys, "createdAt")
	return keys
}
