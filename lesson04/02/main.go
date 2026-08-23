package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func init() {

	viper.SetConfigFile("config.yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.BindEnv("account.from.private", "ACCOUNT_FROM_PRIVATE")
	err = viper.BindEnv("account.from.address", "ACCOUNT_FROM_ADDRESS")
	err = viper.BindEnv("account.to.address", "ACCOUNT_TO_ADDRESS")
	err = viper.BindEnv("transaction.hash", "TRANSACTION_HASH")
	if err = viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}
	fmt.Println("配置文件:", viper.ConfigFileUsed())
}

func main() {

	Balance(viper.GetString("account.from.address"))

	SendTransaction(viper.GetString("account.from.private"), viper.GetString("account.to.address"), 0.4, []byte("test transaction"))

	Record(viper.GetString("transaction.hash"))
}
