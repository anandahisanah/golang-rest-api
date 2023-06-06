package controllers

import (
	"assignment-2/database"
	"assignment-2/models"
	"encoding/json"
	"fmt"
	_ "fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	OrderedAt    time.Time           `json:"orderedAt"`
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
		"code":    200,
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
		"Data": orders,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func UpdateOrderAndItems(c *gin.Context) {
	db := database.GetDB()

	// params
	orderId := c.Param("OrderId")

	// check order by orderId
	var order models.Order
	errFindOrder := db.Preload("Items").First(&order, orderId).Error
	if errFindOrder != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Order not found",
		})
		return
	}

	var CreateOrderRequest CreateOrderRequest
	errJson := c.ShouldBindJSON(&CreateOrderRequest)
	if errJson != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON data",
		})
		return
	}

	// update order

	// order.OrderedAt = CreateOrderRequest.OrderedAt
	// order.CustomerName = CreateOrderRequest.CustomerName
	// errUpdateOrder := db.Save(&order).Error

	errUpdateOrder := db.Model(&order).Where("order_id = ?", orderId).Updates(models.Order{
		CustomerName: CreateOrderRequest.CustomerName,
		OrderedAt:    CreateOrderRequest.OrderedAt,
	}).Error

	if errUpdateOrder != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update Order",
		})
		return
	}

	// update item
	for _, itemRequest := range CreateOrderRequest.Items {
		var item models.Item
		errFindItem := db.First(&item, itemRequest.ItemCode).Error
		if errFindItem != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Item not found",
			})
			return
		}

		item.ItemCode = itemRequest.ItemCode
		item.Description = itemRequest.Description
		item.Quantity = itemRequest.Quantity

		errSaveItem := db.Save(&item).Error
		if errSaveItem != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to update Item",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Order updated successfully",
		"data":    order,
	})
}

func DeleteOrderAndItems(c *gin.Context) {
	db := database.GetDB()
	w := c.Writer

	// params
	orderId := c.Param("OrderId")

	// find order by orderId
	err := db.Where("order_id = ?", orderId).First(&models.Order{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Order with id %s Not Found", orderId),
		})
		return
	}

	// Begin transaction
	tx := db.Begin()

	// delete items
	errDeleteItems := tx.Where("order_id = ?", orderId).Delete(&models.Item{}).Error
	if errDeleteItems != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete Items",
		})
		return
	}

	// delete order
	errDeleteOrder := tx.Where("order_id = ?", orderId).Delete(&models.Order{}).Error

	if errDeleteOrder != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete Order",
		})
		return
	}

	// commit transaction
	errCommit := tx.Commit().Error
	if errCommit != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to commit transaction",
		})
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, _ := json.Marshal(gin.H{
		"code":    200,
		"status":  "success",
		"message": "Delete Success",
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
