package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"io"
	"os"
)

func init() {
	module := "account"
	initInfo(module)

	GetEchoRoot().POST("operator/:operator/"+module, addAccount)

	GetEchoRoot().GET("operator/:operator/"+module+"s", listAccounts)

	GetEchoRoot().GET("operator/:operator/"+module+"/:name", describeAccount)
}

// @Tags			Account
// @Router			/operator/{operator}/account [post]
// @Summary		Add an account
// @Description	Add an account with given operator to the store
// @Param			name		formData	string				true	"Account name"
// @Param			operator	path		string				true	"Operator name"
// @Success		200			{object}	SimpleJSONResponse	"Account added"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func addAccount(c echo.Context) error {
	nsc.GetConfig().Operator = c.Param("operator")

	cmd := nsc.CreateAddAccountCmd()
	if err := cmd.RunE(cmd, []string{c.FormValue("name")}); err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Account added",
	})
}

// @Tags			Account
// @Router			/operator/{operator}/accounts [get]
// @Summary		List accounts
// @Param			operator	path	string	true	"Operator name"
// @Description	Returns json list of existing accounts for given operator
// @Success		200	{object}	[]string	"Operator's accounts list"
// @Failure		500	{object}	string		"Internal error"
func listAccounts(c echo.Context) error {
	config := nsc.GetConfig()

	// TODO:
	//  {
	//   "status": "400",
	//   "message": "`/nsc/store` is not a valid data directory: stat /nsc/store/.nsc: no such file or directory"
	//  }

	config.Operator = c.QueryParam("operator")

	if accounts, err := config.ListAccounts(); err != nil {
		return badRequest(c, err)
	} else {
		return c.JSON(200, map[string][]string{"accounts": accounts})
	}
}

// @Tags			Account
// @Router			/operator/{operator}/account/{name} [get]
// @Param			name		path	string	true	"Account name"
// @Param			operator	query	string	true	"Operator name"
// @Summary		Describes an account
// @Description	Returns json object with account description
// @Success		200	{object}	AccountDescription	"Operator description"
// @Failure		500	{object}	string				"Internal error"
func describeAccount(c echo.Context) error {
	config := nsc.GetConfig()

	config.Operator = c.QueryParam("operator")

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
