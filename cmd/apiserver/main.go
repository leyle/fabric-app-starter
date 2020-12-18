package main

import (
	"flag"
	"fmt"
	"github.com/leyle/fabric-app-starter/chaincodeapi"
	"github.com/leyle/fabric-app-starter/context"
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
	err = checkPathExist(ctx.Cfg.Fabric.CCPath, 4, "fabric connection file")
	if err != nil {
		return err
	}

	// check fabricwallet exist and can read write
	err = checkPathExist(ctx.Cfg.Fabric.WalletPath, 6, "fabric fabricwallet")
	if err != nil {
		return err
	}

	// init system admin user and enroll ca credentials
	err = initAdmin(ctx)
	if err != nil {
		return err
	}

	return nil
}

// minPermission:
// 4 -> only check if can read
// 4 + 2 = 6 -> check if can read and write
func checkPathExist(path string, permission int, desc string) error {
	// first check if exist
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			Logger.Errorf("", "%s[%s] doesn't exist", desc, path)
		} else {
			Logger.Errorf("", "%s[%s] failed, %s", desc, path, err.Error())
		}
		return err
	}

	// then check if can read or read/write
	var bit uint32 = syscall.O_RDWR
	if permission < 6 {
		bit = syscall.O_RDONLY
	}

	err := syscall.Access(path, bit)
	if err != nil {
		Logger.Errorf("", "%s[%s] cannot access, %s", desc, path, err.Error())
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

	return nil
}
