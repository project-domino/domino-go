package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/project-domino/domino-go/db"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/handlers"
	"github.com/project-domino/domino-go/handlers/api"
	"github.com/project-domino/domino-go/handlers/redirect"
	"github.com/project-domino/domino-go/middleware"
	"github.com/project-domino/domino-go/models"
)

func main() {
	db.Open(Config.Database)
	if err := db.Setup(); err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

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
	Must(SetupAssets(r))

	// Routes that require user object
	m := r.Group("/")
	m.Use(middleware.Login())

	// Authentication Routes
	m.GET("/login", handlers.Simple("login.html"))
	m.GET("/register", handlers.Simple("register.html"))
	m.POST("/login", api.Login)
	m.POST("/register", api.Register)
	m.POST("/logout", api.Logout)

	// View Routes
	m.GET("/", handlers.Simple("home.html"))

	m.Group("/account",
		middleware.RequireAuth()).
		GET("/", redirect.Account).
		GET("/profile",
			middleware.AddPageName("profile"),
			handlers.Simple("account-profile.html")).
		GET("/security",
			middleware.AddPageName("security"),
			handlers.Simple("account-security.html")).
		GET("/notifications",
			middleware.AddPageName("notifications"),
			handlers.Simple("account-notifications.html"))

	m.GET("/search/:searchType",
		middleware.LoadSearchItems(),
		middleware.LoadSearchVars(),
		handlers.Simple("search.html"))

	m.Group("/u/:username",
		middleware.LoadUser("Notes", "Collections", "Notes.Tags", "Collections.Tags")).
		GET("/", redirect.User).
		GET("/notes", handlers.Simple("user-notes.html")).
		GET("/collections", handlers.Simple("user-collections.html"))

	m.Group("/note",
		middleware.LoadNote("Author", "Tags"),
		middleware.VerifyNotePublic()).
		GET("/:noteID", handlers.Simple("individual-note.html")).
		GET("/:noteID/:note-name", handlers.Simple("individual-note.html"))

	m.Group("/collection",
		middleware.LoadCollection("Author", "Tags"),
		middleware.VerifyCollectionPublic()).
		GET("/:collectionID",
			handlers.Simple("collection.html")).
		GET("/:collectionID/note/:noteID",
			middleware.LoadNote("Author", "Tags"),
			handlers.Simple("collection-note.html")).
		GET("/:collectionID/note/:noteID/:noteName",
			middleware.LoadNote("Author", "Tags"),
			handlers.Simple("collection-note.html"))

	m.Group("/writer-panel",
		middleware.RequireAuth(),
		middleware.RequireUserType(models.Writer, models.Admin),
		middleware.LoadRequestUser("Notes", "Collections")).
		GET("/", redirect.WriterPanel).
		GET("/note",
			middleware.AddPageName("new-note"),
			handlers.Simple("new-note.html")).
		GET("/note/:noteID/edit",
			middleware.LoadNote("Author", "Tags"),
			middleware.VerifyNoteOwner(),
			handlers.Simple("edit-note.html")).
		GET("/collection",
			middleware.AddPageName("new-collection"),
			handlers.Simple("new-collection.html")).
		GET("/collection/:collectionID/edit",
			middleware.LoadCollection("Author", "Tags"),
			middleware.VerifyCollectionOwner(),
			handlers.Simple("edit-collection.html")).
		GET("/tag",
			middleware.AddPageName("new-tag"),
			handlers.Simple("new-tag.html"))

	// API
	m.GET("/api/version", func(c *gin.Context) {
		handlers.RenderData(c, "debug.html", "data", map[string]interface{}{
			"currentVersion": "v1",
			"versions": []string{
				"v1",
			},
		})
	})
	m.Group("/api/v1").
		GET("/search/:searchType",
			middleware.LoadSearchItems(),
			api.Search).
		POST("/note",
			middleware.RequireAuth(),
			middleware.RequireUserType(models.Writer, models.Admin),
			api.NewNote).
		PUT("/note/:noteID",
			middleware.RequireAuth(),
			middleware.RequireUserType(models.Writer, models.Admin),
			api.EditNote).
		POST("/collection",
			middleware.RequireAuth(),
			middleware.RequireUserType(models.Writer, models.Admin),
			api.NewCollection).
		PUT("/collection/:collectionID",
			middleware.RequireAuth(),
			middleware.RequireUserType(models.Writer, models.Admin),
			api.EditCollection).
		POST("/tag",
			middleware.RequireAuth(),
			middleware.RequireUserType(models.Writer, models.Admin),
			api.NewTag)

	// Debug Routes
	if Config.HTTP.Debug {
		m.Group("/debug").
			GET("/editor", handlers.Simple("editor.html")).
			GET("/error", func(c *gin.Context) {
				errors.Debug.Apply(c)
			}).
			GET("/config", func(c *gin.Context) {
				handlers.RenderData(c, "debug.html", "data", Config)
			}).
			GET("/new/note", handlers.Simple("new-note.html"))
	}

	// Start serving.
	Must(r.Run(fmt.Sprintf(":%d", Config.HTTP.Port)))
}