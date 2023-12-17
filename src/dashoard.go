package src

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func init() {
	root := GetEchoRoot()
	root.PUT("dashboard", updateDashboard)

	root.GET("dashboard", getDashboard)
}

const dashboardFile = "dashboards.json"

type DbDataflow struct {
	Created string `json:"created,omitempty"`
	Name    string `json:"name" validate:"required"`
	Server  string `json:"server" validate:"required"`
}

type DbOperator struct {
	Name string `json:"name" validate:"required"`
}

type DbAccount struct {
	Operator string `json:"operator" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type DbUser struct {
	Operator string `json:"operator" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Account  string `json:"account" validate:"required"`
}

type Dashboard struct {
	Dataflows []DbDataflow `json:"dataflows"`
	Operators []DbOperator `json:"operators"`
	Accounts  []DbAccount  `json:"accounts"`
	Users     []DbUser     `json:"users"`
}

// @Tags			Dashboard
// @Router			/dashboard [put]
// @Summary		Create a new dashboard or put a new one
// @Description	Create a new dashboard or put a new one
// @Param			json		body		Dashboard		true	"Dashboard data in json format"
// @Success		200			{object}	SimpleJSONResponse	"Dashboard updated"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func updateDashboard(c echo.Context) error {
	obj := &Dashboard{}
	err := c.Bind(obj)
	if err != nil {
		return badRequest(c, err)
	}

	err = storeJson(dashboardFile, obj)
	if err != nil {
		return badRequest(c, err)
	}

	// check for accounts are set operator
	for _, account := range obj.Accounts {
		if account.Operator == "" {
			return badRequest(c, fmt.Errorf("account %s has no operator", account.Name))
		}
	}

	// check for users are set operator and account
	for _, user := range obj.Users {
		if user.Operator == "" {
			return badRequest(c, fmt.Errorf("user %s has no operator", user.Name))
		}
		if user.Account == "" {
			return badRequest(c, fmt.Errorf("user %s has no account", user.Name))
		}
	}

	return c.JSON(200, &SimpleJSONResponse{
		Status:  "200",
		Message: "Dashboard updated",
	})
}

// @Tags			Dashboard
// @Router			/dashboard [get]
// @Summary		Get the dashboard
// @Description	Get the dashboard
// @Success		200			{object}	Dashboard	"Get the dashboard"
// @Failure		400			{object}	SimpleJSONResponse	"Bad request"
// @Failure		500			{object}	string				"Internal error"
func getDashboard(c echo.Context) error {
	obj := &Dashboard{}
	err := readJsonFile(dashboardFile, obj)
	if err != nil {
		return badRequest(c, err)
	}
	return c.JSON(200, &obj)
}
