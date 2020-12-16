package publicandprivateapi

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/ginbase/middleware"
)

// create handler
type CreateForm struct {
	// application name
	App string `json:"app" binding:"required"`

	// dataId, it should be unique in entire applications
	DataId string `json:"dataId" binding:"required"`

	// public data
	Public *PublicForm `json:"public"`

	// private data
	Private *PrivateForm `json:"private"`
}

type PublicForm struct {
	Channel   string `json:"channel" binding:"required"`
	ChainCode string `json:"chaincode" binding:"required"`
	DataJson  string `json:"dataJson" binding:"required"`
}

type PrivateForm struct {
	Channel        string `json:"channel" binding:"required"`
	ChainCode      string `json:"chaincode" binding:"required"`
	CollectionName string `json:"collectionName" binding:"required"`
	DataJson       string `json:"dataJson" binding:"required"`
}

func CreateHandler(ctx *context.ApiContext, c *gin.Context) {
	var form CreateForm
	err := c.BindJSON(&form)
	middleware.StopExec(err)

}
