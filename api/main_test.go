package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const server_address = "http://localhost:4000/"

func TestIndex(t *testing.T) {
	resp, err := http.Get(server_address)
	if err != nil {
		return
	}
	defer resp.Body.Close()	

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLoginRequestSuccess(t *testing.T) {
	requestAddress := server_address + "v1/auth/login"

	requestBody, _ := json.Marshal(map[string]string{
		"email": "admin@admin.admin",
		"password": "password",
	})
	requestBodyByte := bytes.NewBuffer(requestBody)

	resp, err := http.Post(requestAddress, "application/json", requestBodyByte)
	if err != nil {
		return
	}

	defer resp.Body.Close()	

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLoginRequestFailEmptyCredentials(t *testing.T) {
	requestAddress := server_address + "v1/auth/login"

	requestBody, _ := json.Marshal(map[string]string{
		"email": "",
		"password": "",
	})
	requestBodyByte := bytes.NewBuffer(requestBody)

	resp, err := http.Post(requestAddress, "application/json", requestBodyByte)
	if err != nil {
		return
	}

	defer resp.Body.Close()	

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestLoginRequestFailWrongCredentials(t *testing.T) {
	requestAddress := server_address + "v1/auth/login"

	requestBody, _ := json.Marshal(map[string]string{
		"email": "admin@adm.admin",
		"password": "password123",
	})
	requestBodyByte := bytes.NewBuffer(requestBody)

	resp, err := http.Post(requestAddress, "application/json", requestBodyByte)
	if err != nil {
		return
	}

	defer resp.Body.Close()	

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegisterRequestSuccess(t *testing.T) {
	requestAddress := server_address + "v1/auth/register"

	requestBody, _ := json.Marshal(map[string]string{
		"name": "John",
		"surname": "Doe",
		"phone": "00000000000",
		"email": "john@doe.me",
		"password": "password",
		"confirm_password": "password",
	})
	requestBodyByte := bytes.NewBuffer(requestBody)

	resp, err := http.Post(requestAddress, "application/json", requestBodyByte)
	if err != nil {
		return
	}

	defer resp.Body.Close()	

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRegisterRequestFailEmailTaken(t *testing.T) {
	requestAddress := server_address + "v1/auth/register"

	requestBody, _ := json.Marshal(map[string]string{
		"name": "John",
		"surname": "Doe",
		"phone": "00000000000",
		"email": "admin@admin.admin",
		"password": "password",
		"confirm_password": "password",
	})
	requestBodyByte := bytes.NewBuffer(requestBody)

	resp, err := http.Post(requestAddress, "application/json", requestBodyByte)
	if err != nil {
		return
	}

	defer resp.Body.Close()	

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegisterRequestFailEmptyCredentials(t *testing.T) {
	requestAddress := server_address + "v1/auth/register"

	requestBody, _ := json.Marshal(map[string]string{
		"name": "",
		"surname": "",
		"phone": "",
		"email": "",
		"password": "",
		"confirm_password": "",
	})
	requestBodyByte := bytes.NewBuffer(requestBody)

	resp, err := http.Post(requestAddress, "application/json", requestBodyByte)
	if err != nil {
		return
	}

	defer resp.Body.Close()	

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
