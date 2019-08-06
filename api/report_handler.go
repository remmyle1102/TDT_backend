package api

import (
	"TDT_backend/app"
	"TDT_backend/db"
	"TDT_backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func FetchReport(c echo.Context) error {
	reportList := make([]*models.Report, 0)
	reportList, err := db.GetAllReport()
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, reportList)
}

func FetchReportData(c echo.Context) error {
	result := make([]*models.ReportTableData, 0)
	location := c.QueryParam("location")

	hosts, err := app.ListDirOrFile(location, 1)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// add ReportTableData
	for _, host := range hosts {
		reportTableData := new(models.ReportTableData)
		reportDataS := make([]*models.ReportData, 0)
		hostLocation := fmt.Sprintf("%s/%s", location, host)
		folders, err := app.ListDirOrFile(hostLocation, 1)
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
		for _, folder := range folders {
			reportData := new(models.ReportData)

			if folder == "Task_Instance_DB" || folder == "Task_OS" {
				// add Report Table suggestions
				suggestionFile, err := os.Open("/etc/ansible/suggestions/" + folder + ".json")
				if err != nil {
					logrus.Error(err)
				}
				defer suggestionFile.Close()
				byteValue, err := ioutil.ReadAll(suggestionFile)
				if err != nil {
					logrus.Error(err)
				}
				err = json.Unmarshal(byteValue, &reportData.TableSuggestion)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusInternalServerError, err)
				}

				// add checkError
				checkErrorFile, err := os.Open("/etc/ansible/checkErrors/" + folder + ".json")
				if err != nil {
					logrus.Error(err)
				}
				defer checkErrorFile.Close()
				byteValue, err = ioutil.ReadAll(checkErrorFile)
				if err != nil {
					logrus.Error(err)
				}
				err = json.Unmarshal(byteValue, &reportData.CheckError)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusInternalServerError, err)
				}
			}

			folderLocation := fmt.Sprintf("%s/%s", hostLocation, folder)

			files, err := app.ListDirOrFile(folderLocation, 2)
			if err != nil {
				logrus.Error(err)
				return c.JSON(http.StatusBadRequest, err)
			}
			for _, file := range files {
				fileLocation := fmt.Sprintf("%s/%s", folderLocation, file)
				jsonFile, err := os.Open(fileLocation)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusBadRequest, err)
				}
				defer jsonFile.Close()
				byteValue, err := ioutil.ReadAll(jsonFile)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusInternalServerError, err)
				}

				if file == "Task_Instance_DB.json" {
					err = json.Unmarshal(byteValue, &reportData.DBTaskData)
					if err != nil {
						logrus.Error(err)
						return c.JSON(http.StatusInternalServerError, err)
					}

				} else {
					var fileData models.FileData
					err = json.Unmarshal(byteValue, &fileData.FileData)
					if err != nil {
						logrus.Error(err)
						return c.JSON(http.StatusInternalServerError, err)
					}
					reportData.FileData = append(reportData.FileData, fileData)
				}

			}

			reportData.File = files
			reportData.Folder = folder
			reportDataS = append(reportDataS, reportData)
		}
		reportTableData.Host = host
		reportTableData.ReportData = reportDataS
		result = append(result, reportTableData)
	}

	return c.JSON(http.StatusOK, result)
}
