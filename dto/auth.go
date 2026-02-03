package dto

import (
	"time"
)

type AuthResponse struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	ErrorCode string   `json:"error_code"`
	TrxID     string   `json:"trx_id"`
	Data      AuthData `json:"data"`
	Message   string   `json:"message"`

	RequestAPICallResult RequestAPICallResult `json:"-"`
}

func (s *AuthResponse) GetAPICall() RequestAPICallResult {
	return s.RequestAPICallResult
}

type AuthData struct {
	Valid      bool         `json:"valid"`
	User       AuthUserData `json:"user"`
	Claims     AuthClaims   `json:"claims"`
	Token      string       `json:"token"`
	EmployeeID string       `json:"employee_id"`
}

type AuthUserData struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	IsActive    bool      `json:"is_active"`
	LastLoginAt time.Time `json:"last_login_at"`
	LastLoginIP string    `json:"last_login_ip"`
	Roles       []struct {
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
	} `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthClaims struct {
	Email    string `json:"email"`
	Exp      int    `json:"exp"`
	Iat      int    `json:"iat"`
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
