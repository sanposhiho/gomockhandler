package interfaces

//go:generate ../../gomockhandler -project_root=/Users/sanposhiho/workspace/gomockhandler/playground -source=$GOFILE -destination=../mock/$GOFILE

type User interface {
	String() string
}
