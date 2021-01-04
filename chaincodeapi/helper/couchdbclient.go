package helper

import (
	"fmt"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/go-api-starter/couchdb"
)

func GetCouchdbClient(ctx *context.ApiContext, dbName string) *couchdb.CouchDBClient {
	cfg := ctx.Cfg.DbState
	opt := &couchdb.CouchDBOption{
		HostPort: cfg.HostPort,
		User:     cfg.User,
		Passwd:   cfg.Passwd,
		Protocol: cfg.Protocol,
	}

	ds := couchdb.New(opt, dbName)
	return ds
}

func GenerateCouchdbDatabaseName(channel, chaincode, collection string) string {
	name := fmt.Sprintf("%s_%s", channel, chaincode)
	if collection != "" {
		// private collection format
		name = fmt.Sprintf("%s$$p%s", name, collection)
	}
	return name
}
