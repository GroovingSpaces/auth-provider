package dto

type CreateUserRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Name     string   `json:"name" validate:"required,min=3"`
	RoleIds  []string `json:"role_ids" validate:"required,min=1"`
}

type UpdateUserRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Name     string   `json:"name" validate:"required,min=3"`
	RoleIds  []string `json:"role_ids" validate:"required,min=1"`
	IsActive bool     `json:"is_active" validate:"required"`
}
