package chaincodeapi

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/authproxy"
	"github.com/leyle/fabric-app-starter/context"
)

func PublicAndPrivateRouter(ctx *context.ApiContext, g *gin.RouterGroup) {
	authR := g.Group("/cc/pp", func(c *gin.Context) {
		authproxy.Auth(ctx, c)
	})
	{
		// create new chain data
		authR.POST("/create", func(c *gin.Context) {
			CreateHandler(ctx, c)
		})
	}
}
