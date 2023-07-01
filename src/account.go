package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"io"
	"os"
)

func init() {
	module := "account"
	initInfo(module)

	GetEchoRoot().POST(module, addAccount)

	GetEchoRoot().GET(module+"s", listAccounts)

	GetEchoRoot().GET(module+"/:name", describeAccount)
}

func addAccount(c echo.Context) error {
	fmt.Println("addAccount to config " + nsc.GetConfig().StoreRoot)

	if operator := c.FormValue("operator"); operator == "" {
		operator = "test_operator" // TODO: get from getperator()
		nsc.GetConfig().Operator = operator
		fmt.Println("operator is set " + nsc.GetConfig().Operator)

	}

	cmd := nsc.CreateAddAccountCmd()
	err := cmd.RunE(cmd, []string{c.FormValue("name")})
	if err != nil {
		return err
	}
	return nil
}

func listAccounts(c echo.Context) error {
	accounts, err := nsc.GetConfig().ListAccounts()
	if err != nil {
		return c.JSON(400, err)
	}
	return c.JSON(200, accounts)
}

func describeAccount(c echo.Context) error {
	describeCmd := lookupCommand(nsc.GetRootCmd(), "describe")
	accountCmd := lookupCommand(describeCmd, "account")

	var r, w, _ = os.Pipe()

	old := os.Stdout
	os.Stdout = w

	nsc.Json = true
	err := accountCmd.RunE(accountCmd, []string{c.Param("name")})
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	var b, _ = io.ReadAll(r)
	os.Stdout = old

	return c.JSONBlob(200, b)
}
