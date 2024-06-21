package main

import (
	persistence "spikepersistence"
	"spikepersistence/server"
)

func main() {
	p := persistence.NewInMemory()
	s := server.NewApi(p)
	s.ListenAndServe()
}
