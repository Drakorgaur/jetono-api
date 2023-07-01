package src

import (
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
