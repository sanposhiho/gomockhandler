package interfaces

//go:generate ../gomockhandler -source=$GOFILE -destination=../mock/$GOFILE

type user interface {
	String() string
	String2() string
}
