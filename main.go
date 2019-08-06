package main

import (
	"TDT_backend/api"
	"TDT_backend/db"
	"TDT_backend/middlewares"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	serverAddress := viper.GetString("server.address")

	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	connection := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=30", dbUser, dbPass, dbHost, dbPort, dbName)
	dbConn, err := sql.Open("sqlserver", connection)
	if err != nil {
		logrus.Error(err)
	}
	db.GetDBConnection(dbConn)

	err = dbConn.Ping()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	e := echo.New()
	e.Use(middlewares.Logging)
	e.HTTPErrorHandler = middlewares.ErrorHandler
	e.Use(middleware.CORS())

	e.POST("/login", api.Login)
	// e.POST("/token", middlewares.RefreshToken)

	// restricted api
	r := e.Group("/api")
	r.Use(middlewares.IsLoggedIn)
	r.GET("/private", api.Private)
	// Playbook API
	r.POST("/upload-playbook", api.UploadPlaybook)
	r.GET("/fetch-playbook", api.FetchPlaybook)
	r.DELETE("/delete-playbook", api.RemovePlaybook)

	// User API
	r.POST("/add-user", api.AddUser)
	r.GET("/fetch-user", api.FetchUser)
	r.DELETE("/remove-user/:id", api.RemoveUser)

	// Host API
	r.POST("/add-host", api.AddHost)
	r.GET("/fetch-host", api.FetchHost)
	r.DELETE("/remove-host/:id", api.RemoveHost)

	// Audit API
	r.POST("/start-audit", api.StartAudit)

	// Report API
	r.GET("/fetch-report", api.FetchReport)
	r.GET("/fetch-report-data", api.FetchReportData)

	lock := make(chan error)
	go func(lock chan error) { lock <- e.Start(serverAddress) }(lock)

	time.Sleep(1 * time.Millisecond)
	middlewares.MakeLogEntry(nil).Warning("application started without ssl/tls enabled")

	err = <-lock
	if err != nil {
		middlewares.MakeLogEntry(nil).Panic("failed to start application")
	}
}
