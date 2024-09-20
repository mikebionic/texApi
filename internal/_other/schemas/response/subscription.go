package response

type Subscription struct {
	ID          int                `json:"id"`
	Title       TranslateWithoutID `json:"title"`
	Description TranslateWithoutID `json:"description"`
	StartAt     string             `json:"start_at"`
	EndAt       string             `json:"end_at"`
	Days        int                `json:"days"`
	Count       int                `json:"count"`
	Price       float32            `json:"price"`
}
