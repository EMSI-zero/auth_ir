package sms_provider

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/supabase/auth/internal/conf"
)

const(
	defaultKavenegarAPIBase = "https://api.kavenegar.com/"
)

type KavenegarProvider struct {
	Config  *conf.KavenegarProviderConfiguration
	APIPath string
	APIVersion string
}

func NewKavenegarProvider(config conf.KavenegarProviderConfiguration) (SmsProvider, error) {
	if err:= config.Validate() ;err != nil {
		return err
	}
	return nil, nil
}

func (k *KavenegarProvider) SendMessage(phone, message, channel, otp string) (string, error) {
	switch channel {
	case SMSProvider:
		return k.SendSMS(phone, message)
	default:
		return "", fmt.Errorf("channel type %q is not supported for TextLocal", channel)
	}
}

func (k *KavenegarProvider) SendSMS(phone, message string) (string, error) {
	body := url.Values{
		"sender":   {k.Config.Sender},
		"receptor": {phone},
		"message":  {message},
	}

	client := &http.Client{Timeout: defaultTimeout}
	
}
