package interfaces

//go:generate ../../gomockhandler -destination=../mock/$GOFILE -check=true . User2

type User2 interface {
	String() string
}
