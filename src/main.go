package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nsc/cmd"
	"github.com/nats-io/nsc/cmd/store"
	"os"

	"github.com/swaggo/echo-swagger"

	_ "github.com/Drakorgaur/jetono-api/docs"
)

var e = echo.New()

func GetEchoRoot() *echo.Echo {
	return e
}

func init() {
	cmd.Json = true

	if err := setUp(); err != nil {
		_ = fmt.Errorf("init failed")
	}

	config := cmd.GetConfig()

	storeDir := os.Getenv("NSC_STORE")
	if storeDir == "" {
		storeDir = "/tmp"
	}

	store.KeyStorePath = "/nsc/keys"
	config.StoreRoot = "/" + storeDir

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
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Logger.Fatal(e.Start(":1323"))
}
