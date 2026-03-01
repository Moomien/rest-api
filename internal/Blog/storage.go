package Blog

import (
	"encoding/json"
	"errors"
	"os"
)

type Storage interface {
	LoadAll() ([]Blog, error)
	LoadById(id int) (Blog, error)
	Save(data Blog) error
	SaveById(id int, data Blog) error
	Delete(id int) error
}

type JsonStorage struct {
	filename string
	cache    map[int]Blog
}

func (j *JsonStorage) NewStorage(filename string) *JsonStorage {
	return &JsonStorage{
		filename: filename + ".json",
		cache:    make(map[int]Blog),
	}
}

func (j *JsonStorage) Delete(id int) error {
	if _, ok := j.cache[id]; !ok {
		return errors.New("blog not found")
	}
	delete(j.cache, id)

	var blogs []Blog
	for _, v := range j.cache {
		blogs = append(blogs, v)
	}
	bytes, err := json.MarshalIndent(blogs, " ", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(j.filename, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonStorage) SaveById(id int, data Blog) error {
	blogs, err := j.LoadAll()
	if err != nil {
		return err
	}

	found := false
	for _, v := range blogs {
		if v.ID == id {
			blogs[id] = v
			found = true
			break
		}
	}

	if !found {
		return errors.New("Not found data")
	}

	newData, err := json.MarshalIndent(blogs, " ", "  ")
	os.WriteFile(j.filename, newData, 0644)
	return nil
}

func (j *JsonStorage) Save(data Blog) error {
	if data.ID == 0 {
		maxID := 0
		for id := range j.cache {
			if id > maxID {
				maxID = id
			}
		}
		data.ID = maxID + 1
	}
	j.cache[data.ID] = data

	var blogs []Blog
	for _, v := range j.cache {
		blogs = append(blogs, v)
	}

	filedata, err := json.MarshalIndent(blogs, " ", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(j.filename, filedata, 0644)
}

func (j *JsonStorage) LoadById(id int) (Blog, error) {
	blogs, err := j.LoadAll()
	if err != nil {
		return Blog{}, err
	}
	for _, v := range blogs {
		j.cache[v.ID] = v
	}

	blog, ok := j.cache[id]
	if !ok {
		return Blog{}, errors.New("blog not found")
	}
	return blog, nil
}

func (j *JsonStorage) LoadAll() ([]Blog, error) {
	data, err := os.ReadFile(j.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var blogs []Blog
	if err := json.Unmarshal(data, &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}
