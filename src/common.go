package src

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"os"
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

func captureStdout() (*os.File, *os.File, *os.File) {
	var r, w, _ = os.Pipe()

	old := os.Stdout
	os.Stdout = w

	return r, w, old
}

func releaseStdoutLock(r, w, old *os.File) []byte {
	w.Close()
	os.Stdout = old
	all, err := io.ReadAll(r)
	if err != nil {
		return []byte{}
	}
	return all
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
