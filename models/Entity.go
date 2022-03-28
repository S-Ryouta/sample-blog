package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Entry struct {
	ID          string     `gorm:"primaryKey" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Body        string     `gorm:"type:text" json:"body"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func SelectEntries(db *gorm.DB) ([]Entry, error) {
	var entries []Entry
	result := db.Find(&entries)
	if result.Error == gorm.ErrRecordNotFound {
		return entries, fmt.Errorf("SelectEntries ErrRecordNotFound")
	}
	return entries, nil
}

func FindEntry(db *gorm.DB, id string) (Entry, error) {
	var entry Entry
	result := db.First(&entry, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return entry, fmt.Errorf("FindEntry ErrRecordNotFound")
	}

	return entry, nil
}

func AddOrUpdateEntry(db *gorm.DB, entry Entry) {
	resultData := Entry{}
	result := db.First(&resultData, "id = ?", entry.ID)
	if result.Error == gorm.ErrRecordNotFound {
		db.Create(&entry)
	} else {
		db.Save(&entry)
	}
}
