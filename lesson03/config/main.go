package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile("config.yaml")
	// viper.AddConfigPath(".")
	viper.SetEnvPrefix("LOCAL")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
}

func loadConfig() *Config {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
	return &config
}

func main() {
	config := loadConfig()
	fmt.Println(config.Server.Name)

	host := viper.Get("database.host")
	fmt.Println(host)

}
