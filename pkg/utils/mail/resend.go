package main

import (
	"fmt"

	"github.com/resendlabs/resend-go"
)

func main() {
    apiKey := "re_cp1shcMv_3DQGo4aSg97Bh9CCFaifVSS2"

    client := resend.NewClient(apiKey)

		link := "https://www.google.com"

    params := &resend.SendEmailRequest{
        From:   	"onboarding@resend.dev",
        To:     	[]string{"german432yo@gmail.com"},
        Html: 		fmt.Sprintf(`
					<div>
						<h1 style='color: #74c27e'>A password recovery has been requested.</h1>
						<h2 style='color: black'>Please click this button to reset your password:</h2>
						<a href='%s' style='width: 200px; text-align: center; vertical-align: center; font-size: 16px; display: grid; align-items: center; text-decoration: none; border-radius: 8px; color: #74c27e; background-color: white; border: 1px solid #74c27e'>
							Recover password
						</a>
						<p style='font-size: 16px; color: black'>If the button doesn&apos;t work, try this link instead:</p>
						<p>%s</p>
					</div>
				`, link, link),
        Subject:	"Hello from Golang",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(sent.Id)
}