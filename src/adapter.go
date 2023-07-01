package src

import (
	"github.com/spf13/cobra"
)

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
