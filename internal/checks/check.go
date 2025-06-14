package checks

type Check interface {
	Run() (pass bool, err error)
	Name() string
}
