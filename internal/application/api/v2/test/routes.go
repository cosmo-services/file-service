package test_api

import (
	_ "main/docs"

	"main/pkg"
)

type TestRoutes struct {
	handler    pkg.RequestHandler
	controller *TestController
}

func NewTestRoutes(handler pkg.RequestHandler, controller *TestController) *TestRoutes {
	return &TestRoutes{
		handler:    handler,
		controller: controller,
	}
}

func (r *TestRoutes) Setup() {
	group := r.handler.Gin.Group("/api/v2/")

	group.GET("/test", r.controller.GetOne)
}
