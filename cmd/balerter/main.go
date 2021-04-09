package main

import (
	"context"
	"fmt"
	apiManager "github.com/balerter/balerter/internal/api/manager"
	"github.com/balerter/balerter/internal/corestorage"
	alertModule "github.com/balerter/balerter/internal/modules/alert"
	"github.com/balerter/balerter/internal/service"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	channelsManager "github.com/balerter/balerter/internal/chmanager"
	"github.com/balerter/balerter/internal/config"
	coreStorageManager "github.com/balerter/balerter/internal/corestorage/manager"
	dsManager "github.com/balerter/balerter/internal/datasource/manager"
	"github.com/balerter/balerter/internal/logger"
	"github.com/balerter/balerter/internal/metrics"
	"github.com/balerter/balerter/internal/modules"
	chartModule "github.com/balerter/balerter/internal/modules/chart"
	httpModule "github.com/balerter/balerter/internal/modules/http"
	"github.com/balerter/balerter/internal/modules/kv"
	logModule "github.com/balerter/balerter/internal/modules/log"
	runtimeModule "github.com/balerter/balerter/internal/modules/runtime"
	"github.com/balerter/balerter/internal/runner"
	scriptsManager "github.com/balerter/balerter/internal/script/manager"
	uploadStorageManager "github.com/balerter/balerter/internal/upload_storage/manager"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

var (
	version = "undefined"

	loggerOptions []zap.Option // for testing purposes
)

const (
	defaultLuaModulesPath = "./?.lua;./modules/?.lua;./modules/?/init.lua;/modules/?.lua;/modules/?/init.lua"
)

func main() {
	cfg, flg, err := config.New()
	if err != nil {
		log.Printf("error configuration load, %v", err)
		os.Exit(1)
	}

	msg, code := run(cfg, flg)

	log.Print(msg)
	os.Exit(code)
}

func run(
	cfg *config.Config,
	flg *config.Flags,
) (string, int) {
	lua.LuaPathDefault = defaultLuaModulesPath

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	wg := &sync.WaitGroup{}

	if err := validateLogLevel(flg.LogLevel); err != nil {
		return err.Error(), 1
	}

	lgr, err := logger.New(flg.LogLevel, flg.Debug, loggerOptions...)
	if err != nil {
		return fmt.Sprintf("error init zap logger, %v", err), 1
	}

	metrics.SetVersion(version)

	lgr.Logger().Info("balerter start", zap.String("version", version))

	lgr.Logger().Debug("loaded configuration", zap.Any("config", cfg), zap.Any("flags", flg))

	if cfg.LuaModulesPath != "" {
		lua.LuaPathDefault = cfg.LuaModulesPath
	}

	lgr.Logger().Debug("lua modules path", zap.String("path", lua.LuaPathDefault))

	// Scripts sources
	lgr.Logger().Info("init scripts manager")
	scriptsMgr := scriptsManager.New()

	// log if use a 465 port and an empty secure string for an email channel
	if cfg.Channels != nil {
		for idx := range cfg.Channels.Email {
			if cfg.Channels.Email[idx].Port == "465" && cfg.Channels.Email[idx].Secure == "" {
				lgr.Logger().Info("secure port 465 with ssl for email channel " + cfg.Channels.Email[idx].Name)
			}
		}
	}

	if err = scriptsMgr.Init(cfg.Scripts); err != nil {
		return fmt.Sprintf("error init scripts manager, %v", err), 1
	}

	// datasources
	lgr.Logger().Info("init datasources manager")
	dsMgr := dsManager.New(lgr.Logger())
	if err = dsMgr.Init(cfg.DataSources); err != nil {
		return fmt.Sprintf("error init datasources manager, %v", err), 1
	}

	// upload storages
	lgr.Logger().Info("init upload storages manager")
	uploadStoragesMgr := uploadStorageManager.New(lgr.Logger())
	if err = uploadStoragesMgr.Init(cfg.StoragesUpload); err != nil {
		return fmt.Sprintf("error init upload storages manager, %v", err), 1
	}

	// core storages
	lgr.Logger().Info("init core storages manager")
	coreStoragesMgr, err := coreStorageManager.New(cfg.StoragesCore, lgr.Logger())
	if err != nil {
		return fmt.Sprintf("error create core storages manager, %v", err), 1
	}
	coreStorageAlert, err := coreStoragesMgr.Get(cfg.StorageAlert)
	if err != nil {
		return fmt.Sprintf("error get core storage: alert '%s', %v", cfg.StorageAlert, err), 1
	}
	coreStorageKV, err := coreStoragesMgr.Get(cfg.StorageKV)
	if err != nil {
		return fmt.Sprintf("error get core storage: kv '%s', %v", cfg.StorageKV, err), 1
	}

	// ChannelsManager
	lgr.Logger().Info("init channels manager")
	channelsMgr := channelsManager.New(lgr.Logger())
	if err = channelsMgr.Init(cfg.Channels); err != nil {
		return fmt.Sprintf("error init channels manager, %v", err), 1
	}
	// TODO: pass channels manager...

	coreModules := initCoreModules(coreStorageAlert, coreStorageKV, channelsMgr, lgr.Logger(), flg)

	// ---------------------
	// |
	// | Runner
	// |
	lgr.Logger().Info("init runner")
	rnr := runner.New(time.Millisecond*time.Duration(cfg.Scripts.UpdateInterval), scriptsMgr, dsMgr, uploadStoragesMgr, coreModules, flg.Script, lgr.Logger())

	lgr.Logger().Info("run runner")
	go rnr.Watch(ctx, ctxCancel, flg.Once)

	// ---------------------
	// |
	// | API
	// |
	if cfg.API != nil {
		var ln net.Listener
		ln, err = net.Listen("tcp", cfg.API.Address)
		if err != nil {
			return fmt.Sprintf("error create api listener, %v", err), 1
		}
		apis := apiManager.New(cfg.API.Address, coreStorageAlert, coreStorageKV, channelsMgr, rnr, lgr.Logger())
		wg.Add(1)
		go apis.Run(ctx, ctxCancel, wg, ln)

		if cfg.API.ServiceAddress != "" {
			var ln net.Listener
			ln, err = net.Listen("tcp", cfg.API.ServiceAddress)
			if err != nil {
				return fmt.Sprintf("error create service listener, %v", err), 1
			}
			srv := service.New(lgr.Logger())
			wg.Add(1)
			go srv.Run(ctx, ctxCancel, wg, ln)
		}
	}

	// ---------------------
	// |
	// | Shutdown
	// |
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)

	var sig os.Signal

	select {
	case sig = <-ch:
		lgr.Logger().Info("got os signal", zap.String("signal", sig.String()))
		ctxCancel()
	case <-ctx.Done():
	}

	rnr.Stop()

	wg.Wait()

	dsMgr.Stop()
	coreStoragesMgr.Stop()

	lgr.Logger().Info("terminate")

	return "", 0
}

func initCoreModules(
	coreStorageAlert corestorage.CoreStorage,
	coreStorageKV corestorage.CoreStorage,
	chManager *channelsManager.ChannelsManager,
	logger *zap.Logger,
	flg *config.Flags,
) []modules.Module {
	coreModules := make([]modules.Module, 0)

	alertMod := alertModule.New(coreStorageAlert.Alert(), chManager, logger)
	coreModules = append(coreModules, alertMod)

	kvModule := kv.New(coreStorageKV.KV())
	coreModules = append(coreModules, kvModule)

	logMod := logModule.New(logger)
	coreModules = append(coreModules, logMod)

	chartMod := chartModule.New(logger)
	coreModules = append(coreModules, chartMod)

	httpMod := httpModule.New(logger)
	coreModules = append(coreModules, httpMod)

	runtimeMod := runtimeModule.New(flg, logger)
	coreModules = append(coreModules, runtimeMod)

	return coreModules
}

func validateLogLevel(level string) error {
	for _, l := range []string{"ERROR", "INFO", "DEBUG"} {
		if l == level {
			return nil
		}
	}

	return fmt.Errorf("wrong log level")
}
