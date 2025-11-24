package provider

import (
	"ecommerce/config"
	"ecommerce/pkg/notification/sms"
	"errors"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSmsClient struct {
	config config.AppConfig
}

func NewTwilioSmsClient(config config.AppConfig) sms.SmsClient {
	return &TwilioSmsClient{
		config: config,
	}
}

func (smsClient *TwilioSmsClient) SendSms(phoneNumber string, message string) error {
	if phoneNumber == "" {
		return errors.New("phone number is required to send sms")
	}

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
