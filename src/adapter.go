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
