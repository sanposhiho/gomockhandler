package mockgen

type Runner interface {
	Run() error

	SetSource(new string)
	SetDestination(new string)
}
