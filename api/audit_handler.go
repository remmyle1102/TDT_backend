package api

import (
	"TDT_backend/app"
	"TDT_backend/db"
	"TDT_backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type audit struct {
	TaskName     string
	HostList     []string
	PlaybookList []string
}

// StartAudit start a new audit
func StartAudit(c echo.Context) error {

	audit := new(audit)
	if err := c.Bind(audit); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// check if the selected ip is avaiable or not
	var hostList []string
	var portList []string
	for _, host := range audit.HostList {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusGatewayTimeout, err)
		}
		conn.Close()
		ipAddr := strings.Split(host, ":")
		hostList = append(hostList, ipAddr[0])
		portList = append(portList, ipAddr[1])
	}

	// write connected ip to hosts
	f, err := os.Create("/etc/ansible/hosts")
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	defer f.Close()
	if _, err := f.WriteString("[win]" + "\n"); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	for _, host := range hostList {
		// create report folder name ip
		err := app.CreateDirIfNotExist("/etc/ansible/temp/" + host)
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}

		if _, err = f.WriteString(host + "\n"); err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}

	}
	if _, err := f.WriteString("\n\n[win:vars]\nansible_user=TDT-Audit\nansible_password=TDT_audit@2019\nansible_connection=winrm\nansible_winrm_server_cert_validation=ignore\n "); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// Execute ansible audit
	for _, playbook := range audit.PlaybookList {
		if playbook == "auditDB" {
			for index, host := range hostList {
				dbHost := host
				dbPort := portList[index]
				dbUser := "tdtaudit"
				dbPass := "TDT_audit@2019"
				connection := fmt.Sprintf("sqlserver://%s:%s@%s:%s?&connection+timeout=30", dbUser, dbPass, dbHost, dbPort)
				dbConn, err := sql.Open("sqlserver", connection)
				if err != nil {
					logrus.Error(err)
				}
				defer func() {
					err := dbConn.Close()
					if err != nil {
						logrus.Error(err)
					}
				}()

				err = app.CreateDirIfNotExist("/etc/ansible/temp/" + host + "/Task_Instance_DB")
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusBadRequest, err)
				}

				data := new(models.AuditDBInstance)
				data, err = db.AuditTaskDBInstance(dbConn)
				if err != nil && err != sql.ErrNoRows {
					logrus.Error(err)
					return c.JSON(http.StatusBadRequest, err)
				}
				file, _ := json.MarshalIndent(data, "", " ")
				fileLocation := fmt.Sprintf("/etc/ansible/temp/%s/Task_Instance_DB/Task_Instance_DB.json", host)
				err = ioutil.WriteFile(fileLocation, file, 0644)
				if err != nil {
					logrus.Error(err)
					return c.JSON(http.StatusBadRequest, err)
				}
			}
		} else {
			cmd := exec.Command("ansible-playbook", playbook)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				logrus.Error(err)
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
	}

	// Move created report to task name folder
	folders, err := app.ListDirOrFile("/etc/ansible/temp/", 1)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = app.CreateDirIfNotExist("/etc/ansible/reports/" + audit.TaskName)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	for _, folder := range folders {
		dest := fmt.Sprintf("/etc/ansible/reports/%s/%s", audit.TaskName, folder)
		err = os.Rename("/etc/ansible/temp/"+folder, dest)
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	// Insert successful task to db
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	addBy := claims["userID"].(float64)
	location := fmt.Sprintf("/etc/ansible/reports/%s", audit.TaskName)
	if err = db.InsertAuditTask(audit.TaskName, location, int(addBy)); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, "Succesfull")
}
