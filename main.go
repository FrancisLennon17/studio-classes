package main

import (
	"flag"
	"os"

	sql "github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/francislennon17/studio-classes/http"
	"github.com/francislennon17/studio-classes/configs"
	"github.com/francislennon17/studio-classes/data/db"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./configs/configs.json", "configuration file path")
	flag.Parse()

	log.SetOutput(os.Stdout)

	conf, err := configs.GetConfig(configFile)
	if err != nil {
		panic(err)
	}

	log.Debug("opening db connection")
	dbCon, err := sql.Open("mysql", conf.GenerateDSN())
	if err != nil {
		panic(err)
	}
	defer dbCon.Close()
	
	dataSource := db.NewDataSource(dbCon)
	http.ListenAndServe(dataSource)
}