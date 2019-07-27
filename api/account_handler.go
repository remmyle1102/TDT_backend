package api

import (
	"TDT_backend/db"
	"TDT_backend/middlewares"
	"TDT_backend/models"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Private(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["userName"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

// Login validate user then create JWT
func Login(c echo.Context) error {
	user := new(models.Account)
	if err := c.Bind(user); err != nil {
		return err
	}
	username := user.Username
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	user, err := db.ValidateUser(username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusForbidden, echo.ErrForbidden)
		}
		return err
	}
	if user.Status == false {
		return c.JSON(http.StatusForbidden, echo.ErrForbidden)
	}

	// Check in your db if the user exists or not

	tokens, err := middlewares.GenerateTokenPair(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tokens)

}

// RemoveUser handle remove user request
func RemoveUser(c echo.Context) error {
	id := c.Param("id")
	err := db.DeleteUser(id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// FetchUser handle fetch user request
func FetchUser(c echo.Context) error {
	users := make([]*models.Account, 0)
	users, err := db.GetAllUser()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// AddUser handle add user request
func AddUser(c echo.Context) error {
	user := new(models.Account)
	if err := c.Bind(user); err != nil {
		logrus.Error(err)
		return err
	}
	fmt.Println(user)
	username := user.Username
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))
	roleID := user.RoleID

	if err := db.InsertUser(username, password, roleID); err != nil {
		logrus.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, "Sucessful")

}
