package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
	"uidGenerator/handler"
	generatorMiddleware "uidGenerator/middleware"
	"uidGenerator/timeprovider"
	"uidGenerator/timeprovider/epoch"
	"uidGenerator/timeprovider/julian"
)

var (
	portNumber   = flag.Int("port", 1323, "Port number")
	workerId     = flag.Int64("workerId", 1, "Worker ID")
	timeProvider = flag.String("timeProvider", "epoch", "Time provider (julian or epoch)")
	offset       = flag.Int64("offset", 1420070400000, "Offset for the time provider")
)

func main() {
	flag.Parse()

	//Time provider
	var provider timeprovider.TimeProvider
	switch *timeProvider {
	case "epoch":
		provider = epoch.New(*offset)
	case "julian":
		provider = julian.New(*offset)
	default:
		panic("Unknown time provider")
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(generatorMiddleware.GeneratorProvider(*workerId, provider))
	e.Use(middleware.Logger())

	// Routes
	e.GET("/", handler.Generator)

	// Start server
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(*portNumber)))
}
