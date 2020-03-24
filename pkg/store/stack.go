package store

type Stack interface {
	NiceContainer
	Push(interface{}) error
	Pop() (interface{}, error)
}
