package main

import (
	"go-web-framework/summer"
)

func registerRouter(core *summer.Core) {
	core.Get("foo", FooControllerHandler)
}
