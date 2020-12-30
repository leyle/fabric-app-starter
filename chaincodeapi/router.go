package chaincodeapi

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/context"
)

func PublicAndPrivateRouter(ctx *context.ApiContext, g *gin.RouterGroup) {
	authR := g.Group("/chaincode/publicandprivate", func(c *gin.Context) {
		context.Auth(ctx, c)
	})
	{
		// create new chain data
		authR.POST("/create", context.HandlerWrapper(CreateHandler, ctx))

		// get data by id
		authR.GET("/public/info", context.HandlerWrapper(GetByIdHandler, ctx))
	}
}
