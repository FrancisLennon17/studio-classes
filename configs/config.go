package configs

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Hostname   string
	Port       int
	DBUser     string
	DBPassword string
	DBName     string
}

func GetConfig(configFile string) (Configuration, error) {
	config := Configuration{}
	err := gonfig.GetConf(configFile, &config)

	return config, err
}

//Generates the data source name for the DB connection
func (c Configuration) GenerateDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.DBUser, c.DBPassword, c.Hostname, c.Port, c.DBName)
}