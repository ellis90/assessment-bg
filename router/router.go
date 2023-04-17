package router

import (
	"github.com/ellis90/assessment-bg/service"
	"github.com/ellis90/assessment-bg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Router(cs *service.CustomerService) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Binder = &utils.CustomBinder{}
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "server running successfully"})
	})
	userRoute := e.Group("/user")
	userRoute.POST("", cs.Create)
	userRoute.GET("", cs.FetchAll)
	userRoute.PUT("", cs.Update)
	userRoute.DELETE("/:id", cs.DeleteById)
	return e
}
