package src

import (
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
