package controllers

import (
	"assignment-2/database"
	"assignment-2/models"
	"encoding/json"
	"net/http"
	_ "time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	OrderedAt    string              `json:"orderedAt"`
	CustomerName string              `json:"customerName"`
	Items        []CreateItemRequest `json:"items"`
}

type CreateItemRequest struct {
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

func CreateOrderAndItems(c *gin.Context) {
	db := database.GetDB()
	w := c.Writer

	var request CreateOrderRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON data",
		})
		return
	}

	Order := models.Order{
		CustomerName: request.CustomerName,
		Items:        []models.Item{},
	}

	// create order
	errOrder := db.Create(&Order).Error

	if errOrder != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Error saving Order data",
		})
		return
	}

	// create item
	for _, itemRequest := range request.Items {
		Item := models.Item{
			ItemCode:    itemRequest.ItemCode,
			Description: itemRequest.Description,
			Quantity:    itemRequest.Quantity,
			OrderId:     Order.OrderId,
		}
		errItem := db.Create(&Item).Error
		if errItem != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Error saving Item data",
			})
			return
		}
		Order.Items = append(Order.Items, Item)
	}

	// response
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(gin.H{
		"status":  "success",
		"message": "Order created successfully",
		"data":    Order,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetOrderAndItems(c *gin.Context) {
	db := database.GetDB()
	w := c.Writer

	orders := []models.Order{}
	err := db.Preload("Items").Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data not found",
		})
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(gin.H{
		"status":  "success",
		"message": "Success",
		"data":    orders,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
