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

	// Integrate front end
	/*
		e.File("/", "frontend/build/index.html")
		e.Use(middleware.Static("frontend/build"))
	*/
	e.POST("/login", api.Login)
	// e.POST("/token", middlewares.RefreshToken)

	// restricted api
	r := e.Group("/api")
	r.Use(middlewares.IsLoggedIn)
	r.GET("/private", api.Private)
	r.POST("/upload-playbook", api.UploadPlaybook)
	r.GET("/fetch-playbook", api.FetchPlaybook)
	r.DELETE("/delete-playbook/:id", api.RemovePlaybook)

	r.POST("/add-user", api.AddUser)
	r.GET("/fetch-user", api.FetchUser)
	r.DELETE("/remove-user/:id", api.RemoveUser)

	r.POST("/add-host", api.AddHost)
	r.GET("/fetch-host", api.FetchHost)
	r.DELETE("/remove-host/:id", api.RemoveHost)

	r.POST("start-audit", api.StartAudit)

	lock := make(chan error)
	go func(lock chan error) { lock <- e.Start(serverAddress) }(lock)

	time.Sleep(1 * time.Millisecond)
	middlewares.MakeLogEntry(nil).Warning("application started without ssl/tls enabled")

	err = <-lock
	if err != nil {
		middlewares.MakeLogEntry(nil).Panic("failed to start application")
	}
}
