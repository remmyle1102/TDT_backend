package db

import (
	"TDT_backend/models"

	"github.com/sirupsen/logrus"
)

func InsertHost(name, ipAdd, description string, addBy, port int) error {
	// playbook := new(models.Playbook)
	query := "INSERT INTO Host (name, ipAdd, port, description, addBy) VALUES (@p1,@p2,@p3,@p4,@p5)"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = stmt.Exec(name, ipAdd, port, description, addBy)
	return err
}

func DeleteHost(id string) error {
	query := "DELETE FROM Host WHERE id = @p1"
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

func GetAllHost() ([]*models.Host, error) {
	hosts := make([]*models.Host, 0)

	query := "SELECT h.id, h.name, h.ipAdd, h.port, h.description, a.username   from Host h INNER JOIN Account a ON h.addBy = a.id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		host := new(models.Host)
		err := rows.Scan(&host.ID, &host.Name, &host.IPAdd, &host.Port, &host.Description, &host.AddBy)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		hosts = append(hosts, host)
	}

	return hosts, err
}
