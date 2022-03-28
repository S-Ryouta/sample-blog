package factories

import (
	"github.com/S-Ryouta/sample-blog/models"
	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
	"log"
)

func CreateEntry(db *gorm.DB) models.Entry {
	entry := models.Entry{}

	entry.ID = faker.UUIDHyphenated()
	entry.Title = faker.TitleMale()
	entry.Description = faker.Word()
	entry.Body = faker.Sentence()

	result := db.Create(&entry)

	if result.Error != nil {
		log.Fatal("Failed to insert entry record. \n", result.Error)
	}

	return entry
}
