package model

import "errors"

//easyjson:json
type Config struct {
	// key: destination
	Mocks map[string]*Mock `json:"mocks"`
}

func NewChunk() *Config {
	return &Config{Mocks: map[string]*Mock{}}
}

func (c *Config) PutMock(mock Mock) {
	c.Mocks[mock.Destination] = &mock
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
