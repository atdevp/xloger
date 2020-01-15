package main

import (
	"flag"
	"os"

	"github.com/atdevp/devlib/file"
	"github.com/xloger/apiserver"
	"github.com/xloger/db"
	"github.com/xloger/g"
	"github.com/xloger/sync"
)

func initDir() {
	storePath := g.Config().XlogNodeConfig.StoragePath
	if file.IsExist(storePath) {
		return
	}

	if err := os.MkdirAll(storePath, 0755); err != nil {
		panic(err)
	}
}

func initConfig() {
	cfg := flag.String("c", "example.cfg.json", "config file")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
}

func main() {

	initConfig()
	initDir()
	db.InitDB()
	g.LogFileCrontabSet.New()
	go sync.SyncLogFileTask()

	apiserver.Start()
}
