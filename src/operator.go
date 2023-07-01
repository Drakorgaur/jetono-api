package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	initInfo("operator")

	GetEchoRoot().POST("operator", addOperator)

	GetEchoRoot().GET("operators", listOperators)

	GetEchoRoot().GET("operator/:name", describeOperator)
}

func addOperator(c echo.Context) error {
	var params nsc.AddOperatorParams // TODO: fill params by api body

	if err := nsc.RunStoreLessAction(&cobra.Command{}, []string{c.FormValue("name")}, &params); err != nil {
		return err
	}
	return nil
}

func listOperators(c echo.Context) error {
	operators := nsc.GetConfig().ListOperators()

	return c.JSON(200, operators)
}

type ExtendedDescribeOperatorParams struct {
	*nsc.DescribeOperatorParams
}

func describeOperator(c echo.Context) error {
	var describeCmd = lookupCommand(nsc.GetRootCmd(), "describe")
	var operatorCmd = lookupCommand(describeCmd, "operator")

	var r, w, _ = os.Pipe()

	old := os.Stdout
	os.Stdout = w

	nsc.Json = true
	err := operatorCmd.RunE(operatorCmd, []string{c.Param("name")})
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
