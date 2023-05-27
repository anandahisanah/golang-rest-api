package routes

import (
	"assignment-2/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.POST("/order", controllers.CreateOrderAndItems)

	router.GET("/order", controllers.GetOrderAndItems)

	router.PATCH("/order/:OrderId", controllers.UpdateOrderAndItems)

	return router
}