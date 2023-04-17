package utils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func resMsg(msg string) string {
	if msg == "successful" {
		return msg
	}
	return fmt.Sprintf("failed to %s user", msg)
}

// JSON serializes the api response properly to json
func JSON(c echo.Context, message string, status int, data any) error {
	switch data.(type) {
	case error:
		return c.JSON(status, map[string]any{
			"message": resMsg(message),
			"errors":  data.(error).Error(),
			"status":  http.StatusText(status),
		})
	default:
		return c.JSON(status, map[string]any{
			"message": resMsg(message),
			"data":    data,
			"status":  http.StatusText(status),
		})
	}
}
