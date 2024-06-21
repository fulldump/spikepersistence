package server

import (
	"errors"
)

var ErrorPersistenceWrite = errors.New("persitence failed to write")
var ErrorPersistenceRead = errors.New("persitence failed to read")
var ErrorNotFound = errors.New("not found")
