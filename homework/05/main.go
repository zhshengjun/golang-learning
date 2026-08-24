package main

import "os"

func main() {

	// 转账合约，这里没有自己部署新的代币，使用了 LINK，合约地址不一样而已
	TransferToken(os.Getenv("ACCOUNT_FROM_PRIVATE"), os.Getenv("ACCOUNT_TO_ADDRESS"), os.Getenv("TOKEN"), 1)
}
