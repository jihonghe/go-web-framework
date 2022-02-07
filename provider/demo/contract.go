package demo

const (
	Key = "summer:demo"
)

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}
