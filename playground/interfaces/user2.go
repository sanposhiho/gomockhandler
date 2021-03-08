package interfaces

//go:generate ../../gomockhandler -project_root=/Users/sanposhiho/workspace/gomockhandler/playground -destination=../mock/$GOFILE . User2

type User2 interface {
	String() string
	String2() string
}
