package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/project-domino/domino-go/api"
	"github.com/project-domino/domino-go/middleware"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/project-domino/domino-go/api/v1"
)

func main() {
	db := Config.Database.Open()
	if err := SetupDB(db); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Enable/disable gin's debug mode.
	if Config.HTTP.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create and set up router.
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	// Add routes.
	r.GET("/", api.Version)
	for version, routes := range api.AllVersionRoutes() {
		routes(db, r.Group("/"+version))
	}

	// Start serving.
	log.Println("Ready to serve on", Config.HTTP.ServeOn())
	if err := r.Run(Config.HTTP.ServeOn()); err != nil {
		log.Println(err)
	}
}
