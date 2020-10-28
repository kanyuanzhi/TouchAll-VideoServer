package utils

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	CViper *viper.Viper
}

func NewConfig() *Config {
	cviper := viper.New()
	cviper.SetConfigName("dataCenterConfig")
	cviper.AddConfigPath("./")
	cviper.SetConfigType("json")
	err := cviper.ReadInConfig()
	if err != nil {
		log.Printf("config file error: %s\n", err)
	}
	return &Config{
		CViper: cviper,
	}
}

func (config *Config) GetValue(key string) interface{} {
	value := config.CViper.Get(key)
	return value
}

func (config *Config) GetSocketConfig() interface{} {
	port := config.GetValue("socket_server.port")
	return port
}

func (config *Config) GetWebSocketConfig() interface{} {
	port := config.GetValue("websocket_server.port")
	return port
}
