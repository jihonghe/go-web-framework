package main

import (
	"go-web-framework/summer"
	"go-web-framework/summer/middleware"
)

func registerRouter(core *summer.Core) {
	// 静态路由
	// core.Get("/user/login", middleware.TestMiddleware(), UserLoginController)
	// core.Get("/user/login", middleware.Timeout(1*time.Second), middleware.Cost, UserLoginController)
	core.Get("/user/login", middleware.Cost, UserLoginController)
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
