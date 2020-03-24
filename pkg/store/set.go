package store

type Set interface {
	NiceContainer

	// direct put the element to set
	Put(interface{}) error
	// random get a element from the set
	Get() (interface{}, error)
	// get all element from the set
	List() ([]interface{}, error)
}
