package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

func initInfo(value string) {
	fmt.Println("Initializing handlers [" + value + "]")
}

func badRequest(c echo.Context, err error) error {
	// TODO: rewrite as echo.Context.Error
	return c.JSON(400, &SimpleJSONResponse{
		Status:  "400",
		Message: err.Error(),
	})
}

// copied from source.
func bodyAsJson(data []byte) ([]byte, error) {
	chunks := bytes.Split(data, []byte{'.'})
	if len(chunks) != 3 {
		return nil, errors.New("data is not a jwt")
	}
	body := chunks[1]
	d, err := base64.RawURLEncoding.DecodeString(string(body))
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %v", err)
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(d, &m); err != nil {
		return nil, fmt.Errorf("error parsing json: %v", err)
	}

	j := &bytes.Buffer{}
	encoder := json.NewEncoder(j)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(m); err != nil {
		return nil, fmt.Errorf("error formatting json: %v", err)
	}
	return j.Bytes(), nil
}

func storeType() (storage.Storage, error) {
	storeType := os.Getenv("JETONO_STORE_TYPE")
	switch storeType {
	case "kubernetes":
		return &storage.KubernetesStore{}, nil
	case "objectbox":
		return &storage.ObjectBoxStore{}, nil
	default:
		return nil, fmt.Errorf("invalid store type: %s", storeType)
	}
}

func setFlagsIfInForm(cmd *cobra.Command, getFlag func(string) string, flags []string) error {
	for _, flag := range flags {
		value := getFlag(flag)
		f := cmd.Flag(flag)
		f.Changed = true
		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{value})
		} else {
			_ = f.Value.Set(value)
		}
	}
	return nil
}

func raiseForRequiredFlags(getFlag func(string) string, flags ...string) (error, string) {
	for _, flag := range flags {
		if getFlag(flag) == "" {
			return fmt.Errorf("required form data is not set"), flag
		}
	}
	return nil, ""
}
