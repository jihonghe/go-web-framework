package demo

import (
	"github.com/jihonghe/go-web-framework/summer"
)

type DemoServiceProvider struct{}

func (d *DemoServiceProvider) Name() string {
	return Key
}

func (d *DemoServiceProvider) Register(container summer.Container) summer.NewInstance {
	return NewDemoService
}

func (d *DemoServiceProvider) Boot(container summer.Container) error {
	println("demo service boot")
	return nil
}

func (d *DemoServiceProvider) IsDefer() bool {
	return true
}

func (d *DemoServiceProvider) Params(container summer.Container) []interface{} {
	// 将DemoServiceProvider.Params()的参数返回即可
	return []interface{}{container}
}
