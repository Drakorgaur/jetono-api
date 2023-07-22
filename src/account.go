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

	GetEchoRoot().PATCH("operator/:operator/"+module+"/:name", updateAccount)
}

// @Tags			Account
// @Router			/operator/{operator}/account [post]
// @Summary		Add an account
// @Description	Add an account with given operator to the store
// @Param			name				formData	string				true	"Account name"
// @Param			operator			path		string				true	"Operator name"
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

	if err, flag := raiseForRequiredFlags(c.FormValue, "name"); err != nil {
		return c.JSON(400, map[string]string{"code": "400", "message": "required form data no filled", "field": flag})
	}

	if err := nsc.GetConfig().SetOperator(c.Param("operator")); err != nil {
		return err
	}

	err := setFlagsIfInForm(addAccountCmd, c.FormValue, []string{
		// "url", is not supported yet.
		"name",
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

// @Tags			Account
// @Router			/operator/{operator}/account/{name} [patch]
// @Param			name					path		string	true	"Account name"
// @Param			operator				path		string	true	"Operator name"
// @Param			tag						formData	string	false	"add tags for user - comma separated list or option can be specified multiple times"
// @Param			rm-tag					formData	string	false	"remove tag - comma separated list or option can be specified multiple times"
// @Param			conns					formData	string	false	"set maximum active connections for the account (-1 is unlimited)"
// @Param			leaf-conns				formData	string	false	"set maximum active leaf node connections for the account (-1 is unlimited)"
// @Param			data					formData	string	false	"set maximum data in bytes for the account (-1 is unlimited)"
// @Param			exports					formData	string	false	"set maximum number of exports for the account (-1 is unlimited)"
// @Param			imports					formData	string	false	"set maximum number of imports for the account (-1 is unlimited)"
// @Param			payload					formData	string	false	"set maximum message payload in bytes for the account (-1 is unlimited)"
// @Param			subscriptions			formData	string	false	"set maximum subscription for the account (-1 is unlimited)"
// @Param			wildcard-exports		formData	bool	false	"exports can contain wildcards"
// @Param			disallow-bearer			formData	bool	false	"require user jwt to not be bearer token"
// @Param			rm-sk					formData	string	false	"remove signing key - comma separated list or option can be specified multiple times"
// @Param			description				formData	string	false	"Description for this account"
// @Param			info-url				formData	string	false	"Link for more info on this account"
// @Param			js-tier					formData	string	false	"JetStream: replication tier (0 creates a configuration that applies to all assets) "
// @Param			rm-js-tier				formData	string	false	"JetStream: remove replication limits for the specified tier (0 is the global tier) this flag is exclusive of all other js flags"
// @Param			js-mem-storage			formData	string	false	"JetStream: set maximum memory storage in bytes for the account (-1 is unlimited / 0 disabled) (units: k/m/g/t kib/mib/gib/tib)"
// @Param			js-disk-storage			formData	string	false	"JetStream: set maximum disk storage in bytes for the account (-1 is unlimited / 0 disabled) (units: k/m/g/t kib/mib/gib/tib)"
// @Param			js-streams				formData	string	false	"JetStream: set maximum streams for the account (-1 is unlimited)"
// @Param			js-consumer				formData	string	false	"JetStream: set maximum consumer for the account (-1 is unlimited)"
// @Param			js-max-mem-stream		formData	string	false	"JetStream: set maximum size of a memory stream for the account (-1 is unlimited / 0 disabled) (units: k/m/g/t kib/mib/gib/tib)"
// @Param			js-max-disk-stream		formData	string	false	"JetStream: set maximum size of a disk stream for the account (-1 is unlimited / 0 disabled) (units: k/m/g/t kib/mib/gib/tib)"
// @Param			js-max-bytes-required	formData	string	false	"JetStream: set whether max stream is required when creating a stream"
// @Param			js-max-ack-pending		formData	string	false	"JetStream: set number of maximum acks that can be pending for a consumer in the account"
// @Param			name					formData	string	false	"account to edit"
// // @Param			js-disable				formData	string	false	"disables all JetStream limits in the account by deleting any limits"
// @Summary		Updates an account
// @Description	Updates an account and returns json with status ok if successful
// @Success		200	{object}	SimpleJSONResponse	"Status ok"
// @Failure		500	{object}	string				"Internal error"
func updateAccount(c echo.Context) error {
	var updateCmd = lookupCommand(nsc.GetRootCmd(), "edit")
	var updateAccountCmd = lookupCommand(updateCmd, "account")

	if err := nsc.GetConfig().SetOperator(c.Param("operator")); err != nil {
		return badRequest(c, err)
	}

	err := setFlagsIfInForm(updateAccountCmd, c.FormValue, []string{
		"start", "expiry", "tag", "rm-tag", "conns", "leaf-conns", "exports", "imports", "subscriptions",
		"payload", "data", "wildcard-exports", "sk", "rm-sk", "description", "info-url", "response-ttl", "allow-pub-response",
		"allow-pub-response", "allow-pub", "allow-pubsub", "allow-sub", "deny-pub", "deny-pubsub", "deny-sub",
		"rm-response-perms", "rm", "max-responses", "disallow-bearer",
		"js-tier",
		"rm-js-tier",
		"js-mem-storage",
		"js-disk-storage",
		"js-streams",
		"js-consumer",
		"js-max-mem-stream",
		"js-max-disk-stream",
		"js-max-bytes-required",
		"js-max-ack-pending",
		// "js-disable",
	})

	if err != nil {
		return badRequest(c, err)
	}

	if err := updateAccountCmd.RunE(updateAccountCmd, []string{c.Param("name")}); err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Account updated",
	})
}
