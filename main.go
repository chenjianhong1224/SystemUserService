// SaleUserService project main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var statObj *stator

//VERSION : module's version
var VERSION string

func showVersion() {
	fmt.Println(VERSION)
}

var signalNames = map[syscall.Signal]string{
	syscall.SIGINT:  "SIGINT",
	syscall.SIGQUIT: "SIGQUIT",
	syscall.SIGTERM: "SIGTERM",
	syscall.SIGKILL: "SIGKILL",
}

func signalName(s syscall.Signal) string {
	if name, ok := signalNames[s]; ok {
		return name
	}
	return fmt.Sprintf("SIG %d", s)
}

func getConfig() (*Config, error) {
	opts := newOptions()
	opts.InstallFlags()
	if opts.version {
		showVersion()
		os.Exit(0)
	}

	cfg := newConfig()
	if err := cfg.load(opts.configFile); err != nil {
		return nil, fmt.Errorf("failed to load config file: %s", err)
	}

	return cfg, nil
}

func main() {
	defer zap.L().Sync()
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("cfg:%+v\n", cfg)

	db := newDbOperator(cfg)

	// init local stat gloag obj
	statObj = newStator(&cfg.Stat)
	statObj.start()

	BuildLogger(&cfg.Logger)

	svc := &httpHandler{cfg: cfg, systemUserSv: &system_user_service{d: db}, systemRoleSv: &system_role_service{d: db}, systemMenuSv: &system_menu_service{d: db}, systemPrivilegeSv: &system_privilege_service{d: db}}
	if err = svc.start(); err != nil {
		panic(err)
	}

	sigc := make(chan os.Signal, 4)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)
	shutdown := make(chan bool)
	go func() {
		for s := range sigc {
			name := s.String()
			if sig, ok := s.(syscall.Signal); ok {
				name = signalName(sig)
			}
			zap.L().Info(fmt.Sprintf("Received %v, initiating shutdown...", name))
			select {
			case shutdown <- true:
			}
		}
	}()
	<-shutdown

	zap.L().Info("system manager service stopped")
}
