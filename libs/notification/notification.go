package notification

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"study/config"
	"time"

	"github.com/maddevsio/fcm"

	"github.com/mailgun/mailgun-go/v3"
)

func mailGun(sender, subject, body, recipient string) {
	yourDomain := config.Env.MailGunDomain
	privateAPIKey := config.Env.MailGunApiKey

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

//SendMail takes in sender, subject, body, recipient and sends an smtp mail to the described destination
func SendMail(sender, subject, body, recipient string) {
	mailGun(sender, subject, body, recipient)
}

// SendText takes in the recipients number and sends a text
func SendText(recipient string, message int, sender string) {
	email := "test@gmail.com"
	ebulk := ""
	response := fmt.Sprintf("test.com/sendsms?username=%s&apikey=%s&sender=%s&messagetext=%d&flash=0&recipients=%s", email, ebulk, sender, message, recipient)
	resp, err := http.Get(response)
	fmt.Printf(response)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
	}
}

// SendPushNotification takes in the recipients number and sends a text
func SendPushNotification(title string, body string, token string, id string, model interface{}, types string) {
	data := map[string]interface{}{
		"title":        title,
		"body":         body,
		"type":         types,
		"request_id":   id,
		"request_body": model,
	}
	c := fcm.NewFCM("")
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  []string{token},
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title:       title,
			Body:        body,
			ClickAction: "OPEN_ACTIVITY_1",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Status Code   :", response.StatusCode)
	fmt.Println("Success       :", response.Success)
	fmt.Println("Fail          :", response.Fail)
	fmt.Println("Canonical_ids :", response.CanonicalIDs)
	fmt.Println("Topic MsgId   :", response.MsgID)

}

func SendTwilloNotification(recipient string, message int) {

	// Set initial variables
	accountSid := ""
	authToken := ""
	urlStr := ""

	fmt.Println(recipient, message)
	// Build out the data for our message
	v := url.Values{}
	v.Set("To", recipient)
	v.Set("From", "")
	v.Set("Body", strconv.Itoa(message))
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
