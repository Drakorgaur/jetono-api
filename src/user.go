package src

import (
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
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user [post]
// @Summary		Add user
// @Description	Add user with given operator and account to the store
// @Param			name		formData	string				true	"Username"
// @Param			account		path		string				true	"Account name"
// @Param			operator	path		string				true	"Operator name"
// @Success		200			{object}	SimpleJSONResponse	"User added"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func addUser(c echo.Context) error {
	nsc.GetConfig().Operator = c.Param("operator")
	nsc.GetConfig().Account = c.Param("account")

	cmd := nsc.CreateAddUserCmd()
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
	store, err := nsc.GetStoreForOperator(c.Param("operator"))
	if err != nil {
		return badRequest(c, err)
	}

	claim, err := store.ReadRawUserClaim(c.Param("account"), c.Param("name"))
	if err != nil {
		return badRequest(c, err)
	}

	body, err := bodyAsJson(claim)
	if err != nil {
		return badRequest(c, err)
	}

	return c.JSONBlob(200, body)
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
	operator := c.Param("operator")
	account := c.Param("account")
	user := c.Param("name")

	s, err := nsc.GetStoreForOperator(operator)
	if err != nil {
		return badRequest(c, err)
	}

	entityJwt, err := s.Read(store.Accounts, account, store.Users, store.JwtName(user))
	if err != nil {
		return err
	}

	uc, err := jwt.DecodeUserClaims(string(entityJwt))

	keyStore := store.NewKeyStore(operator)
	entityKP, err := keyStore.GetKeyPair(uc.Subject)

	if entityKP == nil {
		return c.JSON(404, map[string]string{"error": "user was not found - please specify it"})
	}

	d, err := nsc.GenerateConfig(s, account, user, entityKP)

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
