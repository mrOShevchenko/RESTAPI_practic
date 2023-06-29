package repository

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func PingHandler(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "PING OK")
}
