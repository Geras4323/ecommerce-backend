package cloud

import (
	"github.com/geras4323/ecommerce-backend/pkg/utils"
	"github.com/mailjet/mailjet-apiv3-go"
)

type SenderAddress struct {
	Email string
	Name  string
}

type EmailProduct struct {
	Article  string
	Quantity uint
	Price    float64
}

var DefaultSender = SenderAddress{
	Email: "geras4323info@gmail.com",
	Name:  "Swiftharbour",
}

var Mjt *mailjet.Client

func ConnectMailjet() {
	mailjetClient := mailjet.NewMailjetClient(utils.GetEnvVar("MAILJET_PUBLIC_KEY"), utils.GetEnvVar("MAILJET_PRIVATE_KEY"))

	Mjt = mailjetClient
}

func SendMail(messagesInfo []mailjet.InfoMessagesV31) (*mailjet.ResultsV31, error) {
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := Mjt.SendMailV31(&messages)
	if err != nil {
		return nil, err
	}
	return res, nil
}
