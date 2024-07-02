package sms_provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/emsi-zero/auth_ir/internal/conf"
	"github.com/emsi-zero/auth_ir/internal/utilities"
)

const (
	defaultKavenegarAPIBase = "https://api.kavenegar.com"
)

type KavenegarProvider struct {
	Config     *conf.KavenegarProviderConfiguration
	APIPath    string
	APIVersion string
}

func NewKavenegarProvider(config conf.KavenegarProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	apiPath := defaultKavenegarAPIBase + "/v1/" + config.ApiKey + "/verify/lookup.json"
	return &KavenegarProvider{
		Config:     &config,
		APIPath:    apiPath,
		APIVersion: "v1",
	}, nil
}

func (k *KavenegarProvider) SendMessage(phone, message, channel, otp string) (string, error) {
	switch channel {
	case SMSProvider:
		return k.SendSMS(phone, message)
	default:
		return "", fmt.Errorf("channel type %q is not supported for TextLocal", channel)
	}
}

type KavenegarResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []*struct {
		MessageID  int64  `json:"messageid"`
		Message    string `json:"message"`
		Status     int    `json:"status"`
		StatusText string `json:"statustext"`
		Sender     string `json:"sender"`
		Receptor   string `json:"receptor"`
		Date       int64  `json:"date"`
		Cost       int64  `json:"cost"`
	} `json:"entries"`
}

func (k *KavenegarProvider) SendSMS(phone, message string) (string, error) {
	body := url.Values{
		"sender":   {k.Config.Sender},
		"receptor": {phone},
		"token":    {message},
		"template": {k.Config.OTPTemplateKey},
	}

	client := &http.Client{Timeout: defaultTimeout}

	request, err := http.NewRequest("POST", k.APIPath, strings.NewReader(body.Encode()))
	request.PostForm = body
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}

	resbody, err := client.Do(request)
	if err != nil {
		return "", err
	}

	resp := &KavenegarResponse{}
	derr := json.NewDecoder(resbody.Body).Decode(resp)
	if derr != nil {
		return "", derr
	}
	defer utilities.SafeClose(resbody.Body)

	if resp.Return.Status != 200 {
		return "", fmt.Errorf("twilio error: %v %v for message %s", resp.Return.Message, resp.Return.Status, resp.Entries[0].Message)
	}

	return strconv.Itoa(int(resp.Entries[0].MessageID)), nil

}
