package main

import (
	"flag"
	"fmt"
	"github.com/leyle/fabric-app-starter/chaincodeapi"
	"github.com/leyle/fabric-app-starter/context"
	"github.com/leyle/fabric-app-starter/fabricwallet"
	"github.com/leyle/fabric-app-starter/jwtserver"
	. "github.com/leyle/ginbase/consolelog"
	"github.com/leyle/ginbase/middleware"
	"github.com/leyle/ginbase/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var port string
	var cfile string
	var conf context.Config
	flag.StringVar(&cfile, "c", "", "-c /config/file/path")
	flag.StringVar(&port, "p", "", "-p 8200")
	flag.Parse()

	if cfile == "" {
		fmt.Println("No config file, please use -c /config/path/file")
		os.Exit(1)
	}

	err := conf.LoadConf(cfile)
	if err != nil {
		os.Exit(1)
	}

	if port != "" {
		conf.Server.Port = port
	}

	// sqlite init and check
	ctx := &context.ApiContext{
		Cfg:    &conf,
		DbFile: conf.Sqlite.DbPath,
	}

	err = preCheck(ctx)
	if err != nil {
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
	e := middleware.SetupGin()
	if ctx.Cfg.Debug {
		middleware.PrintHeader = true
	}

	apiRouter := e.Group("/api")

	// auth server
	jwtserver.JWTRouter(ctx, apiRouter.Group(""))

	// chaincode api
	chaincodeapi.PublicAndPrivateRouter(ctx, apiRouter.Group(""))

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
		Logger.Infof("", "wallet path %s doesn't exist, try to create it", ctx.Cfg.Fabric.WalletPath)
		err = os.MkdirAll(ctx.Cfg.Fabric.WalletPath, os.ModePerm)
		if err != nil {
			Logger.Errorf("", "try to create wallet path[%s] failed, %s", ctx.Cfg.Fabric.WalletPath, err.Error())
			return err
		}
	}

	// init system admin user and enroll ca credentials
	err = initAdmin(ctx)
	if err != nil {
		return err
	}

	return nil
}

func initAdmin(ctx *context.ApiContext) error {
	ds := context.NewDs(ctx.DbFile)
	defer ds.Close()

	// create table
	ds.Db.AutoMigrate(&jwtserver.UserAccount{})
	ds.Db.AutoMigrate(&jwtserver.CaAccount{})

	salt := util.GenerateDataId()
	account := &jwtserver.UserAccount{
		Id:        ctx.Cfg.Admin.Username,
		Username:  ctx.Cfg.Admin.Username,
		Role:      jwtserver.UserRoleAdmin,
		Salt:      salt,
		Valid:     jwtserver.UserStatusValid,
		PassHash:  util.GenerateHashPasswd(salt, ctx.Cfg.Admin.Password),
		CreatedAt: util.CurUnixTime(),
	}
	account.UpdatedAt = account.CreatedAt

	// query or insert
	var dbAccount jwtserver.UserAccount
	tx := ds.Db.Where(" username = ?", account.Username).First(&dbAccount)
	if tx.RowsAffected == 0 {
		Logger.Infof("", "System don't have admin user, create it")
		tx = ds.Db.Create(account)
		if tx.Error != nil {
			Logger.Errorf("", "Init system admin user failed", tx.Error.Error())
			return tx.Error
		}
	}

	// enroll admin
	return enrollAdmin(ctx)
}

func enrollAdmin(ctx *context.ApiContext) error {
	// check if wallets/{orgadmin}.id exists
	userId := ctx.Cfg.Admin.Username
	passwd := ctx.Cfg.Admin.Password

	err := fabricwallet.EnrollUser(ctx, userId, passwd)

	return err
}
