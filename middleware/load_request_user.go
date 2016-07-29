package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/project-domino/domino-go/db"
	"github.com/project-domino/domino-go/models"
)

// LoadRequestUser loads certain objects in the request context's user
func LoadRequestUser(objects ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Acquire user from the request context.
		user := c.MustGet("user").(models.User)

		// Set objects to be preloaded to db
		preloadedDB := db.DB.Where("id = ?", user.ID)
		for _, object := range objects {
			preloadedDB = preloadedDB.Preload(object)
		}

		// Query for user and set context
		var loadedUser models.User
		if err := preloadedDB.First(&loadedUser).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Set("user", loadedUser)

		c.Next()
	}
}