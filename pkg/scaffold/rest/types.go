package rest

import (
	"fmt"
	"path"
	"strings"
)

const resourceTypeSuffix = "Manager"

type Package struct {
	Name string
}

func NewPackage(name string) *Package {
	return &Package{Name: name}
}

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

func NewService(kind string) *Service {
	s := &Service{
		Kind: kind,
	}

	s.Complete()
	return s
}

// Complete
// set default value for Service
func (s *Service) Complete() {

	s.Type = fmt.Sprintf("%s%s", s.Kind, resourceTypeSuffix)
	s.Title = fmt.Sprintf("%sService", UpperKind(s.Type))
	s.Description = fmt.Sprintf("resource for managing %s", s.Kind)
	s.Version = "v1.0"

	// best practice is /apis/{apiversion}/{kind}
	s.RootURLPrefix = path.Join("/api", s.Version, s.Kind)
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

func NewAuthor(name string, email string, url string) *Author {
	a := &Author{}
	a.Complete()
	return a
}

func (a *Author) Complete() {
	a.Name = "scott.wang"
	a.Email = "scottwangsxll@gmail.com"
	a.URL = "https://github.com/sxllwx"
}

type Model struct {
	Name string
}

func NewModel(name string) *Model {

	a := &Model{
		Name: name,
	}
	return a
}

func UpperKind(kind string) string {
	exportPrefix := strings.ToUpper(string(kind[0]))
	return exportPrefix + string(kind[1:])
}
