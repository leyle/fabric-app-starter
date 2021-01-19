package contract

import (
	"errors"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

var ErrNoIdData = errors.New("no id data")
