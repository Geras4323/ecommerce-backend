package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/geras4323/ecommerce-backend/pkg/database"
	"github.com/geras4323/ecommerce-backend/pkg/models"
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

var PaymentErrors = map[string]string{
	"Internal": "Ocurrió un error durante la carga de los comprobantes",

	"NotFound": "Comprobante no encontrado",
	"Create":   "Error al crear comprobante",
	"Update":   "Error al actualizar comprobante",
	"Delete":   "Error al eliminar comprobante",
}

const (
	StatusAccepted string = "accepted"
	StatusPending  string = "pending"
	StatusRejected string = "rejected"

	PlatformMP         string = "mercadopago"
	PlatformAttachment string = "attachment"
)

// GET /api/v1/payments //////////////////////////////////////////////////////
func GetPayments(c echo.Context) error {
	payments := make([]models.Payment, 0)

	if err := database.Gorm.Find(&payments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(PaymentErrors[utils.Internal], err.Error()))
	}

	return c.JSON(http.StatusOK, payments)
}

// GET /api/v1/payments/:id ?statusOnly=
func GetPayment(c echo.Context) error {
	paymentID := c.Param("id")

	var payment models.Payment
	if err := database.Gorm.First(&payment, paymentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(PaymentErrors[utils.NotFound], err.Error()))
	}

	statusOnly := c.QueryParam("statusOnly") == "true"
	var returnValue interface{}

	if statusOnly {
		returnValue = payment.Status
	} else {
		returnValue = payment
	}

	return c.JSON(http.StatusOK, returnValue)

}

// POST /api/v1/payments //////////////////////////////////////////////////////
func CreatePayment(c echo.Context) error {
	orderID := c.Param("orderID")

	var order models.Order
	if err := database.Gorm.First(&order, orderID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(OrderErrors[utils.NotFound], err.Error()))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buff := bytes.NewBuffer(nil)
	io.Copy(buff, src)

	// Upload
	mimetype := http.DetectContentType(buff.Bytes())
	base64file := fmt.Sprintf("data:%s;base64,%s", mimetype, toBase64(buff.Bytes()))
	fileFolder := fmt.Sprintf("%s/%s", utils.GetEnvVar("CLOUDINARY_ENV_FOLDER"), "payments")
	fileName := uuid.New().String()
	filePath := fmt.Sprintf("%s/%s", fileFolder, fileName)

	res, err := cloud.Cld.Upload.Upload(cloud.Ctx, base64file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   fileFolder,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.SCTMake(utils.CommonErrors[utils.FileUpload], err.Error()))
	}

	// Save to DB
	var newPayment models.Payment
	newPayment.OrderID = order.ID
	newPayment.Url = null.StringFrom(res.SecureURL)
	newPayment.Path = null.StringFrom(filePath)
	newPayment.Status = StatusAccepted
	newPayment.Platform = PlatformAttachment

	if err := database.Gorm.Create(&newPayment).Error; err != nil {
		return c.JSON(http.StatusConflict, utils.SCTMake(PaymentErrors[utils.Create], err.Error()))
	}

	return c.JSON(http.StatusOK, newPayment)
}

// PUT /api/v1/payments/:id //////////////////////////////////////////////////////
func UpdatePayment(c echo.Context) error {
	// paymentID := c.Param("id")

	// var newPayment models.UpdatePayment
	// c.Bind(&newPayment)

	// var oldPayment models.Payment
	// if err := database.Gorm.First(&oldPayment, paymentID).Error; err != nil {
	// 	return c.JSON(http.StatusNotFound, utils.SCTMake(PaymentErrors[utils.NotFound], err.Error()))
	// }

	// oldPayment.OrderID = newPayment.OrderID

	// if err := database.Gorm.Where("id = ?", paymentID).Save(&oldPayment).Error; err != nil {
	// 	return c.JSON(http.StatusInternalServerError, utils.SCTMake(PaymentErrors[utils.Update], err.Error()))
	// }

	// return c.JSON(http.StatusOK, oldPayment)
	return nil
}

// DELETE /api/v1/payments/:id //////////////////////////////////////////////////////
func DeletePayment(c echo.Context) error {
	paymentID := c.Param("id")

	var payment models.Payment
	if err := database.Gorm.First(&payment, paymentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(PaymentErrors[utils.NotFound], err.Error()))
	}

	if err := database.Gorm.Delete(&payment, paymentID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(PaymentErrors[utils.Delete], err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]any{"id": paymentID, "message": "Payment deleted successfully"})
}

// POST /api/v1/payments/mercadopago/add //////////////////////////////////////////////////////
func AddMPPayment(c echo.Context) error {
	body := models.NewMPPayment{}
	c.Bind(&body)

	var order models.Order
	if err := database.Gorm.First(&order, body.OrderID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(OrderErrors[utils.NotFound], err.Error()))
	}

	var newPayment models.Payment
	newPayment.OrderID = body.OrderID
	newPayment.Status = StatusPending
	newPayment.Platform = PlatformMP

	if err := database.Gorm.Create(&newPayment).Error; err != nil {
		return c.JSON(http.StatusConflict, utils.SCTMake(PaymentErrors[utils.Create], err.Error()))
	}

	return c.JSON(http.StatusOK, newPayment)
}

// POST /api/v1/payments/mercadopago/:id/end //////////////////////////////////////////////////////
func EndMPPayment(c echo.Context) error {
	paymentID := c.Param("id")

	body := models.EndMPPayment{}
	c.Bind(&body)

	var oldPayment models.Payment
	if err := database.Gorm.First(&oldPayment, paymentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.SCTMake(PaymentErrors[utils.NotFound], err.Error()))
	}

	oldPayment.Url = null.StringFrom(fmt.Sprintf("https://www.mercadopago.com.ar/tools/receipt-view/%d", body.PaymentNumber))
	oldPayment.Paid = body.Paid
	oldPayment.Received = body.Received
	oldPayment.Status = body.Status

	if err := database.Gorm.Where("id = ?", paymentID).Save(&oldPayment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.SCTMake(PaymentErrors[utils.Update], err.Error()))
	}

	fmt.Printf("RECEIVED PAYMENT:\nDate: %s\n\n", time.Now().Format(time.RFC1123Z))

	return c.JSON(http.StatusOK, oldPayment)
}
