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
	db, err := OpenDB()
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	backend, err := backend.New(db)
	if err != nil {
		log.Fatal("Failed to initialize backend:", err)
	}

	router := gin.New()

	router.RedirectTrailingSlash = false
	router.RemoveExtraSlash = true
	router.HandleMethodNotAllowed = true

	router.GET("/", handlers.NewWelcomePage())
	router.GET("/:username", handlers.NewShowProfile(backend))

	// Make no difference between "/foo" and "/foo/".
	handler := middleware.RemoveTrailingSlashFromPath(router)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenDB() (*sql.DB, error) {
	config := mysql.Config{
		DBName: "fakebook",
		User:   "fakebook",
		Passwd: "password",
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
