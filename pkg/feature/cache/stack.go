package cache

type Stack interface {
	NiceContainer
	Push(interface{}) error
	Pop() (interface{}, error)
}
