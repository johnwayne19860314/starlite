package features

type ServiceUrlMaker interface {
	Make(name string, runLevel string) string
}
