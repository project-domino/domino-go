package search

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/models"
)

// Users returns all users that match a given query
func Users(db *gorm.DB, q string, items int, page int) ([]models.User, error) {
	var users []models.User

	searchQuery, err := ParseQuery(q)
	if err != nil {
		return users, err
	}
	// qSelectors := searchQuery.Selectors
	qText := strings.Join(searchQuery.Text, " & ")

	if err := db.Where(queryFormat, qText).
		Find(&users).
		Limit(items).
		Offset(page * items).
		Error; err != nil {
		return users, err
	}
	return users, nil
}
