package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sanposhiho/gomockhandler/internal/model"
)

type Repository struct{}

func NewRepository() Repository {
	return Repository{}
}

func (r *Repository) Put(m *model.Config, path string) error {
	d, err := json.MarshalIndent(m, "", "	")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	return os.WriteFile(path, append(d, '\n'), 0644)
}

func (r *Repository) Get(path string) (*model.Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var m model.Config
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return &m, err
}
