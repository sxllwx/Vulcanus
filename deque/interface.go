package deque

type StopErr string

func (s StopErr) Error() string {
	return "double end queue already stop"
}

func CheckStopError(err error) bool {
	_, is := err.(StopErr)
	return is
}

type Interface interface {
	Len() (int, error)

	Push(interface{}) error
	Shift() (interface{}, error)
	Pop() (interface{}, error)
	Done(interface{})
}

