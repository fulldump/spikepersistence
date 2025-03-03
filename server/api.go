package server

import (
	"github.com/fulldump/box"

	persistence "spikepersistence"
)

func NewApi(p persistence.Persistencer) *box.B {

	b := box.NewBox()

	b.Use(
		MiddlewareAccessLog(),
		MiddlewareSetPersistence(p),
	)

	b.Handle("POST", "/projects", CreateProject)
	b.Handle("GET", "/projects", ListProjects)
	b.Handle("GET", "/projects/{project_id}", GetProject)
	b.Handle("DELETE", "/projects/{project_id}", DeleteProject)
	b.Handle("PATCH", "/projects/{project_id}", ModifyProject)

	b.Handle("GET", "/openapi.json", OpenApi(b)).WithName("OpenApi")

	return b
}
