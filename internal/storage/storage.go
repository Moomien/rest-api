package storage

import (
	"blog-restapi/internal/models"
	"encoding/json"
	"os"
)

type Storage interface {
	Load() ([]models.Blog, error)
	Save(data models.Blog) error
}

type JsonStorage struct {
	filename string
}

func (j *JsonStorage) Save(data models.Blog) error {
	filename := j.filename + ".json"
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, d, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonStorage) Load() ([]models.Blog, error) {
	filename := j.filename + ".json"
	reader, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result []models.Blog
	err = json.Unmarshal(reader, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
