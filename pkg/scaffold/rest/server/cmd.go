package server

import "github.com/spf13/cobra"

// author info
type author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// webservice info
// webservice name will named modelWebService
type ws struct {

	// used for open-api config and for swagger-ui
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`

	// src-code package name
	Pkg string `json:"pkg"`
	// the model name
	// if model name is User, the webservice name is UsersWebService
	Model string `json:"model"`
}

type option struct {
	Author author `json:"author"`
	Ws     ws     `json:"ws"`
}

func New() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "server",
		Short: "use to generate server code",
	}
	return cmd
}
