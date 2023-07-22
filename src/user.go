package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/jwt/v2"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/nats-io/nsc/cmd/store"
	"github.com/spf13/cobra"
)

func init() {
	module := "user"
	initInfo(module)

	GetEchoRoot().POST("operator/:operator/account/:account/"+module, addUser)

	GetEchoRoot().GET("operator/:operator/account/:account/"+module+"s", listUsers)

	GetEchoRoot().GET("operator/:operator/account/:account/"+module+"/:name", describeUser)

	GetEchoRoot().GET("creds/operator/:operator/account/:account/"+module+"/:name", generateUser)

	GetEchoRoot().DELETE("operator/:operator/account/:account/"+module+"/:name", revokeUser)

	GetEchoRoot().PATCH("operator/:operator/account/:account/"+module+"/:name", updateUser)
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user [post]
// @Summary		Add user
// @Description	Add user with given operator and account to the store
// @Param			name		formData	string				sstrue	"Username"
// @Param			account		path		string				sstrue	"Account name"
// @Param			operator	path		string				sstrue	"Operator name"
// @Success		200			{object}	SimpleJSONResponse	"User added"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func addUser(c echo.Context) error {
	cmd := nsc.CreateAddUserCmd()
	nsc.GetConfig().Operator = c.Param("operator")
	nsc.GetConfig().Account = c.Param("account")

	if err, flag := raiseForRequiredFlags(c.FormValue, "name"); err != nil {
		return c.JSON(400, map[string]string{"code": "400", "message": "required form data no filled", "field": flag})
	}

	err := setFlagsIfInForm(cmd, c.FormValue, []string{
		// "url", is not supported yet.
		"name",
		"tag",
		"start",
		"expiry",
	})
	if err != nil {
		return badRequest(c, err)
	}

	if err := cmd.Flags().Set("account", nsc.GetConfig().Account); err != nil {
		return badRequest(c, err)
	}

	if err := cmd.RunE(cmd, []string{c.FormValue("name")}); err != nil {
		return badRequest(c, err)
	}
	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "User added",
	})
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/users [get]
// @Summary		List users
// @Param			account		path	string	true	"Account name"
// @Param			operator	path	string	true	"Operator name"
// @Description	Returns json list of existing users for given operator's account
// @Success		200	{object}	[]string	"List of users for given operator's account"
// @Failure		500	{object}	string		"Internal error"
func listUsers(c echo.Context) error {
	nsc.GetConfig().Operator = c.Param("operator")
	nsc.GetConfig().Account = c.Param("account")

	ctx, err := nsc.NewActx(&cobra.Command{}, []string{})
	if err != nil {
		return badRequest(c, err)
	}

	entries, err := nsc.ListUsers(ctx.StoreCtx().Store, c.Param("account"))
	if err != nil {
		return badRequest(c, err)
	}

	var users []string
	var errors []string

	for _, e := range entries {
		if e.Err == nil {
			users = append(users, e.Name)
		} else {
			errors = append(errors, e.Err.Error())
		}
	}

	if len(errors) > 0 {
		return c.JSON(200, errors)
	}

	return c.JSON(200, map[string][]string{"users": users})
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user/{name} [get]
// @Param			name		path	string	true	"Username"
// @Param			account		path	string	true	"Account name"
// @Param			operator	path	string	true	"Operator name"
// @Summary		Describes user
// @Description	Returns json object with user description
// @Success		200	{object}	UserDescription	"Operator description"
// @Failure		500	{object}	string			"Internal error"
func describeUser(c echo.Context) error {
	s, err := nsc.GetStoreForOperator(c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	claim, err := s.ReadRawUserClaim(c.Param("account"), c.Param("name"))
	if err != nil {
		return badRequest(c, err)
	}

	body, err := bodyAsJson(claim)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSONBlob(200, body)
}

func GetUserCreds(operator string, account string, user string) ([]byte, error) {
	s, err := nsc.GetStoreForOperator(operator)
	if err != nil {
		return nil, err
	}

	entityJwt, err := s.Read(store.Accounts, account, store.Users, store.JwtName(user))
	if err != nil {
		return nil, err
	}

	uc, err := jwt.DecodeUserClaims(string(entityJwt))

	keyStore := store.NewKeyStore(operator)
	entityKP, err := keyStore.GetKeyPair(uc.Subject)

	if entityKP == nil {
		return nil, fmt.Errorf("user was not found - please specify it")
	}

	d, err := nsc.GenerateConfig(s, account, user, entityKP)

	if err != nil {
		return nil, err
	}

	return d, nil
}

// @Tags		User
// @Router		/creds/operator/{operator}/account/{account}/user/{name} [get]
// @Summary	Generate user credentials
// @Param		name		path	string	true	"Username"
// @Param		account		path	string	true	"Account name"
// @Param		operator	path	string	true	"Operator name"
// @Description
// @Success	200	{object}	map[string]string	"Operators list"
// @Success	404	{object}	map[string]string	"User was not found"
// @Failure	500	{object}	string				"Internal error"
func generateUser(c echo.Context) error {
	d, err := GetUserCreds(c.Param("operator"), c.Param("account"), c.Param("user"))
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, map[string]string{"creds": string(d)})
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user/{name} [delete]
// @Param			name		path	string	true	"Username"
// @Param			account		path	string	true	"Account name"
// @Param			operator	path	string	true	"Operator name"
// @Summary		Revokes a user
// @Description	Revokes a user
// @Success		200	{object}	map[string]string	"Operator description"
// @Failure		500	{object}	string				"Internal error"
func revokeUser(c echo.Context) error {
	var revokeCmd = lookupCommand(nsc.GetRootCmd(), "revocations")
	var operatorCmd = lookupCommand(revokeCmd, "add-user")

	if err := nsc.GetConfig().SetOperator(c.Param("operator")); err != nil {
		return badRequest(c, err)
	}
	if err := operatorCmd.Flags().Set("account", c.Param("account")); err != nil {
		return badRequest(c, err)
	}
	if err := operatorCmd.Flags().Set("name", c.Param("name")); err != nil {
		return badRequest(c, err)
	}

	if err := operatorCmd.RunE(operatorCmd, []string{}); err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, map[string]string{"status": "ok"})
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user/{name} [patch]
// @Param			name				path		string	true	"Username"
// @Param			account				path		string	true	"Account name"
// @Param			operator			path		string	true	"Operator name"
// @Param			tag					formData	string	false	"add tags for user - comma separated list or option can be specified multiple times"
// @Param			rm-tag				formData	string	false	"remove tag - comma separated list or option can be specified multiple times"
// @Param			start				formData	string	false	"valid from ('0' is always, '3d' is three days) - yyyy-mm-dd, #m(inutes), #h(ours), #d(ays), #w(eeks), #M(onths), #y(ears)"
// @Param			expiry				formData	string	false	"valid until ('0' is always, '2M' is two months) - yyyy-mm-dd, #m(inutes), #h(ours), #d(ays), #w(eeks), #M(onths), #y(ears)"
// @Param			time				formData	string	false	"`add start-end time range of the form "%s-%s" (option can be specified multiple times)`, timeFormat, timeFormat))"
// @Param			rm-time				formData	string	false	"`remove start-end time by start time "%s" (option can be specified multiple times)`, timeFormat))"
// @Param			locale				formData	string	false	"set the locale with which time values are interpreted")
// @Param			source-network		formData	string	false	"add source network for connection - comma separated list or option can be specified multiple times")
// @Param			rm-source-network	formData	string	false	"remove source network for connection - comma separated list or option can be specified multiple times")
// @Param			conn-type			formData	string	false	"set	allowed	connection	types:	%s	%s	%s	%s	%s	%s	-	comma	separated	list	or	option	can	be	specified	multiple	times, jwt.ConnectionTypeLeafnode, jwt.ConnectionTypeMqtt, jwt.ConnectionTypeStandard, jwt.ConnectionTypeWebsocket, jwt.ConnectionTypeLeafnodeWS, jwt.ConnectionTypeMqttWS))"
// @Param			rm-conn-type		formData	string	false	"remove connection types - comma separated list or option can be specified multiple times")
// @Param			subs				formData	string	false	"set maximum number of subscriptions (-1 is unlimited)")
// @Param			data				formData	string	false	"set maximum data in bytes for the user (-1 is unlimited)")
// @Param			payload				formData	string	false	"set maximum message payload in bytes for the account (-1 is unlimited)")
// @Param			bearer				formData	bool 	false	"no connect challenge required for user")
// @Param			response-ttl		formData	string	false	"the amount of stime sthe s%s sis svalid s(global) - s[#ms(millis) | #s(econds) | m(inutes) | h(ours)] - sDefault sis no time limit., typeName))"
// @Param			allow-pub-response	formData	string	false	"%s to slimit how soften sa sclient scan spublish sto sreply ssubjects	[with an optional count, --allow-pub-response=n] (global), typeName))"
// @Param			allow-pub			formData	string	false	"add publish s%s s- scomma sseparated slist sor soption scan sbe sspecified smultiple stimes, typeName))"
// @Param			allow-pubsub		formData	string	false	"add publish sand ssubscribe s%s s- scomma sseparated slist sor soption scan sbe sspecified	multiple stimes, typeName))"
// @Param			allow-sub			formData	string	false	"add subscribe s%s s- scomma sseparated slist sor soption scan sbe sspecified multiple stimes, typeName))"
// @Param			deny-pub			formData	string	false	"add deny spublish s%s s- scomma sseparated slist sor soption scan sbe sspecified smultiple	times, typeName))"
// @Param			deny-pubsub			formData	string	false	"add deny spublish sand ssubscribe s%s s- scomma sseparated slist sor soption scan sbe sspecified multiple times, typeName))"
// @Param			deny-sub			formData	string	false	"add deny ssubscribe s%s s- scomma sseparated slist sor soption scan sbe sspecified smultiple	times, typeName))"
// @Param			rm-response-perms	formData	string	false	"remove sresponse ssettings sfrom s%s, typeName))"
// @Param			rm					formData	string	false	"remove spublish/subscribe sand sdeny s%s - comma separated	list or	option can be specified multiple times, typeName))"
// @Summary		Updates an user
// @Description	Updates an user and returns json with status ok if successful
// @Success		200	{object}	SimpleJSONResponse	"Status ok"
// @Failure		500	{object}	string				"Internal error"
func updateUser(c echo.Context) error {
	var updateCmd = lookupCommand(nsc.GetRootCmd(), "edit")
	var updateUserCmd = lookupCommand(updateCmd, "user")

	if err := nsc.GetConfig().SetOperator(c.Param("operator")); err != nil {
		return badRequest(c, err)
	}
	if err := nsc.GetConfig().SetAccount(c.Param("account")); err != nil {
		return badRequest(c, err)
	}
	if err := updateUserCmd.Flags().Set("account", c.Param("account")); err != nil {
		return badRequest(c, err)
	}

	err := setFlagsIfInForm(updateUserCmd, c.FormValue, []string{
		"start", "expiry", "rm", "allow-pub", "allow-sub", "allow-pubsub",
		"deny-pub", "deny-sub", "deny-pubsub", "tag", "rm-tag", "source-network", "rm-source-network", "payload",
		"rm-response-perms", "max-responses", "response-ttl", "allow-pub-response", "bearer", "rm-time", "time", "conn-type",
		"rm-conn-type", "subs", "data",
	})
	if err != nil {
		return badRequest(c, err)
	}

	if err := updateUserCmd.RunE(updateUserCmd, []string{c.Param("name")}); err != nil {
		return badRequest(c, err)
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "User updated",
	})
}
