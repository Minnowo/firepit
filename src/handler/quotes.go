package handler

import (
	"net/http"

	"github.com/EZCampusDevs/firepit/database"
	"github.com/EZCampusDevs/firepit/database/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GETCreateNewQuote(c echo.Context) error {

	var quote models.Quote

	if err := c.Bind(&quote); err != nil {

		log.Error(err)

		return echo.NewHTTPError(http.StatusBadRequest, "Could not read input")
	}

	result := database.GetDB().Create(&quote)

	if result.Error != nil {

		log.Error(result.Error)

		return echo.NewHTTPError(http.StatusInternalServerError, "There was a server error")
	}

	return c.String(http.StatusOK, "Quote Created")
}

func GETRandomQuote(c echo.Context) error {

	var quote models.Quote

	result := database.GetDB().Order("RAND()").Take(&quote)

	if result.Error != nil {
		log.Error(result.Error)
		return echo.NewHTTPError(http.StatusNotFound, "Could not get random quote")
	}

	return c.JSON(http.StatusOK, quote)
}
