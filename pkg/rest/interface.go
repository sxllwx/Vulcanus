package rest

type Interface interface {
	GET() *request
	POST() *request
	DELETE() *request
	PUT() *request
}
