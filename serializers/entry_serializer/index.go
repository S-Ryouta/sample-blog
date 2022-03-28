package entry_serializer

import (
	"github.com/S-Ryouta/sample-blog/models"
	"time"
)

type Index struct {
	id          string
	title       string
	description string
	body        string
	createdAt   string
	updatedAt   string
}

func IndexSerializer(entries []models.Entry) []Index {
	var index []Index
	for _, entry := range entries {
		t := Index{
			id:          entry.ID,
			title:       entry.Title,
			description: entry.Description,
			body:        entry.Body,
			createdAt:   entry.CreatedAt.Format(time.RFC3339),
			updatedAt:   entry.CreatedAt.Format(time.RFC3339),
		}
		index = append(index, t)
	}

	return index
}
