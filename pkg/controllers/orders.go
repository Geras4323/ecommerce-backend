package controllers

import (
	"fmt"
	"math"
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/auth"
	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"gorm.io/gorm"
)

var OrderErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de los pedidos",

	"NotFound": "Orden no encontrada",
	"Create":   "Error al crear orden",
	"Update":   "Error al actualizar orden",
	"Delete":   "Error al eliminar orden",

	"Empty": "Esta orden no tiene productos",
}

// GET /api/v1/orders //////////////////////////////////////////////////////
func GetOrders(c echo.Context) error {
	orders := make([]models.Order, 0)

	if err := database.Gorm.
		Preload("User").
		Select("orders.*, COUNT(DISTINCT order_products.id) as products").
		Joins("INNER JOIN order_products ON orders.id = order_products.order_id").
		Group("orders.id").
		Order("orders.created_at DESC").
		Find(&orders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(OrderErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, orders)
}

// GET /api/v1/orders/my-orders
func GetOrdersByUser(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	var orders []models.Order

	if err := database.Gorm.
		Select("orders.*, COUNT(DISTINCT order_products.id) as products").
		Joins("INNER JOIN order_products ON orders.id = order_products.order_id").
		Where("orders.user_id = ?", c.User.ID).
		Order("created_at DESC").
		Group("orders.id").
		Find(&orders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(OrderErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, orders)
}

// GET /api/v1/orders/:id
func GetOrder(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	orderID := c.Context.Param("id")

	var order models.Order
	if err := database.Gorm.Preload("Payments", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = 'accepted'").Order("created_at DESC")
	}).Preload("User").Preload("OrderProducts.Product.Images").First(&order, orderID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(OrderErrors[utils.NotFound], err.Error()))
	}

	if order.UserID != c.User.ID && c.User.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, utils.SCTMake(OrderErrors[utils.NotFound], "inexistent order")) // Order does not belong to user - to not tell this order exists
	}

	return c.JSON(http.StatusOK, order)
}

// POST /api/v1/orders //////////////////////////////////////////////////////
func CreateOrder(baseContext echo.Context) error {
	c := baseContext.(*auth.AuthContext)

	items := []models.NewOrderProduct{}
	c.Bind(&items)

	if len(items) == 0 {
		return c.JSON(http.StatusForbidden, utils.SCTMake(OrderErrors["Empty"], "no products in order"))
	}

	if err := database.Gorm.First(&models.User{}, c.User.ID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(AuthErrors[utils.NotFound], err.Error()))
	}

	emailItems := make([]map[string]interface{}, len(items))
	orderProducts := make([]models.OrderProduct, len(items))
	var total float64 = 0

	order := models.Order{
		UserID: uint(c.User.ID),
		Total:  0,
	}

	// BEGIN TRANSACTION
	txError := database.Gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(OrderErrors[utils.Create], err.Error()))
		}

		for i, item := range items {
			var prod models.Product
			var unit models.Unit
			if err := database.Gorm.First(&prod, item.ProductID).Error; err != nil {
				return c.JSON(http.StatusNotFound, utils.SCTMake(fmt.Sprintf("Producto %d no encontrado", item.ProductID), err.Error()))
			}
			if err := database.Gorm.First(&unit, item.UnitID).Error; err != nil {
				return c.JSON(http.StatusNotFound, utils.SCTMake(fmt.Sprintf("Unidad %d no encontrada", item.UnitID), err.Error()))
			}

			total += float64(item.Quantity) * unit.Price

			emailItems[i] = map[string]interface{}{
				"article":  prod.Name,
				"quantity": item.Quantity,
				"unit":     unit.Unit,
				"price":    math.Round(unit.Price*100) / 100,
				"total":    math.Round(float64(item.Quantity)*unit.Price*100) / 100,
			}

			orderProduct := models.OrderProduct{
				OrderID:  order.ID,
				Product:  prod,
				Unit:     unit,
				Quantity: item.Quantity,
			}

			orderProducts = append(orderProducts, orderProduct)
			if err := tx.Create(&orderProduct).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, utils.SCTMake(fmt.Sprintf("Error al agregar producto %d a la orden", orderProduct.ID), err.Error()))
			}
		}

		if err := tx.Model(&models.Order{}).Where("id = ?", order.ID).Update("total", total).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, utils.SCTMake(AuthErrors[utils.Update], err.Error()))
		}

		return nil
	})

	if txError != nil {
		return txError
	}

	// CLEAR SHOPPING CART
	if err := database.Gorm.Where("user_id = ?", c.User.ID).Unscoped().Delete(&models.CartItem{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(CartErrors["Clear"], err.Error()))
	}

	// SEND CONFIRMATION EMAIL
	variables := map[string]interface{}{
		"name":  c.User.Name,
		"items": emailItems,
		"total": order.Total,
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cloud.DefaultSender.Email,
				Name:  cloud.DefaultSender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: c.User.Email,
				},
			},
			Subject:          "Confirmación de pedido",
			TemplateLanguage: true,
			TemplateID:       4839487,
			Variables:        variables,
		},
	}

	_ = messagesInfo
	// _, err := cloud.SendMail(messagesInfo)
	// if err != nil {
	// 	c.JSON(http.StatusBadGateway, utils.SCTMake(utils.CommonErrors[utils.Email], err.Error()))
	// }

	order.OrderProducts = append(order.OrderProducts, orderProducts...)
	return c.JSON(http.StatusCreated, order)
}

// PATCH /api/v1/orders/:id/state //////////////////////////////////////////////////////
func UpdateOrderState(c echo.Context) error {
	orderID := c.Param("id")

	var newOrder models.UpdateOrder
	c.Bind(&newOrder)

	var oldOrder models.Order
	if err := database.Gorm.First(&oldOrder, orderID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(OrderErrors[utils.NotFound], err.Error()))
	}

	oldOrder.State = newOrder.State

	if err := database.Gorm.Where("id = ?", orderID).Save(&oldOrder).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(OrderErrors[utils.Update], err.Error()))
	}

	return c.JSON(http.StatusOK, oldOrder)
}

// DELETE /api/v1/orders/:id //////////////////////////////////////////////////////
func DeleteOrder(c echo.Context) error {
	orderID := c.Param("id")

	if err := database.Gorm.First(&models.Order{}, orderID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(OrderErrors[utils.NotFound], err.Error()))
	}

	if err := database.Gorm.Delete(&models.Order{}, orderID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(OrderErrors[utils.Delete], err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]any{"id": orderID, "message": "Order deleted successfully"})
}
