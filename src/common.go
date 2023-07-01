package src

import "github.com/spf13/cobra"
import "fmt"

func notACommand() *cobra.Command {
	return &cobra.Command{}
}

func initInfo(value string) {
	fmt.Println("Initializing handlers [" + value + "]")
}
