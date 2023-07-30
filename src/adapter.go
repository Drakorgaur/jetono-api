package src

import (
	"github.com/labstack/echo/v4"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/spf13/cobra"
)

func init() {
	resetCliMiddleware()
}

func resetCliMiddleware() {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			resetGlobalVars()
			return next(c)
		}
	})
}

// resetGlobalVars resets global variables from nsc package as one
// command can set it up, which will affect other commands.
func resetGlobalVars() {
	nsc.KeyPathFlag = ""
}

func lookupCommand(root *cobra.Command, name string) *cobra.Command {
	for _, e := range root.Commands() {
		if e.Name() == name {
			return e
		}
	}
	return nil
}

func omitCmdOutput() error {
	// TODO: implement
	// function should set cobra.Command outputFlag to /dev/null
	return nil
}
