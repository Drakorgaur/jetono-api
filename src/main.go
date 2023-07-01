package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nsc/cmd"
	"github.com/nats-io/nsc/cmd/store"
	"os"
)

var e = echo.New()

func GetEchoRoot() *echo.Echo {
	return e
}

func init() {
	if err := setUp(); err != nil {
		_ = fmt.Errorf("init failed")
	}

	config := cmd.GetConfig()

	store_dir := os.Getenv("NSC_STORE")
	if store_dir == "" {
		store_dir = "/tmp"
	}

	store.KeyStorePath = "/nsc/keys"
	config.StoreRoot = "/" + store_dir

	if _, err := os.Stat(config.StoreRoot); os.IsNotExist(err) {
		fmt.Println("StoreRoot does not exist")
		err := os.MkdirAll(config.StoreRoot, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println(config.StoreRoot + " is the store_dir root")
}

func Api() {
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
