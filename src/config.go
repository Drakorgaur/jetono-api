package src

import (
	"fmt"
	nsc "github.com/nats-io/nsc/v2/cmd"
	"os"
)

var nkeys = "NKEYS_PATH"

func checkForEnv(env string, dirName string) error { // TODO: rename
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to configure application. Can not get pwd")
	}

	if os.Getenv(env) == "" {
		err := os.Setenv(env, dir+"/"+nsc.GetToolName()+"/"+dirName)
		if err != nil {
			return err
		}
	}

	return nil
}

func setUp() error {
	// TODO: set this up
	if err := checkForEnv(nsc.NscHomeEnv, "home"); err != nil {
		return err
	}

	if err := checkForEnv(nkeys, "keys"); err != nil {
		return err
	}

	return nil
}
