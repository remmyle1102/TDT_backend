package db

import (
	"TDT_backend/models"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

func ValidateUser(username, password string) (*models.Account, error) {
	user := new(models.Account)
	query := "SELECT a.id, a.username, r.roleName, a.status from Account a INNER JOIN Role r ON a.roleId = r.id WHERE a.username = @p1 AND a.password = @p2"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
	}
	row := stmt.QueryRow(username, password)
	err = row.Scan(&user.ID, &user.Username, &user.Role, &user.Status)
	if err == sql.ErrNoRows {
		return user, echo.ErrUnauthorized
	}
	return user, err
}

func InsertUser(username, password string, roleID int) error {
	query := "INSERT INTO Account (username, password, roleID) VALUES (@p1,@p2,@p3)"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = stmt.Exec(username, password, roleID)
	return err
}

func DeleteUser(id string) error {
	query := "DELETE FROM Account WHERE id = @p1"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func GetAllUser() ([]*models.Account, error) {
	users := make([]*models.Account, 0)

	query := "SELECT a.id, a.username, r.roleName, a.status from Account a INNER JOIN Role r ON a.roleId = r.id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := new(models.Account)
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.Status)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}
