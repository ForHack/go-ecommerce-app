package notification

import (
	"encoding/json"
	"fmt"
	"go-ecommerce-app/configs"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type NotificationClient interface {
	SendSMS(phone string, message string) error
}

type notificationClient struct {
	config configs.AppConfig
}

// Twilio
func (c notificationClient) SendSMS(phone string, message string) error {
	accountSid := c.config.TwilioAccountSid
	authToken := c.config.TwilioAuthToken

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(c.config.TwilioFromPhoneNumber)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

	return nil
}

func NewNotificationClient(config configs.AppConfig) NotificationClient {
	return &notificationClient{
		config: config,
	}
}
