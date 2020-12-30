package private

import (
	"github.com/leyle/fabric-app-starter/context"
)

type PrivateForm struct {
	App    string `json:"-"`
	DataId string `json:"-"`

	Channel        string `json:"channel" binding:"required"`
	ChainCode      string `json:"chaincode" binding:"required"`
	CollectionName string `json:"collectionName" binding:"required"`

	DataJson string `json:"dataJson" binding:"required"`
}

func CallPrivateChainCode(ctx *context.ApiContext, form *PrivateForm) error {
	return nil
}
