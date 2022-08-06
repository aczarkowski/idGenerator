package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"uidGenerator/handler"
	generatorMiddleware "uidGenerator/middleware"
	"uidGenerator/timeprovider/julian"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(threadid.ThreadId(epoch.New(int64(1420070400000))))
	e.Use(generatorMiddleware.GeneratorProvider(1, julian.New(2200100000)))
	e.Use(middleware.Logger())

	// Routes
	e.GET("/", handler.Generator)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
