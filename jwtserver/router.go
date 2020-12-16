package jwtserver

import (
	"github.com/gin-gonic/gin"
	"github.com/leyle/fabric-app-starter/authproxy"
	"github.com/leyle/fabric-app-starter/context"
)

func JWTRouter(ctx *context.ApiContext, g *gin.RouterGroup) {
	// no auth
	nR := g.Group("/jwt")
	{
		// login
		nR.POST("/login", func(c *gin.Context) {
			LoginHandler(ctx, c)
		})
	}

	// need auth
	authR := g.Group("/jwt", func(c *gin.Context) {
		authproxy.Auth(ctx, c)
	})
	{
		authR.POST("/create", func(c *gin.Context) {
			CreateUserHandler(ctx, c)
		})
	}
}
