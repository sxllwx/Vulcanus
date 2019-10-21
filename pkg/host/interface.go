package host

type Interface interface {
	Execute(rootCmd string, args ...string) ([]byte, error)
}
