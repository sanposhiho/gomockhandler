package interfaces

//go:generate ../../gomockhandler -destination=../mock/$GOFILE . User2

type User2 interface {
	String() string
	String2() string
	String3() string
}
