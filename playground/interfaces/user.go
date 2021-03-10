package interfaces

// Since this is a playground, the gomockhandler is mentioned in the comment of go generate for testing config generation, but it is not necessary.
// You don't need go generate comments anymore with gomockhandler
//go:generate ../../gomockhandler -project_root=../ -source=$GOFILE -destination=../mock/$GOFILE

type User interface {
	String() string
}
