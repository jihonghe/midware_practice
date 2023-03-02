package config

import (
	"flag"

	"github.com/spf13/viper"
)

var (
	Config *viper.Viper
)

func Init() {
	var confPath string
	flag.StringVar(&confPath, "config", "", "Config path")
	flag.Parse()

	Config = viper.New()
	Config.SetConfigFile(confPath)
	Config.WatchConfig()
	err := Config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
