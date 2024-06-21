package server

import (
	"context"
	"log"

	"github.com/fulldump/box"
)

func DeleteProject(ctx context.Context) error {

	project_id := box.GetUrlParameter(ctx, "project_id")
	err := ContextPersistence(ctx).Delete(ctx, project_id)
	if err != nil {
		log.Println(err.Error())
		return ErrorPersistenceWrite
	}

	return nil
}
