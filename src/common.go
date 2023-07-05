package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

func initInfo(value string) {
	fmt.Println("Initializing handlers [" + value + "]")
}

func badRequest(c echo.Context, err error) error {
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
