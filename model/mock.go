package model

type Mock struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	CheckSum    [16]byte `json:"checksum"`
}

func NewMock(source, destination string, checksum [16]byte) Mock {
	return Mock{
		Source:      source,
		Destination: destination,
		CheckSum:    checksum,
	}
}
