package model

import "errors"

//easyjson:json
type Chunk struct {
	// key: destination
	Mocks map[string]*Mock `json:"mocks"`
}

func NewChunk() *Chunk {
	return &Chunk{Mocks: map[string]*Mock{}}
}

func (c *Chunk) PutMock(mock Mock) {
	c.Mocks[mock.Destination] = &mock
}

var (
	ErrNotFound = errors.New("chunk is not found in chunk")
)

func (c *Chunk) Find(destination string) (*Mock, error) {
	if m, ok := c.Mocks[destination]; ok {
		return m, nil
	}
	return nil, ErrNotFound
}
