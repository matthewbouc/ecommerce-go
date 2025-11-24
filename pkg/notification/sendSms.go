package notification

import (
	"ecommerce/config"
	"errors"
	"fmt"
	"os"

	api "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/twilio/twilio-go"
)

type SmsClient interface {
	SendSms(phoneNumber string, message string) error
}

type smsClient struct {
	config config.AppConfig
}

// Twilio client
func (smsClient smsClient) SendSms(phoneNumber string, message string) error {
	if phoneNumber == "" {
		return errors.New("phone number is required to send sms")
	}

	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	// Make sure TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN exists in your environment
	client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	params.SetBody(message)
	params.SetFrom(smsClient.config.TwilioNumber)
	params.SetTo(phoneNumber)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if resp.Body != nil {
		fmt.Println(*resp.Body)
	} else {
		fmt.Println(resp.Body)
	}

	return nil
}

func NewSmsClient(config config.AppConfig) SmsClient {
	return &smsClient{
		config: config,
	}
}
