package api

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type Audit struct {
	HostList     []string
	PlaybookList []string
}

func StartAudit(c echo.Context) error {

	audit := new(Audit)
	if err := c.Bind(audit); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	fmt.Println(audit.PlaybookList)

	// check if the selected ip is avaiable or not
	var hostList []string
	for _, host := range audit.HostList {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			logrus.Error(err)
		} else {
			conn.Close()
			fmt.Println(host)
			ipAddr := strings.Split(host, ":")
			hostList = append(hostList, ipAddr[0])
		}

	}

	// write connected ip to hosts
	f, err := os.Create("ansible/hosts")
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
		if _, err := f.WriteString(host + "\n"); err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusBadRequest, err)
		}
	}

	return c.JSON(http.StatusOK, "Succesfull")
}
