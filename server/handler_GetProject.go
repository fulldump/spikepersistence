package server

import (
	"context"
	"log"

	"github.com/fulldump/box"

	persistence "spikepersistence"
)

func GetProject(ctx context.Context) (*persistence.ItemWithId, error) {

	project_id := box.GetUrlParameter(ctx, "project_id")
	item, err := ContextPersistence(ctx).Get(ctx, project_id)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrorPersistenceRead
	}
	if item == nil {
		return nil, ErrorNotFound
	}

	return item, nil
}
