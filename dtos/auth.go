package dtos

type RegisterRequest struct {
	UserName string `json:"user_name" validate:"required,gte=6"`
	Password string `json:"password" validate:"required,gte=6"`
}

type RegisterResponse struct {
	Meta Meta    `json:"meta"`
	Data *Tokens `json:"data,omitempty"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	UserName string `json:"user_name" validate:"required,gte=6"`
	Password string `json:"password" validate:"required,gte=6"`
}

type LoginResponse struct {
	Meta Meta    `json:"meta"`
	Data *Tokens `json:"data,omitempty"`
}
