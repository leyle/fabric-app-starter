package chaincodeapi

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/chaincodeapi/dbstate"
	"github.com/leyle/fabric-app-starter/context"
)

func ChaincodeMiddlewareRouter(ctx *context.ApiContext, g *gin.RouterGroup) {
	authR := g.Group("", func(c *gin.Context) {
		context.Auth(ctx, c)
	})

	// world state api
	stateR := authR.Group("/state")
	{
		stateR.POST("/index", context.HandlerWrapper(dbstate.CreateCouchdbIndexHandler, ctx))
		stateR.GET("/search", context.HandlerWrapper(SearchHandler, ctx))
	}

	// chaincode api
	ccR := authR.Group("/chaincode")
	{
		// create public/private chaincode
		ccR.POST("/publicandprivate/create", context.HandlerWrapper(CreateHandler, ctx))

		// get by id
		ccR.GET("/id/info", context.HandlerWrapper(GetByIdHandler, ctx))
	}
}

/*
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

// db index create api
func DbIndexCreateRouter(ctx *context.ApiContext, g *gin.RouterGroup) {
	authR := g.Group("/state", func(c *gin.Context) {
		context.Auth(ctx, c)
	})
	{
		// create index on world state database(couchdb)
		authR.POST("/index", context.HandlerWrapper(dbstate.CreateCouchdbIndexHandler, ctx))
	}
}
*/
