package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
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
// @Param			name				formData	string				true	"Account name"
// @Param			operator			path		string				true	"Operator name"
//
// @Param			response-ttl		formData	string				false	"the amount of time the default permission is valid (global) - [#ms(millis) | #s(econds) | m(inutes) | h(ours)] - Default is no time limit."
// @Param			allow-pub-response	formData	string				false	"%s to limit how often a client can publish to reply subjects [with an optional count, --allow-pub-response=n] (global)"
// @Param			allow-pub			formData	string				false	"add publish %s - comma separated list or option can be specified multiple times"
// @Param			allow-pubsub		formData	string				false	"add publish and subscribe %s - comma separated list or option can be specified multiple times"
// @Param			allow-sub			formData	string				false	"add subscribe %s - comma separated list or option can be specified multiple times"
// @Param			deny-pub			formData	string				false	"add deny publish %s - comma separated list or option can be specified multiple times"
// @Param			deny-pubsub			formData	string				false	"add deny publish and subscribe %s - comma separated list or option can be specified multiple times"
// @Param			deny-sub			formData	string				false	"add deny subscribe %s - comma separated list or option can be specified multiple times"
// @Param			start				formData	string				false	"valid from ('0' is always, '3d' is three days) - yyyy-mm-dd, #m(inutes), #h(ours), #d(ays), #w(eeks), #M(onths), #y(ears)"
// @Param			expiry				formData	string				false	"valid until ('0' is always, '2M' is two months) - yyyy-mm-dd, #m(inutes), #h(ours), #d(ays), #w(eeks), #M(onths), #y(ears)"
// @Success		200					{object}	SimpleJSONResponse	"Account added"
// @Failure		400					{object}	SimpleJSONResponse	"Bad request"
// @Failure		500					{object}	string				"Internal error"
func addAccount(c echo.Context) error {
	var addCmd = lookupCommand(nsc.GetRootCmd(), "add")
	var addAccountCmd = lookupCommand(addCmd, "account")

	if err := nsc.GetConfig().SetOperator(c.Param("operator")); err != nil {
		return err
	}

	err := setFlagsIfInForm(addAccountCmd, c.FormValue, []string{
		// "url", is not supported yet.
		"response-ttl",
		"allow-pub-response",
		"allow-pub",
		"allow-pubsub",
		"allow-sub",
		"deny-pub",
		"deny-pubsub",
		"deny-sub",
		"start",
		"expiry",
	})
	if err != nil {
		return badRequest(c, err)
	}
	if err := addAccountCmd.RunE(addAccountCmd, []string{c.FormValue("name")}); err != nil {
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
// @Param			operator	path	string	true	"Operator name"
// @Summary		Describes an account
// @Description	Returns json object with account description
// @Success		200	{object}	AccountDescription	"Operator description"
// @Failure		500	{object}	string				"Internal error"
func describeAccount(c echo.Context) error {
	store, err := nsc.GetStoreForOperator(c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	claim, err := store.ReadRawAccountClaim(c.Param("name"))
	if err != nil {
		return badRequest(c, err)
	}

	body, err := bodyAsJson(claim)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSONBlob(200, body)
}
