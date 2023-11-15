package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
	"os"
	"os/exec"
	"reflect"
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

func addStoreConfig() []string {
	return []string{
		"--config-dir", os.Getenv("NSC_HOME"),
		"--data-dir", os.Getenv("NSC_STORE"),
		"--keystore-dir", os.Getenv("NKEYS_PATH"),
	}
}

func runNsc(s interface{}, ctx echo.Context, args ...string) error {
	if err := ctx.Bind(&s); err != nil {
		return err
	}
	m := structs.Map(s)
	fmt.Println(m)

	for flag, value := range m {
		fmt.Printf("flag: %s, value (%s): %s\n", flag, reflect.TypeOf(value), value)
		if value == nil || value == "" || value == false {
			continue
		}
		args = append(args, "--"+strcase.KebabCase(flag))
		if val, ok := value.(string); ok {
			args = append(args, val)
		}
	}

	args = append(args, addStoreConfig()...)

	fmt.Println("nsc", args)

	var stderr bytes.Buffer
	cmd := exec.Command("nsc", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nsc error[%s]: %s", err, stderr.String())
	}

	return nil
}

// copied from source. (nsc)
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

func setFlagsIfInJson(cmd *cobra.Command, s interface{}, ctx echo.Context) error {
	if err := ctx.Bind(&s); err != nil {
		return err
	}
	m := structs.Map(s)

	for flag, value := range m {
		if value == nil || value == "" {
			continue
		}
		formatterFlagName := strcase.KebabCase(flag)
		f := cmd.Flag(formatterFlagName)
		if f == nil {
			fmt.Printf("flag %s not found\n", formatterFlagName)
			continue
		}
		f.Changed = true
		var finalValue string
		if val, ok := value.(string); ok {
			finalValue = val
		}
		if val, ok := value.(int); ok {
			finalValue = fmt.Sprintf("%d", val)
		}
		if val, ok := value.(bool); ok {
			finalValue = fmt.Sprintf("%t", val)
		}

		if val, ok := f.Value.(pflag.SliceValue); ok {
			_ = val.Replace([]string{finalValue})
		} else {
			_ = f.Value.Set(finalValue)
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
