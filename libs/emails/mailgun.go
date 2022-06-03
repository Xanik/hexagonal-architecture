package emails

import (
	"bytes"
	"context"
	"fmt"
	"study/config"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

func SendEmail(subject, body, recipient string) {
	MAILGUN_DOMAIN := config.Env.MailGunDomain
	API_KEY := config.Env.MailGunApiKey
	SENDER := "test@gmail.com"

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(MAILGUN_DOMAIN, API_KEY)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(SENDER, subject, "", recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	// Send the message	with a 10 second timeout

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

type data = map[string]interface{}

type Message struct {
	subject       string
	message       string
	isHTML        bool
	hasAttachment bool
	attachment    interface{}
}

func SendTemplate(subject, to, name string, code int) {
	message := new(Message)
	body := new(bytes.Buffer)

	err := email.Execute(body, data{
		"Name": name,
		"Code": code,
	})

	fmt.Println(err)

	message.subject = subject
	message.isHTML = true
	message.message = body.String()

	SendEmail(subject, body.String(), to)
}
