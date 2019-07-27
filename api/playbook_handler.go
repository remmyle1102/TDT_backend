package api

import (
	"TDT_backend/db"
	"TDT_backend/models"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func FetchPlaybook(c echo.Context) error {
	playbookList := make([]*models.Playbook, 0)
	playbookList, err := db.GetAllPlaybook()
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, playbookList)
	}
	return c.JSON(http.StatusOK, playbookList)
}

func UploadPlaybook(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	addBy := claims["userID"].(float64)

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	name := c.FormValue("name")
	if name == "" {
		name = file.Filename
	} else {
		name = name + ".yml"
	}
	description := c.FormValue("description")

	//	Source
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	location := fmt.Sprintf("ansible/playbook/%s", name)
	dst, err := os.Create(location)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	err = db.InsertPlaybook(name, description, location, int(addBy))
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Successful")
}

func UpdatePlaybook(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")
	description := c.FormValue("description")

	err := db.UpdatePlaybook(id, name, description)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func RemovePlaybook(c echo.Context) error {
	playbook := new(models.Playbook)
	if err := c.Bind(playbook); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	id := playbook.ID
	location := playbook.Location
	fmt.Println(location)
	err := db.DeletePlaybook(id)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	err = os.Remove(location)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.NoContent(http.StatusNoContent)
}
