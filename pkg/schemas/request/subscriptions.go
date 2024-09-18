package request

type CreateSubscription struct {
	Title   Translate `json:"title" binding:"required"`
	Desc    Translate `json:"description" binding:"required"`
	StartAt string    `json:"start_at" binding:"required"`
	EndAt   string    `json:"end_at" binding:"required"`
	Days    int       `json:"days" binding:"gt=0"`
	Count   int       `json:"count" binding:"gt=0"`
	Price   float32   `json:"price" binding:"gt=0"`
}

type UpdateSubscription struct {
	ID      int       `json:"id" binding:"gt=0"`
	Title   Translate `json:"title" binding:"required"`
	Desc    Translate `json:"description" binding:"required"`
	StartAt string    `json:"start_at" binding:"required"`
	EndAt   string    `json:"end_at" binding:"required"`
	Days    int       `json:"days" binding:"gt=0"`
	Count   int       `json:"count" binding:"gt=0"`
	Price   float32   `json:"price" binding:"gt=0"`
}
