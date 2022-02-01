package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-web-framework/summer"
)

func UserLoginController(c *summer.Context) error {
	log.Printf("recv a new req, start executing work at time: %s", time.Now().String())
	time.Sleep(3 * time.Second)
	log.Printf("finished executing work at time: %s", time.Now().String())
	c.Json(http.StatusOK, "login success")
	return nil
}

type SubjectController struct{}

func (s SubjectController) List(c *summer.Context) error {
	return nil
}

func (s SubjectController) GetName(c *summer.Context) error {
	return nil
}

func (s SubjectController) Get(c *summer.Context) error {
	c.Json(http.StatusOK, fmt.Sprintf("this is %s", c.GetRequest().URL.Path))
	return nil
}

func (s SubjectController) Delete(c *summer.Context) error {
	return nil
}

func (s SubjectController) Add(c *summer.Context) error {
	return nil
}

func (s SubjectController) Update(c *summer.Context) error {
	return nil
}
