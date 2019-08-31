package server

import (
	"fmt"
	"path"
	"strings"
)

const (
	ResourceTypeSuffix = "Manager"
)

// the service tag for swagger
type Tag struct {
	Name        string
	Description string
}

type Service struct {

	// which kind of resource be managed by this web-service
	// eg: books, users, please set this filed to plural
	Kind string
	// this flied is resource + Manager
	Type string

	// url config
	RootURLPrefix string

	//  the service meta
	Title       string
	Description string
	Version     string
	Tag         *Tag
}

// Complete
// set default value for Service
func (s *Service) Complete() {

	s.Type = fmt.Sprintf("%s%s", s.Kind, ResourceTypeSuffix)
	s.Title = fmt.Sprintf("%sService", strings.ToUpper(s.Type))
	s.Description = fmt.Sprintf("resource for managing %s", s.Kind)
	s.Version = "v1.0.0"

	s.RootURLPrefix = path.Join("/", s.RootURLPrefix, s.Kind)
	s.Tag = &Tag{
		Name:        s.Kind,
		Description: fmt.Sprintf("Managing %s", s.Kind),
	}
}

// the service auth
type Author struct {
	Name  string
	Email string
	URL   string
}

func (a *Author) Complete() {
	a.Name = "scott.wang"
	a.Email = "scottwangsxll@gmail.com"
	a.URL = "https://github.com/sxllwx"
}

type Model struct {
	Name string
}
