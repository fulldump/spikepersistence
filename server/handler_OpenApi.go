package server

import (
	"encoding/json"
	"net/http"

	"github.com/fulldump/box"
	"github.com/fulldump/box/boxopenapi"
)

func OpenApi(b *box.B) any {
	// openapi
	// Openapi automatic spec
	spec := boxopenapi.Spec(b)
	spec.Info.Title = "Spike service"
	spec.Info.Description = "Spike service for learning purposes"
	spec.Info.Contact = &boxopenapi.Contact{
		Url: "https://github.com/fulldump/spikepersistence/issues/new",
	}
	spec.Servers = []boxopenapi.Server{
		{Url: "http://localhost:8080/"},
	}

	o, _ := json.MarshalIndent(spec, "", "    ")

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(o)
	}
}
