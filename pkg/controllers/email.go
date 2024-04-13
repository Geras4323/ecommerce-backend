package controllers

import (
	"net/http"

	"github.com/geras4323/ecommerce-backend/pkg/cloud"
	"github.com/labstack/echo/v4"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type EmailData struct {
	Name  string
	Items []cloud.EmailProduct
}

func GetEmail(c echo.Context) error {
	var body EmailData
	c.Bind(&body)

	items := make([]map[string]interface{}, len(body.Items))
	for i, item := range body.Items {
		items[i] = map[string]interface{}{
			"article":  item.Article,
			"quantity": item.Quantity,
			"price":    item.Price,
		}
	}

	variables := map[string]interface{}{
		"name":  body.Name,
		"items": items,
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cloud.DefaultSender.Email,
				Name:  cloud.DefaultSender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "german432yo@gmail.com",
				},
			},
			Subject:          "Confirmación de pedido",
			TemplateLanguage: true,
			TemplateID:       4839487,
			Variables:        variables,
		},
	}

	res, err := cloud.SendMail(messagesInfo)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
