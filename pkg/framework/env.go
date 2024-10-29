package framework

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`

	MaxMultipartMemory int64 `mapstructure:"MAX_MULTIPART_MEMORY"`

	TimeZone string `mapstructure:"TIMEZONE"`
}

var globalEnv = Env{
	MaxMultipartMemory: 10 << 20, // 10 MB
}

func GetEnv() Env {
	return globalEnv
}

func NewEnv() Env {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read cofiguration", err)
	}

	viper.SetDefault("TIMEZONE", "UTC")

	err = viper.Unmarshal(&globalEnv)
	if err != nil {
		log.Fatal("environment cant be loaded: ", err)
	}

	return globalEnv
}
