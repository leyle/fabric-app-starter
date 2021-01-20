package chaincodeapi

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi/dbstate"
	"github.com/leyle/fabric-app-starter/api/context"
)

func ChaincodeRouter(ctx *context.ApiContext, g *gin.RouterGroup) {
	authR := g.Group("", func(c *gin.Context) {
		context.Auth(ctx, c)
	})

	// world state api
	stateR := authR.Group("/state")
	{
		stateR.POST("/index", context.HandlerWrapper(dbstate.CreateCouchdbIndexHandler, ctx))
	}

	ccR := authR.Group("/cc")
	{
		ccR.POST("/new", context.HandlerWrapper(CreateStateHandler, ctx))
	}
}
