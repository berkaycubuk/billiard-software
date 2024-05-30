package vipps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/shopspring/decimal"
)

const (
	PAYMENT_COMPLETED string = "AUTHORIZED"
	PAYMENT_CANCELED	 = "TERMINATED"
	PAYMENT_WAITING	         = "CREATED"
)

type VippsInitiateResponse struct {
	Token			string
	CheckoutFrontendUrl	string
	PollingUrl		string
}

type VippsSessionPaymentDetails struct {
	State			string		`json:"state"`
}

type VippsSessionResponse struct {
	SessionState		string		`json:"sessionState"`
	PaymentDetails		VippsSessionPaymentDetails `json:"paymentDetails"`
}

type VippsPaymentResponse struct {
	State		string		`json:"state"`
}

type VippsAccessTokenResponse struct {
	AccessToken	string		`json:"access_token"`
}

func AccessToken() (string, error) {
	endpoint := os.Getenv("VIPPS_ENDPOINT")
	billiardSoftwareVersion := os.Getenv("VERSION")
	merchantSerialNumber := os.Getenv("VIPPS_MERCHANT_SERIAL_NUMBER")
	clientId := os.Getenv("VIPPS_CLIENT_ID")
	clientSecret := os.Getenv("VIPPS_CLIENT_SECRET")
	ocpApim := os.Getenv("VIPPS_OCP_APIM_SUBSCRIPTION_KEY")

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/accesstoken/get", endpoint), nil)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Vipps-System-Name", "Billiard Software")
	request.Header.Set("Vipps-System-Version", billiardSoftwareVersion)
	request.Header.Set("Vipps-System-Plugin-Name", "Billiard Software Vipps")
	request.Header.Set("Vipps-System-Plugin-Version", "1.0.0")
	request.Header.Set("Merchant-Serial-Number", merchantSerialNumber)
	request.Header.Set("client_id", clientId)
	request.Header.Set("client_secret", clientSecret)
	request.Header.Set("Ocp-Apim-Subscription-Key", ocpApim)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New(string(responseBody))
	}

	var vippsResponse VippsAccessTokenResponse
	err = json.Unmarshal([]byte(responseBody), &vippsResponse)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	return fmt.Sprintf("bearer %v", vippsResponse.AccessToken), nil
}

func Initiate(orderID uint, price decimal.Decimal, ref string, description string) (*VippsInitiateResponse, error) {
	endpoint := os.Getenv("VIPPS_ENDPOINT")
	callbackUrl := os.Getenv("VIPPS_CALLBACK_URL")
	returnUrl := fmt.Sprintf("%v/%v", os.Getenv("VIPPS_RETURN_URL"), orderID)
	termsAndConditionsUrl := os.Getenv("VIPPS_TERMS_AND_CONDITIONS_URL")
	billiardSoftwareVersion := os.Getenv("VERSION")
	merchantSerialNumber := os.Getenv("VIPPS_MERCHANT_SERIAL_NUMBER")
	clientId := os.Getenv("VIPPS_CLIENT_ID")
	clientSecret := os.Getenv("VIPPS_CLIENT_SECRET")
	ocpApim := os.Getenv("VIPPS_OCP_APIM_SUBSCRIPTION_KEY")

	// convert price
	priceMultiplied := price.Mul(decimal.NewFromFloat(100))
	convertedPrice := priceMultiplied.IntPart()

	bodyData := []byte(fmt.Sprintf(`
	{
		"merchantInfo": {
			"callbackUrl": "%v",
			"returnUrl": "%v",
			"callbackAuthorizationToken": "--",
			"termsAndConditionsUrl": "%v"
		},
		"transaction": {
			"amount": {
				"value": %v,
				"currency": "NOK"
			},
			"reference": "%v",
			"paymentDescription": "%v"
		},
		"configuration": {
			"customerInteraction": "CUSTOMER_NOT_PRESENT",
			"elements": "PaymentOnly"
		}
	}
	`, callbackUrl, returnUrl, termsAndConditionsUrl, convertedPrice, ref, description))

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/checkout/v3/session", endpoint), bytes.NewBuffer(bodyData))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Vipps-System-Name", "Billiard Software")
	request.Header.Set("Vipps-System-Version", billiardSoftwareVersion)
	request.Header.Set("Vipps-System-Plugin-Name", "Billiard Software Vipps")
	request.Header.Set("Vipps-System-Plugin-Version", "1.0.0")
	request.Header.Set("Merchant-Serial-Number", merchantSerialNumber)
	request.Header.Set("client_id", clientId)
	request.Header.Set("client_secret", clientSecret)
	request.Header.Set("Ocp-Apim-Subscription-Key", ocpApim)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(string(responseBody))
	}

	var vippsResponse VippsInitiateResponse
	err = json.Unmarshal([]byte(responseBody), &vippsResponse)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return &vippsResponse, nil
}

