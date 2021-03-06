package interfaces

//go:generate ../../gomockhandler -source=$GOFILE -destination=../mock/$GOFILE

type User interface {
	String() string
	String2() string
}
