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
	result := make([]*models.ReportData, 0)
	location := c.QueryParam("location")
	folders, err := app.ListDirOrFile(location, 1)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	for _, folder := range folders {
		folderLocation := fmt.Sprintf("%s/%s", location, folder)
		reportData := new(models.ReportData)
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
				err = json.Unmarshal([]byte(byteValue), &reportData.DBTaskData)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusInternalServerError, err)
				}

			} else {
				var fileData models.FileData
				err = json.Unmarshal([]byte(byteValue), &fileData.FileData)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusInternalServerError, err)
				}
				reportData.FileData = append(reportData.FileData, fileData)
			}

		}

		reportData.File = files
		reportData.Folder = folder
		result = append(result, reportData)
	}
	return c.JSON(http.StatusOK, result)
}
