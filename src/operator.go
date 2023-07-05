package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/spf13/cobra"
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
		return badRequest(c, err)
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

	return c.JSON(200, map[string][]string{"operators": operators})
}

// @Tags			Operator
// @Router			/operator/{name} [get]
// @Param			name	path	string	true	"Operator name"
// @Summary		Describes an operator
// @Description	Returns json object with operator description
// @Success		200	{object}	OperatorDescription	"Operator description"
// @Failure		500	{object}	string				"Internal error"
func describeOperator(c echo.Context) error {
	s, err := nsc.GetStoreForOperator(c.Param("name"))
	if err != nil {
		return badRequest(c, err)
	}

	claim, err := s.ReadRawOperatorClaim()
	if err != nil {
		return badRequest(c, err)
	}

	body, err := bodyAsJson(claim)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSONBlob(
		200,
		body,
	)
}
