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

func VerifyTokenWithMiddleware(token string, permissionName string) (dto.VerifyTokenData, error) {
	verifyTokenResponse, err := VerifyToken(token)
	if err != nil {
		return dto.VerifyTokenData{}, err
	}

	if verifyTokenResponse.Data.Valid != true {
		return dto.VerifyTokenData{}, fmt.Errorf("%v", dto.ErrInvalidToken)
	}

	for _, role := range verifyTokenResponse.Data.User.Roles {
		if !role.IsActive {
			return dto.VerifyTokenData{}, fmt.Errorf("%v", dto.ErrRoleInactive)
		}
		for _, permission := range role.Permissions {
			if permission.Slug == permissionName {
				if !permission.IsActive {
					return dto.VerifyTokenData{}, fmt.Errorf("%v", dto.ErrPermissionInactive)
				}
				return verifyTokenResponse.Data, nil
			}
		}
	}

	return dto.VerifyTokenData{}, fmt.Errorf("%v", dto.ErrRoleForbidden)
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

func CreateUser(token string, payload dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	var response dto.CreateUserResponse

	url := fmt.Sprintf("%s/api/v1/users", authHost)

	requestBody, _ := json.Marshal(payload)

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)

	// Set default values
	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = "POST"
	response.RequestAPICallResult.RequestBody = string(requestBody)
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
		} else {
			return response, err
		}
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
	} else {
		return response, fmt.Errorf("%v", response.ErrorCode)
	}
}

func UpdateUser(token string, id string, payload dto.UpdateUserRequest) (dto.UpdateUserResponse, error) {
	var response dto.UpdateUserResponse

	url := fmt.Sprintf("%s/api/v1/users/%s", authHost, id)

	requestBody, _ := json.Marshal(payload)

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Put(url)

	// Set default values
	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = "PUT"
	response.RequestAPICallResult.RequestBody = string(requestBody)
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
		} else {
			return response, err
		}
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
	} else {
		return response, fmt.Errorf("%v", response.ErrorCode)
	}
}

func DeleteUser(token string, id string) (dto.DeleteUserResponse, error) {
	var response dto.DeleteUserResponse

	url := fmt.Sprintf("%s/api/v1/users/%s", authHost, id)

	resp, err := restyClient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		Delete(url)

	response.RequestAPICallResult.RequestURL = url
	response.RequestAPICallResult.Method = "DELETE"
	response.RequestAPICallResult.RequestBody = ""
	response.RequestAPICallResult.ResponseStatusCode = 0

	if resp != nil {
		if resp.Request != nil && resp.Request.Header != nil {
			reqHeaders, _ := json.Marshal(resp.Request.Header)
			response.RequestAPICallResult.RequestHeaders = string(reqHeaders)
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
		} else {
			return response, err
		}
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
	} else {
		return response, fmt.Errorf("%v", response.ErrorCode)
	}
}