func Status(orderID uint, ref string) (*VippsSessionResponse, error) {
	endpoint := os.Getenv("VIPPS_ENDPOINT")
	billiardSoftwareVersion := os.Getenv("VERSION")
	merchantSerialNumber := os.Getenv("VIPPS_MERCHANT_SERIAL_NUMBER")
	clientId := os.Getenv("VIPPS_CLIENT_ID")
	clientSecret := os.Getenv("VIPPS_CLIENT_SECRET")
	ocpApim := os.Getenv("VIPPS_OCP_APIM_SUBSCRIPTION_KEY")

	request, err := http.NewRequest("GET", fmt.Sprintf("%v/checkout/v3/session/%v", endpoint, ref), nil)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Vipps-System-Name", "Billiard Software")
	request.Header.Set("Vipps-System-Version", billiardSoftwareVersion)
	request.Header.Set("Vipps-System-Plugin-Name", "Billiard Software Vipps")
	request.Header.Set("Vipps-System-Plugin-Version", "1.0.0")
	request.Header.Set("Merchant-Serial-Number", merchantSerialNumber)
	request.Header.Set("client_id", clientId)
	request.Header.Set("client_secret", clientSecret)
	request.Header.Set("Ocp-Apim-Subscription-Key", ocpApim)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(string(responseBody))
	}

	var vippsResponse VippsSessionResponse
	err = json.Unmarshal([]byte(responseBody), &vippsResponse)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return &vippsResponse, nil
}

func GetPayment(ref string) (*VippsPaymentResponse, error) {
	endpoint := os.Getenv("VIPPS_ENDPOINT")
	merchantSerialNumber := os.Getenv("VIPPS_MERCHANT_SERIAL_NUMBER")
	clientId := os.Getenv("VIPPS_CLIENT_ID")
	clientSecret := os.Getenv("VIPPS_CLIENT_SECRET")
	ocpApim := os.Getenv("VIPPS_OCP_APIM_SUBSCRIPTION_KEY")

	accessToken, err := AccessToken()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%v/epayment/v1/payments/%v", endpoint, ref), nil)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", accessToken)
	request.Header.Set("Merchant-Serial-Number", merchantSerialNumber)
	request.Header.Set("client_id", clientId)
	request.Header.Set("client_secret", clientSecret)
	request.Header.Set("Ocp-Apim-Subscription-Key", ocpApim)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(string(responseBody))
	}

	var vippsResponse VippsPaymentResponse
	err = json.Unmarshal([]byte(responseBody), &vippsResponse)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return &vippsResponse, nil
}

func Capture(ref string, price decimal.Decimal) error {
	endpoint := os.Getenv("VIPPS_ENDPOINT")
	merchantSerialNumber := os.Getenv("VIPPS_MERCHANT_SERIAL_NUMBER")
	ocpApim := os.Getenv("VIPPS_OCP_APIM_SUBSCRIPTION_KEY")

	// convert price
	priceMultiplied := price.Mul(decimal.NewFromFloat(100))
	convertedPrice := priceMultiplied.IntPart()

	// access token stuff
	accessToken, err := AccessToken()
	if err != nil {
		return err
	}

	bodyData := []byte(fmt.Sprintf(`
	{
		"modificationAmount": {
			"currency": "NOK",
			"value": %v
		}
	}
	`, convertedPrice))

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/epayment/v1/payments/%v/capture", endpoint, ref), bytes.NewBuffer(bodyData))
	request.Header.Set("Authorization", accessToken)
	request.Header.Set("Merchant-Serial-Number", merchantSerialNumber)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Ocp-Apim-Subscription-Key", ocpApim)
	request.Header.Set("Idempotency-Key", utils.RandomString(10))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New(string(responseBody))
	}

	var vippsResponse VippsSessionResponse
	err = json.Unmarshal([]byte(responseBody), &vippsResponse)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}
