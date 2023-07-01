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

// @Tags			Operator
// @Router			/operator [post]
// @Summary		Add an operator
// @Description	Add an operator to the store
// @Param			name	formData	string				true	"Operator name"
// @Success		200		{object}	SimpleJSONResponse	"Operator added"
// @Failure		400		{object}	SimpleJSONResponse	"Bad request"
// @Failure		500		{object}	string				"Internal error"
func addOperator(c echo.Context) error {
	var params nsc.AddOperatorParams // TODO: fill params by api body

	if err := nsc.RunStoreLessAction(&cobra.Command{}, []string{c.FormValue("name")}, &params); err != nil {
		return c.JSON(400, &SimpleJSONResponse{
			Status:  "400",
			Message: err.Error(),
		})
	}
	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Operator added",
	})
}

// @Tags			Operator
// @Router			/operators [get]
// @Summary		List operators
// @Description	Returns json list of existing operators
// @Success		200	{object}	[]string	"Operators list"
// @Failure		500	{object}	string		"Internal error"
func listOperators(c echo.Context) error {
	operators := nsc.GetConfig().ListOperators()

	return c.JSON(200, operators)
}

// @Tags			Operator
// @Router			/operator/{name} [get]
// @Param			name	path	string	true	"Operator name"
// @Summary		Describes an operator
// @Description	Returns json object with operator description
// @Success		200	{object}	OperatorDescription	"Operator description"
// @Failure		500	{object}	string				"Internal error"
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
