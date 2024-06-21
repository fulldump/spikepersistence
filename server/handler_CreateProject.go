package server

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	persistence "spikepersistence"
)

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateProject(ctx context.Context, input *CreateProjectRequest) (*persistence.ItemWithId, error) {

	newProject := &persistence.ItemWithId{
		Id: uuid.NewString(),
		Item: persistence.Item{
			Title:       input.Title,
			Description: input.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Tasks:       []*persistence.Task{},
		},
	}

	err := ContextPersistence(ctx).Put(ctx, newProject)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrorPersistenceWrite
	}

	return newProject, nil
}
