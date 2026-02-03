package authprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GroovingSpaces/auth-provider/dto"
	"github.com/go-resty/resty/v2"
)

var restyClient *resty.Client
var authHost string

func Init(host string) {
	restyClient = resty.NewWithClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	authHost = host
}

func VerifyToken(token string) (dto.VerifyTokenResponse, error) {

	var response dto.VerifyTokenResponse

	url := fmt.Sprintf("%s/api/v1/auth/verify-token", authHost)

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Post(url)

	// Set default values
	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = "POST"
	response.RequestAPICallResult.RequestBody = ""
	response.RequestAPICallResult.ResponseStatusCode = 0

	// Safely access response properties
	if resp != nil {
		if resp.Request != nil && resp.Request.Header != nil {
			reqHeaders, _ := json.Marshal(resp.Request.Header)
			response.RequestAPICallResult.RequestHeaders = string(reqHeaders)
			if resp.Request.RawRequest != nil {
				response.RequestAPICallResult.Method = resp.Request.RawRequest.Method
			}
		}
		if resp.Header() != nil {
			respHeaders, _ := json.Marshal(resp.Header())
			response.RequestAPICallResult.ResponseHeaders = string(respHeaders)
		}
		response.RequestAPICallResult.RequestLatency = resp.Time().String()
		response.RequestAPICallResult.ResponseBody = string(resp.Body())
		response.RequestAPICallResult.ResponseStatusCode = resp.StatusCode()
	}

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "excedeed") {
			return response, errors.New(dto.ErrTimeoutError)
		}
		return response, err
	}

	if resp == nil {
		return response, errors.New("response is nil")
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return response, err
	}

	if response.Status == "OK" {
		return response, nil
	}
	return response, fmt.Errorf("%v", response.ErrorCode)
}

func GetCurrentUser(token string) (dto.GetCurrentUserResponse, error) {
	var response dto.GetCurrentUserResponse

	url := fmt.Sprintf("%s/api/v1/auth/me", authHost)

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get(url)

	// Set default values
	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = "GET"
	response.RequestAPICallResult.RequestBody = ""
	response.RequestAPICallResult.ResponseStatusCode = 0

	// Safely access response properties
	if resp != nil {
		if resp.Request != nil && resp.Request.Header != nil {
			reqHeaders, _ := json.Marshal(resp.Request.Header)
			response.RequestAPICallResult.RequestHeaders = string(reqHeaders)
			if resp.Request.RawRequest != nil {
				response.RequestAPICallResult.Method = resp.Request.RawRequest.Method
			}
		}
		if resp.Header() != nil {
			respHeaders, _ := json.Marshal(resp.Header())
			response.RequestAPICallResult.ResponseHeaders = string(respHeaders)
		}
		response.RequestAPICallResult.RequestLatency = resp.Time().String()
		response.RequestAPICallResult.ResponseBody = string(resp.Body())
		response.RequestAPICallResult.ResponseStatusCode = resp.StatusCode()
	}

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "excedeed") {
			return response, errors.New(dto.ErrTimeoutError)
		}
		return response, err
	}

	if resp == nil {
		return response, errors.New("response is nil")
	}

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return response, err
	}

	if response.Status == "OK" {
		return response, nil
	}
	return response, fmt.Errorf("%v", response.ErrorCode)
}

func GetRoles(token string) (dto.GetRolesResponse, error) {
	var response dto.GetRolesResponse

	url := fmt.Sprintf("%s/api/v1/roles", authHost)

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Get(url)

	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = "GET"
	response.RequestAPICallResult.RequestBody = ""
	response.RequestAPICallResult.ResponseStatusCode = 0

	if resp != nil {
		if resp.Request != nil && resp.Request.Header != nil {
			reqHeaders, _ := json.Marshal(resp.Request.Header)
			response.RequestAPICallResult.RequestHeaders = string(reqHeaders)
			if resp.Request.RawRequest != nil {
				response.RequestAPICallResult.Method = resp.Request.RawRequest.Method
			}
		}
		if resp.Header() != nil {
			respHeaders, _ := json.Marshal(resp.Header())
			response.RequestAPICallResult.ResponseHeaders = string(respHeaders)
		}
		response.RequestAPICallResult.RequestLatency = resp.Time().String()
		response.RequestAPICallResult.ResponseBody = string(resp.Body())
		response.RequestAPICallResult.ResponseStatusCode = resp.StatusCode()
	}

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "excedeed") {
			return response, errors.New(dto.ErrTimeoutError)
		}
		return response, err
	}

	if resp == nil {
		return response, errors.New("response is nil")
	}

	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return response, err
	}

	if response.Status == "OK" {
		return response, nil
	}
	return response, fmt.Errorf("%v", response.ErrorCode)
}
