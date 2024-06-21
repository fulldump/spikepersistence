package server

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/fulldump/box"

	persistence "spikepersistence"
)

type ModifyProjectRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func ModifyProject(ctx context.Context, input *ModifyProjectRequest) (*persistence.ItemWithId, error) {

	project_id := box.GetUrlParameter(ctx, "project_id")
	item, err := ContextPersistence(ctx).Get(ctx, project_id)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrorPersistenceRead
	}
	if item == nil {
		return nil, ErrorNotFound
	}

	// Modify object
	mergeLeft(item, input)
	item.UpdatedAt = time.Now()

	err = ContextPersistence(ctx).Put(ctx, item)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrorPersistenceWrite
	}

	return item, nil
}

func mergeLeft(a, b any) {
	bjson, _ := json.Marshal(b)
	json.Unmarshal(bjson, &a)
}
