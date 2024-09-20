package request

type CreateAboutUs struct {
	Text Translate `json:"text" binding:"required"`
}

type UpdateAboutUs struct {
	ID       int       `json:"id" binding:"gt=0"`
	Text     Translate `json:"text" binding:"required"`
	IsActive bool      `json:"is_active"`
}
