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

	GetEchoRoot().PATCH("operator/:name", updateOperator)
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

	if err, flag := raiseForRequiredFlags(c.FormValue, "name"); err != nil {
		return c.JSON(400, map[string]string{"code": "400", "message": "required form data no filled", "field": flag})
	}

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

// @Tags			Operator
// @Router			/operator/{name} [patch]
// @Param			name						path		string	true	"Operator name"
// @Param			tag							formData	string	false	"add tags for user - comma separated list or option can be specified multiple times"
// @Param			rm-tag						formData	string	false	"remove tag - comma separated list or option can be specified multiple times"
// @Param			account-jwt-server-url		formData	string	false	"set account jwt server url for nsc sync (only http/https/nats urls supported if updating with nsc)"
// @Param			system-account				formData	string	false	"set system account by account by public key or name"
// @Param			service-url					formData	string	false	"add an operator service url - comma separated list or option can be specified multiple times"
// @Param			rm-service-url				formData	string	false	"remove an operator service url - comma separated list or option can be specified multiple times"
// @Param			require-signing-keys		formData	bool    false	"require accounts/user to be signed with a signing key"
// @Param			rm-account-jwt-server-url	formData	string	false	"clear account server url"
// @Summary		Updates an operator
// @Description	Updates an operator and returns json with status ok if successful
// @Success		200	{object}	SimpleJSONResponse	"Status ok"
// @Failure		500	{object}	string				"Internal error"
func updateOperator(c echo.Context) error {
	var updateCmd = lookupCommand(nsc.GetRootCmd(), "edit")
	var updateOperatorCmd = lookupCommand(updateCmd, "operator")

	if err := nsc.GetConfig().SetOperator(c.Param("name")); err != nil {
		return badRequest(c, err)
	}

	err := setFlagsIfInForm(updateOperatorCmd, c.FormValue, []string{
		"tag",
		"rm-tag",
		"account-jwt-server-url",
		"system-account",
		"service-url",
		"rm-service-url",
		"require-signing-keys",
		"rm-account-jwt-server-url",
		//
		"start",
		"expiry",
	})
	if err != nil {
		return badRequest(c, err)
	}

	if err := updateOperatorCmd.RunE(updateOperatorCmd, []string{c.Param("name")}); err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Operator updated",
	})
}
