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

	GetEchoRoot().PUT("operator/:name", updateOperator)
}

type addOperatorForm struct {
	Name               string `json:"name" validate:"required" `
	GenerateSigningKey bool   `json:"generate_signing_key,omitempty"`
	Sys                bool   `json:"sys,omitempty"`
	Force              bool   `json:"force,omitempty"`
	Start              string `json:"start,omitempty"`
	Expiry             string `json:"expiry,omitempty"`
}

// @Tags			Operator
// @Router			/operator [post]
// @Summary		Add an operator
// @Description	Add an operator to the store
// @Param			json	body		addOperatorForm		true	"request body"
// @Success		200		{object}	SimpleJSONResponse	"Operator added"
// @Failure		400		{object}	SimpleJSONResponse	"Bad request"
// @Failure		500		{object}	string				"Internal error"
func addOperator(c echo.Context) error {
	var addCmd = lookupCommand(nsc.GetRootCmd(), "add")
	var addOperatorCmd = lookupCommand(addCmd, "operator")

	err := setFlagsIfInJson(addOperatorCmd, &addOperatorForm{}, c)
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
	operator := c.Param("name")
	s, err := nsc.GetStoreForOperator(operator)
	if err != nil {
		return badRequest(c, err)
	}

	s.Info.Name = operator
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

type updateOperatorForm struct {
	Tag                   string `json:"tag,omitempty"`
	RmTag                 string `json:"rm_tag,omitempty"`
	AccountJwtServerUrl   string `json:"account_jwt_server_url,omitempty"`
	SystemAccount         string `json:"system_account,omitempty"`
	ServiceUrl            string `json:"service_url,omitempty"`
	RmServiceUrl          string `json:"rm_service_url,omitempty"`
	RequireSigningKeys    bool   `json:"require_signing_keys,omitempty"`
	RmAccountJwtServerUrl string `json:"rm_account_jwt_server_url,omitempty"`
}

// @Tags			Operator
// @Router			/operator/{name} [patch]
// @Param			name	path	string				true	"Operator name"
// @Param			json	body	updateOperatorForm	true	"request body"
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

	err := setFlagsIfInJson(updateOperatorCmd, &updateOperatorForm{}, c)
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
