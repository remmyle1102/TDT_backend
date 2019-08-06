package db

import (
	"TDT_backend/models"

	"github.com/sirupsen/logrus"
)

// GetAllReport get report from DB
func GetAllReport() ([]*models.Report, error) {
	reports := make([]*models.Report, 0)

	query := "SELECT r.id, r.name , r.location, r.date, a.username from Report r INNER JOIN Account a ON r.addBy = a.id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		report := new(models.Report)
		err := rows.Scan(&report.ID, &report.Name, &report.Location, &report.Date, &report.AddBy)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, err
}
