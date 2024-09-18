package response

type AboutUs struct {
	ID       int                `json:"id"`
	Text     TranslateWithoutID `json:"text"`
	IsActive bool               `json:"is_active"`
}
