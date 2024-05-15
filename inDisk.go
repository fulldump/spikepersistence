package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type InDisk struct {
	dataDir string
	cache   *InMemory
}

func NewInDisk(dataDir string) (*InDisk, error) {

	// ensure dir
	err := os.MkdirAll(dataDir, 0777)
	if err != nil {
		return nil, fmt.Errorf("ERROR: ensure data dir '%s': %s\n", dataDir, err.Error())
	}

	cache := NewInMemory()

	// load dir
	err = filepath.WalkDir(dataDir, func(filename string, d fs.DirEntry, walkErr error) error {

		if walkErr != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if !strings.EqualFold(".json", path.Ext(filename)) {
			return nil
		}

		f, err := os.Open(filename)
		if err != nil {
			log.Printf("error loading '%s': %s\n", filename, err.Error())
			return nil // todo. check if err should be returned
		}

		item := &ItemWithId{}
		err = json.NewDecoder(f).Decode(item)
		if err != nil {
			log.Printf("error decoding '%s': %s\n", filename, err.Error())
			return nil // todo. check if err should be returned
		}

		return cache.Put(context.Background(), item)
	})
	if err != nil {
		return nil, fmt.Errorf("load items: %s", err.Error())
	}

	return &InDisk{
		dataDir: dataDir,
		cache:   cache,
	}, nil
}

func (f *InDisk) List(ctx context.Context) ([]*ItemWithId, error) {
	return f.cache.List(ctx)
}

func (f *InDisk) Put(ctx context.Context, item *ItemWithId) error {

	filename := path.Join(f.dataDir, item.Id+".json")

	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	e := json.NewEncoder(fd)
	e.SetIndent("", "    ")
	err = e.Encode(item)
	if err != nil {
		return fmt.Errorf("persisting %s: %s\n", filename, err.Error())
	}

	return f.cache.Put(ctx, item)
}

func (f *InDisk) Get(ctx context.Context, id string) (*ItemWithId, error) {
	return f.cache.Get(ctx, id)
}

func (f *InDisk) Delete(ctx context.Context, id string) error {

	item, err := f.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("item '%s' does not exist", id)
	}

	filename := path.Join(f.dataDir, item.Id+".json")
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("item '%s' persistence error: %s", id, err.Error())
	}

	return f.cache.Delete(ctx, id)
}
