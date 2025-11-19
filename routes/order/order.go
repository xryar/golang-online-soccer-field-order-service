package routes

import (
	"order-service/clients"
	"order-service/constants"
	controllers "order-service/controllers/http"
	"order-service/middlewares"

	"github.com/gin-gonic/gin"
)

type OrderRoute struct {
	controller controllers.IRegistryController
	client     clients.IRegistryClient
	group      *gin.RouterGroup
}

type IOrderRoute interface {
	Run()
}

func NewOrderRoute(group *gin.RouterGroup, controller controllers.IRegistryController, client clients.IRegistryClient) IOrderRoute {
	return &OrderRoute{
		controller: controller,
		client:     client,
		group:      group,
	}
}

func (or *OrderRoute) Run() {
	group := or.group.Group("/order")
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, or.client), or.controller.GetOrder().GetAllWithPagination)
	group.GET("/:uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, or.client), or.controller.GetOrder().GetByUUID)
	group.GET("/user", middlewares.CheckRole([]string{
		constants.Customer,
	}, or.client), or.controller.GetOrder().GetOrderByUserID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, or.client), or.controller.GetOrder().Create)
}
