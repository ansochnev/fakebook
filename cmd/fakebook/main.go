package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"fakebook/internal/backend"
	"fakebook/internal/handlers"
	"fakebook/internal/middleware"
)

func main() {
	config, err := ReadConfigFile("fakebook.yaml")
	if err != nil {
		log.Fatalln("Failed to read config file:", err)
	}

	db, err := OpenDB(config)
	if err != nil {
		log.Fatalln("Failed to open DB:", err)
	}
	defer db.Close()

	backend, err := backend.New(db)
	if err != nil {
		log.Fatal("Failed to initialize backend:", err)
	}

	gin.SetMode(gin.ReleaseMode)
	if config.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.RedirectTrailingSlash = false
	router.RemoveExtraSlash = true
	router.HandleMethodNotAllowed = true

	router.Use(gin.Logger())

	router.Static("/css", "site/css")

	welcomePage := handlers.WelcomePage{
		BasicURL: config.BasicURL(),
	}
	router.GET("/", welcomePage.Handle)

	router.GET("/:username", handlers.NewShowProfile(backend))

	// Make no difference between "/foo" and "/foo/".
	handler := middleware.RemoveTrailingSlashFromPath(router)

	if config.UseHTTPS {
		err = http.ListenAndServeTLS(config.ListenAddress,
			config.CertFile, config.KeyFile, handler)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = http.ListenAndServe(config.ListenAddress, handler)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func OpenDB(config *Config) (*sql.DB, error) {
	dbConfig := mysql.Config{
		DBName: config.MySQL.Database,
		User:   config.MySQL.User,
		Passwd: config.MySQL.Password,
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
