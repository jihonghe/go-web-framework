package main

import (
	"time"

	"go-web-framework/summer"
)

func registerRouter(core *summer.Core) {
	// 静态路由
	core.Get("/user/login", summer.TimeoutHandler(UserLoginController, time.Second*2))
	// 路由组
	group := core.Group("/subject")
	{
		group.Get("/list", SubjectController{}.List)
		// 动态路由
		group.Delete("/:id", SubjectController{}.Delete)
		group.Get("/:id", SubjectController{}.Get)
		group.Get("/:id/name", SubjectController{}.GetName)
		group.Post("/:id", SubjectController{}.Add)
	}
}
