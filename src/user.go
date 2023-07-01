package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	module := "user"
	initInfo(module)

	GetEchoRoot().POST(module, addUser)

	GetEchoRoot().GET(module+"s", listUsers)

	GetEchoRoot().GET("account/"+":account/"+module+"/:name", describeUser)
}

func addUser(c echo.Context) error {
	if operator := c.FormValue("operator"); operator == "" {
		operator = "test_operator" // TODO: get from getperator()
		nsc.GetConfig().Operator = operator
		fmt.Println("operator is set " + nsc.GetConfig().Operator)
	}
	account := c.FormValue("account")
	if account == "" {
		return c.JSON(400, "account is required")
	}

	cmd := nsc.CreateAddUserCmd()
	err := cmd.Flags().Set("account", account)
	if err != nil {
		return err
	}
	err = cmd.RunE(cmd, []string{c.FormValue("name")})
	if err != nil {
		return err
	}
	return nil
}

func listUsers(c echo.Context) error {
	// TODO: fix this

	ctx, err := nsc.NewActx(&cobra.Command{}, []string{})
	if err != nil {
		return c.JSON(400, err)
	}
	entries, err := nsc.ListUsers(ctx.StoreCtx().Store, c.QueryParam("account"))
	if err != nil {
		return c.JSON(400, err)
	}

	var users []string

	for _, e := range entries {
		if e.Err == nil {
			_ = append(users, e.Name)
		}
	}

	return c.JSON(200, users)
}

func describeUser(c echo.Context) error {
	var describeCmd = lookupCommand(nsc.GetRootCmd(), "describe")
	var operatorCmd = lookupCommand(describeCmd, "user")

	var r, w, _ = os.Pipe()

	old := os.Stdout
	os.Stdout = w

	nsc.Json = true

	err := operatorCmd.Flags().Lookup("account").Value.Set(c.Param("account"))
	if err != nil {
		return err
	}

	err = operatorCmd.RunE(operatorCmd, []string{c.Param("name")})

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
