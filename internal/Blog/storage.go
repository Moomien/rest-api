package Blog

type Storage interface {
	LoadAll() ([]Blog, error)
	LoadById(id int) (Blog, error)
	Save(data Blog) error
	SaveById(id int, data Blog) error
	Delete(id int) error
}
