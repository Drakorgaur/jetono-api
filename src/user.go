package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	module := "user"
	initInfo(module)

	GetEchoRoot().POST("operator/:operator/account/:account/"+module, addUser)

	GetEchoRoot().GET("operator/:operator/account/:account/"+module+"s", listUsers)

	GetEchoRoot().GET("operator/:operator/account/:account/"+module+"/:name", describeUser)
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
	if operator := c.Param("operator"); operator == "" {
		return c.JSON(400, &SimpleJSONResponse{
			Status:  "400",
			Message: "operator is required",
		})
	} else {
		nsc.GetConfig().Operator = operator
	}

	if account := c.Param("account"); account == "" {
		return c.JSON(400, &SimpleJSONResponse{
			Status:  "400",
			Message: "account is required",
		})
	} else {
		nsc.GetConfig().Account = account
	}

	cmd := nsc.CreateAddUserCmd()
	err := cmd.Flags().Set("account", nsc.GetConfig().Account)
	if err != nil {
		return err
	}
	err = cmd.RunE(cmd, []string{c.FormValue("name")})
	if err != nil {
		return err
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
		return c.JSON(400, err)
	}

	entries, err := nsc.ListUsers(ctx.StoreCtx().Store, c.Param("account"))
	if err != nil {
		return c.JSON(400, err)
	}

	var users []string

	for _, e := range entries {
		if e.Err == nil {
			_ = append(users, e.Name)
		}
	}

	// TODO: returns null
	return c.JSON(200, users)
}

// @Tags			User
// @Router			/operator/{operator}/account/{account}/user/{name} [get]
// @Param			name		path	string	true	"Account name"
// @Param			account		path	string	true	"Account name"
// @Param			operator	path	string	true	"Operator name"
// @Summary		Describes user
// @Description	Returns json object with user description
// @Success		200	{object}	UserDescription	"Operator description"
// @Failure		500	{object}	string				"Internal error"
func describeUser(c echo.Context) error {
	nsc.GetConfig().Operator = c.Param("operator")

	var describeCmd = lookupCommand(nsc.GetRootCmd(), "describe")
	var operatorCmd = lookupCommand(describeCmd, "user")

	var r, w, _ = os.Pipe()

	old := os.Stdout
	os.Stdout = w

	nsc.Json = true

	err := operatorCmd.Flags().Lookup("account").Value.Set(c.Param("account"))
	if err != nil {
		return err
	}

	err = operatorCmd.RunE(operatorCmd, []string{c.Param("name")})

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
