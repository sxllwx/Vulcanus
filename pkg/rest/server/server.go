package server

//go:generate code-gen -go-type User  -rest-type users
type User struct {
	Name string `code-gen:"key"`
	Age  string
}

func (u *User)GetName()string{
	return ""
}
