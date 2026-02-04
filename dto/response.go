package dto

import (
	"time"
)

type VerifyTokenResponse struct {
	Code      int             `json:"code"`
	Status    string          `json:"status"`
	ErrorCode string          `json:"error_code"`
	TrxID     string          `json:"trx_id"`
	Data      VerifyTokenData `json:"data"`

	// APICall Result
	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *VerifyTokenResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

type VerifyTokenData struct {
	Valid  bool         `json:"valid"`
	User   AuthUserData `json:"user"`
	Claims AuthClaims   `json:"claims"`
	Token  string       `json:"token"`
}

type UserData struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	IsActive    bool      `json:"is_active"`
	LastLoginAt time.Time `json:"last_login_at"`
	LastLoginIP string    `json:"last_login_ip"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClaimsResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Exp      int    `json:"exp"`
}

type GetCurrentUserResponse struct {
	Code      int                `json:"code"`
	Status    string             `json:"status"`
	ErrorCode string             `json:"error_code"`
	TrxID     string             `json:"trx_id"`
	Data      GetCurrentUserData `json:"data"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *GetCurrentUserResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

type GetCurrentUserData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CreateUserResponse struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	ErrorCode string   `json:"error_code"`
	TrxID     string   `json:"trx_id"`
	Data      UserData `json:"data"`
	Message   string   `json:"message"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *CreateUserResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

type GetUsersResponse struct {
	Code      int          `json:"code"`
	Status    string       `json:"status"`
	ErrorCode string       `json:"error_code"`
	TrxID     string       `json:"trx_id"`
	Data      GetUsersData `json:"data"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

type GetUsersData struct {
	Users      []UserData `json:"users"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

func (s *GetUsersResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

type GetUserResponse struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	ErrorCode string   `json:"error_code"`
	TrxID     string   `json:"trx_id"`
	Data      UserData `json:"data"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *GetUserResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

type UpdateUserResponse struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	ErrorCode string   `json:"error_code"`
	TrxID     string   `json:"trx_id"`
	Data      UserData `json:"data"`
	Message   string   `json:"message"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *UpdateUserResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

// GetRolesResponse for listing roles from auth service.
// API returns data as array directly: "data": [{...}, {...}]
type GetRolesResponse struct {
	Code                 int                  `json:"code"`
	Status               string               `json:"status"`
	ErrorCode            string               `json:"error_code"`
	TrxID                string               `json:"trx_id"`
	Data                 []RoleData           `json:"data"`
	RequestAPICallResult RequestAPICallResult `json:"-"`
}

type RoleData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	Permissions []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Module      string `json:"module"`
		Action      string `json:"action"`
		IsActive    bool   `json:"is_active"`
	} `json:"permissions"`
}

type DeleteUserResponse struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	ErrorCode string `json:"error_code"`
	TrxID     string `json:"trx_id"`
	Data      any    `json:"data"`
	Message   string `json:"message"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *DeleteUserResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}
