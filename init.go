package main

import (
	"fmt"
	"os"
	"path"
	"priceupdater/stocks-api/bootconfig"
	"priceupdater/stocks-api/cron"
	"priceupdater/stocks-api/db"
	"priceupdater/stocks-api/utils/logging"
)

const (
	app = "stocks-api"
)

func init() {
	pwd, _ := os.Getwd()
	os.Setenv("APP", app)
	if len(os.Getenv("PORT")) == 0 {
		os.Setenv("PORT", "19093")
	}
	if len(os.Getenv("CONFIGPATH")) == 0 {
		os.Setenv("CONFIGPATH", path.Join(pwd, "config/config.local.json"))
	}
	if len(os.Getenv("LOGDIR")) == 0 {
		os.Setenv("LOGDIR", path.Join(pwd, "logs"))
	}
	fmt.Println("CONFIGPATH=", os.Getenv("CONFIGPATH"))
	fmt.Println("LOGDIR=", os.Getenv("LOGDIR"))
	bootconfig.InitConfig()
	// Loading DB connections
	db.Init()

	logging.NewLogger()
	cron.NewCron()
}
