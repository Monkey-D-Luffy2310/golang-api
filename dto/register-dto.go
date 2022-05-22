package dto

type RegisterDTO struct {
	Name     string `json:"name" form:"name" binding:"min=3"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"min=6"`
}
