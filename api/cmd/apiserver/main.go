package main

import (
	"flag"
	"fmt"
	"github.com/leyle/fabric-app-starter/api/chaincodeapi"
	"github.com/leyle/fabric-app-starter/api/context"
	"github.com/leyle/fabric-user-manager/apirouter"
	"github.com/leyle/fabric-user-manager/model"
	"github.com/leyle/go-api-starter/couchdb"
	"github.com/leyle/go-api-starter/ginhelper"
	"github.com/leyle/go-api-starter/logmiddleware"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var port string
	var cfile string
	var conf context.Config
	flag.StringVar(&cfile, "c", "", "-c /config/file/path")
	flag.StringVar(&port, "p", "", "-p 9000")
	flag.Parse()

	if cfile == "" {
		fmt.Println("No config file, please use -c /config/path/file")
		os.Exit(1)
	}

	err := conf.LoadConf(cfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if port != "" {
		conf.Server.Port = port
	}

	jwtCtx := setupJWTContext(&conf)
	ctx := &context.ApiContext{
		Cfg:    &conf,
		JWTCtx: jwtCtx,
	}

	err = preCheck(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go httpServer(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// do something todo
	<-signalChan
	fmt.Println("Shutdown api server")
}

func httpServer(ctx *context.ApiContext) {
	var err error
	logger := logmiddleware.GetLogger(logmiddleware.LogTargetStdout)

	e := ginhelper.SetupGin(logger)
	if ctx.Cfg.Debug {
		ginhelper.PrintHeaders = true
	}

	apiRouter := e.Group("/api")

	// jwt/ca api
	err = apirouter.Init(ctx.JWTCtx)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}
	apirouter.JWTRouter(ctx.JWTCtx, apiRouter.Group(""))

	// chaincode api
	chaincodeapi.ChaincodeRouter(ctx, apiRouter.Group(""))

	addr := ctx.Cfg.Server.GetServerAddr()
	err = e.Run(addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func preCheck(ctx *context.ApiContext) error {
	var err error

	// check fabric connection file exist
	err = context.CheckPathExist(ctx.Cfg.Fabric.CCPath, 4, "fabric connection file")
	if err != nil {
		return err
	}

	// check fabricwallet exist and can read write
	err = context.CheckPathExist(ctx.Cfg.Fabric.WalletPath, 6, "fabric fabricwallet")
	if err != nil {
		// try to create it
		// Logger.Infof("", "wallet path %s doesn't exist, try to create it", ctx.Cfg.Fabric.WalletPath)
		err = os.MkdirAll(ctx.Cfg.Fabric.WalletPath, os.ModePerm)
		if err != nil {
			// Logger.Errorf("", "try to create wallet path[%s] failed, %s", ctx.Cfg.Fabric.WalletPath, err.Error())
			return err
		}
	}

	return nil
}

func setupJWTContext(conf *context.Config) *model.JWTContext {
	dbOpt := &couchdb.CouchDBOption{
		HostPort: conf.Couchdb.HostPort,
		User:     conf.Couchdb.User,
		Passwd:   conf.Couchdb.Passwd,
		Protocol: conf.Couchdb.Protocol,
	}

	registrarOpt := &model.FabricCARegistrar{
		EnrollId: conf.Admin.Username,
		Secret:   conf.Admin.Password,
	}

	gwOpt := &model.FabricGWOption{
		CCPath:     conf.Fabric.CCPath,
		WalletPath: conf.Fabric.WalletPath,
		OrgName:    conf.Fabric.OrgName,
	}

	jwtOpt := &model.JWTOption{
		Secret:      []byte(conf.JWT.Secret),
		ExpireHours: conf.JWT.ExpireHours,
	}

	opt := &model.Option{
		CouchDBOpt:     dbOpt,
		Registrar:      registrarOpt,
		FabricGWOption: gwOpt,
		JWTOpt:         jwtOpt,
	}

	ctx := &model.JWTContext{
		Opt: opt,
	}
	return ctx
}
