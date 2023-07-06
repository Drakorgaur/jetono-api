package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
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
// @Param			name					formData	string				true	"Operator name"
// @Param			generate-signing-key	formData	bool				false	"generate a signing key with the operator"
// @Param			sys						formData	bool				false	"generate system account with the operator (if specified will be signed with signing key)"
// @Param			force					formData	bool				false	"on import, overwrite existing when already present"
// @Param			start					formData	string				false	"valid from ('0' is always, '3d' is three days) - yyyy-mm-dd, #m(inutes), #h(ours), #d(ays), #w(eeks), #M(onths), #y(ears)"
// @Param			expiry					formData	string				false	"valid until ('0' is always, '2M' is two months) - yyyy-mm-dd, #m(inutes), #h(ours), #d(ays), #w(eeks), #M(onths), #y(ears)"
// @Success		200						{object}	SimpleJSONResponse	"Operator added"
// @Failure		400						{object}	SimpleJSONResponse	"Bad request"
// @Failure		500						{object}	string				"Internal error"
func addOperator(c echo.Context) error {
	var addCmd = lookupCommand(nsc.GetRootCmd(), "add")
	var addOperatorCmd = lookupCommand(addCmd, "operator")

	err := setFlagsIfInForm(addOperatorCmd, c.FormValue, []string{
		// "url", is not supported yet.
		"generate-signing-key",
		"sys",
		"force",
		"start",
		"expiry",
	})
	if err != nil {
		return badRequest(c, err)
	}

	if err := addOperatorCmd.RunE(addOperatorCmd, []string{c.FormValue("name")}); err != nil {
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
