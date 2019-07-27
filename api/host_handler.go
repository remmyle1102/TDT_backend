package api

import (
	"TDT_backend/db"
	"TDT_backend/models"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// AddHost handle add host request
func AddHost(c echo.Context) error {
	host := new(models.Host)
	if err := c.Bind(host); err != nil {
		logrus.Error(err)
		return err
	}
	name := host.Name
	ipAdd := host.IPAdd
	port := host.Port
	description := host.Description
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	addBy := claims["userID"].(float64)

	if err := db.InsertHost(name, ipAdd, description, int(addBy), port); err != nil {
		logrus.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, "Sucessful")
}

// RemoveHost handle remove user request
func RemoveHost(c echo.Context) error {
	id := c.Param("id")
	err := db.DeleteHost(id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// FetchHost fetch all host
func FetchHost(c echo.Context) error {
	hosts := make([]*models.Host, 0)
	hosts, err := db.GetAllHost()
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, hosts)
}
