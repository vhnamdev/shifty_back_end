package dto

type RegisterRequest struct {
	FullName string `json:"full_name" form:"full_name" validate:"required"`
	Email    string `json:"email"     form:"email"     validate:"required,email"`
	Password string `json:"password"  form:"password"  validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"    form:"email"    validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type SendOTP struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}
