package db

import (
	"TDT_backend/models"

	"github.com/sirupsen/logrus"
)

func GetAllPlaybook() ([]*models.Playbook, error) {
	playbookList := make([]*models.Playbook, 0)

	query := "SELECT p.id, p.name, a.username, p.description, p.location from PlayBook p INNER JOIN Account a ON p.addBy = a.id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		playbook := new(models.Playbook)
		err := rows.Scan(&playbook.ID, &playbook.Name, &playbook.AddBy, &playbook.Description, &playbook.Location)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		playbookList = append(playbookList, playbook)
	}

	return playbookList, err
}

func InsertPlaybook(name, description, location string, addBy int) error {
	// playbook := new(models.Playbook)
	query := "INSERT INTO Playbook (name, addBy, description, location) VALUES (@p1,@p2,@p3,@p4)"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = stmt.Exec(name, addBy, description, location)
	return err
}

func UpdatePlaybook(id, name, description string) error {
	query := "UPDATE PlayBook SET name=@p1, description = @p2 WHERE id = @p3"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = stmt.Exec(name, description, id)
	return err
}

func DeletePlaybook(id int) error {
	query := "DELETE FROM PlayBook WHERE id = @p1"
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
