package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"uidGenerator/generator"
)

func Generator(c echo.Context) error {
	worker := c.Get("worker").(*generator.WorkerVariant)
	idn := c.QueryParam("numberOfIds")
	numberOfIds := 1
	if idn != "" {
		v, err := strconv.Atoi(idn)
		if err == nil {
			numberOfIds = v
		}
	}

	ids, err := worker.GenerateID(numberOfIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ids": ids,
	})
}
