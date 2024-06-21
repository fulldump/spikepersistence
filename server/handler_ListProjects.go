package server

import (
	"context"
	"log"
)

type ListProjectItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func ListProjects(ctx context.Context) ([]*ListProjectItem, error) {

	items, err := ContextPersistence(ctx).List(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrorPersistenceRead
	}

	result := make([]*ListProjectItem, len(items))
	for i, item := range items {
		result[i] = &ListProjectItem{
			Id:    item.Id,
			Title: item.Title,
		}
	}

	return result, nil
}
