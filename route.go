package main

import (
	"github.com/jihonghe/go-web-framework/summer/gin"
	"github.com/jihonghe/go-web-framework/summer/middleware"
)

func registerRouter(core *gin.Engine) {
	// 静态路由
	// core.Get("/user/login", middleware.TestMiddleware(), UserLoginController)
	// core.Get("/user/login", middleware.Timeout(500*time.Millisecond), middleware.Cost, UserLoginController)
	core.GET("/user/login", middleware.Cost, UserLoginController)
	// 路由组
	group := core.Group("/subject")
	{
		group.GET("/list", SubjectController{}.List)
		// 动态路由
		group.DELETE("/:id", SubjectController{}.Delete)
		group.GET("/:id", SubjectController{}.Get)
		group.GET("/:id/name", SubjectController{}.GetName)
		group.POST("/:id", SubjectController{}.Add)
	}
}
