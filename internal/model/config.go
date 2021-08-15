package model

import "errors"

type Config struct {
	// key: destination
	Mocks map[string]*Mock `json:"mocks"`
}

func NewChunk() *Config {
	return &Config{Mocks: map[string]*Mock{}}
}

func (c *Config) PutMock(destination string, mock Mock) {
	c.Mocks[destination] = &mock
}

func (c *Config) DeleteMock(dest string) {
	delete(c.Mocks, dest)
}

var (
	ErrNotFound = errors.New("config is not found in config")
)

func (c *Config) Find(destination string) (*Mock, error) {
	if m, ok := c.Mocks[destination]; ok {
		return m, nil
	}
	return nil, ErrNotFound
}
