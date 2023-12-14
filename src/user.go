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
	root := GetEchoRoot()

	root.POST("operator/:operator/account/:account/user", addUser)

	root.GET("operator/:operator/account/:account/users", listUsers)

	root.GET("operator/:operator/account/:account/user/:name", describeUser)

	root.GET("creds/operator/:operator/account/:account/user/:name", generateUser)

	root.DELETE("operator/:operator/account/:account/user/:name", revokeUser)

	root.PATCH("operator/:operator/account/:account/user/:name", updateUser)
}

type addUserForm struct {
	Name    string `json:"name"`
	Account string `json:"account"`
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user [post]
// @Summary		Add user
// @Description	Add user with given operator and account to the store
// @Param			json		body		addUserForm			true	"json"
// @Param			account		path		string				true	"Account name"
// @Param			operator	path		string				true	"Operator name"
// @Success		200			{object}	SimpleJSONResponse	"User added"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func addUser(c echo.Context) error {
	cmd := nsc.CreateAddUserCmd()
	if err := nsc.GetConfig().SetOperator(c.Param("operator")); err != nil {
		return err
	}
	if err := nsc.GetConfig().SetAccount(c.Param("account")); err != nil {
		return err
	}

	s := &addUserForm{}
	err := setFlagsIfInJson(cmd, s, c)
	if err != nil {
		return badRequest(c, err)
	}

	if err := cmd.Flags().Set("account", nsc.GetConfig().Account); err != nil {
		return badRequest(c, err)
	}

	if err := cmd.RunE(cmd, []string{s.Name}); err != nil {
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
// @Router		/creds/operator/{operator}/account/{account}/user/{user} [get]
// @Summary	Generate user credentials
// @Param		user		path	string	true	"Username"
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

type updateUserForm struct {
	Tag              string `json:"tag,omitempty"`
	RmTag            string `json:"rm_tag,omitempty"`
	Start            string `json:"start,omitempty"`
	Expiry           string `json:"expiry,omitempty"`
	Time             string `json:"time,omitempty"`
	RmTime           string `json:"rm_time,omitempty"`
	Locale           string `json:"locale,omitempty"`
	SourceNetwork    string `json:"source_network,omitempty"`
	RmSourceNetwork  string `json:"rm_source_network,omitempty"`
	ConnType         string `json:"conn_type,omitempty"`
	RmConnType       string `json:"rm_conn_type,omitempty"`
	Subs             string `json:"subs,omitempty"`
	Data             string `json:"data,omitempty"`
	Payload          string `json:"payload,omitempty"`
	Bearer           bool   `json:"bearer,omitempty"`
	ResponseTTL      string `json:"response_ttl,omitempty"`
	AllowPubResponse string `json:"allow_pub_response,omitempty"`
	AllowPub         string `json:"allow_pub,omitempty"`
	AllowPubSub      string `json:"allow_pubsub,omitempty"`
	AllowSub         string `json:"allow_sub,omitempty"`
	DenyPub          string `json:"deny_pub,omitempty"`
	DenyPubSub       string `json:"deny_pubsub,omitempty"`
	DenySub          string `json:"deny_sub,omitempty"`
	RmResponsePerms  string `json:"rm_response_perms,omitempty"`
	Rm               string `json:"rm,omitempty"`
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user/{name} [patch]
// @Param			name		path	string			true	"Username"
// @Param			account		path	string			true	"Account name"
// @Param			operator	path	string			true	"Operator name"
// @Param			json		body	updateUserForm	true	"add tags for user - comma separated list or option can be specified multiple times"
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

	err := setFlagsIfInJson(updateUserCmd, &updateUserForm{}, c)
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
