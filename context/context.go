package context

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type ApiContext struct {
	Cfg    *Config
	Wallet *gateway.Wallet
	DbFile string
}

func DumpReqHeaders(c *gin.Context) map[string]string {
	cheaders := c.Request.Header
	headers := make(map[string]string)
	for key, vals := range cheaders {
		headers[key] = vals[0]
	}
	return headers
}
