package chunk

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mailru/easyjson"
	"github.com/sanposhiho/gomockhandler/model"
)

const filename = "gomockhandler.json"

type Repository struct {
}

func NewRepository() Repository {
	return Repository{}
}

func (r *Repository) Put(m *model.Chunk) error {
	d, err := easyjson.Marshal(m)
	if err != nil {
		return fmt.Errorf("easyjson marshal: %w", err)
	}

	return ioutil.WriteFile(filename, d, 0644)
}

func (r *Repository) Get() (*model.Chunk, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var m model.Chunk
	err = easyjson.Unmarshal(raw, &m)
	if err != nil {
		return nil, fmt.Errorf("easyjson unmarshal: %w", err)
	}

	return &m, err
}
