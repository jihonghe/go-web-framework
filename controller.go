package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jihonghe/go-web-framework/provider/demo"
	"github.com/jihonghe/go-web-framework/summer/gin"
)

func UserLoginController(c *gin.Context) {
	log.Printf("recv a new req, start executing work at time: %s", time.Now().String())
	time.Sleep(1 * time.Second)
	log.Printf("finished executing work at time: %s", time.Now().String())
	c.ISetStatusOk().IJson("login success")
}

type SubjectController struct{}

func (s SubjectController) List(c *gin.Context) {
	demoService := c.MustMake(demo.Key).(demo.Service)
	foo := demoService.GetFoo()
	c.ISetStatusOk().IJson(foo)
}

func (s SubjectController) GetName(c *gin.Context) {
}

func (s SubjectController) Get(c *gin.Context) {
	c.ISetStatusOk().IJson(fmt.Sprintf("this is %s", c.Request.URL.Path))
}

func (s SubjectController) Delete(c *gin.Context) {
}

func (s SubjectController) Add(c *gin.Context) {
}

func (s SubjectController) Update(c *gin.Context) {
}
