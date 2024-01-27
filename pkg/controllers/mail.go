package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type EmailAddress struct {
	Email string
	Name	string
}

var emailSender = EmailAddress{
	Email: "geras4323info@gmail.com",
	Name: "Geraldo",
}

func GetEmail(c echo.Context) error {
	mailjetClient := mailjet.NewMailjetClient("b66eee0c45cd1d455f2dd52864d19b60", "bd1baf7dae632355cb54b56b12339984")

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: emailSender.Email,
				Name: emailSender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "german432yo@gmail.com",
					Name: "Nombrecito",
				},
			},
			Subject: "EmailTest",
			TemplateLanguage: true,
			TemplateID: 4839487,
			Variables: map[string]interface{}{"name": "Geras", "items":[]map[string]interface{}{{"id": 1, "article": "Bike"}, {"id": 2, "article": "Couch"}}},
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}