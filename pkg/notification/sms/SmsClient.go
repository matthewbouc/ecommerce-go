package sms

type SmsClient interface {
	SendSms(phoneNumber string, message string) error
}
