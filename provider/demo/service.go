package demo

import (
	"github.com/jihonghe/go-web-framework/summer"
)

type DemoService struct {
	// 需要实现Service接口
	Service

	container summer.Container
}

func NewDemoService(params ...interface{}) (interface{}, error) {
	container := params[0].(summer.Container)
	println("new demo service")
	return &DemoService{container: container}, nil
}

func (d *DemoService) GetFoo() Foo {
	return Foo{Name: "I am foo"}
}
