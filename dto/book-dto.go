package dto

type BookUpdateDTO struct {
	ID          uint64 `json:"id" form:"id"`
	Title       string `json:"title" form:"title" binding:"omitempty,min=2"`
	Description string `json:"description" form:"description" binding:"omitempty,min=10"`
	UserID      uint64 `json:"user_id,omitempty" form:"user_id,omitempty"`
}

type BookCreateDTO struct {
	Title       string `json:"title" form:"title" binding:"min=2"`
	Description string `json:"description" form:"description" binding:"required,min=10"`
	UserID      uint64 `json:"user_id" form:"user_id"`
}
