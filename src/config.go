package src

import (
	"fmt"
	nsc "github.com/nats-io/nsc/cmd"
	"os"
)

func checkForEnv(env string, dirName string) error {
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

	if err := checkForEnv("NKEYS_PATH", "keys"); err != nil {
		return err
	}

	return nil
}
