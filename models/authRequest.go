package models

type AuthRequest struct {
	Password      string				`json:"password" validate:"required,min=6"`
	Email         string				`json:"email" validate:"required,email"`
}