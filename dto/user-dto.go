package dto

type UserUpdateDTO struct {
	ID       uint64 `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"min=3"`
	Email    string `json:"email,omitempty" form:"email,omitempty" binding:"omitempty,email"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"omitempty,min=6"`
}
