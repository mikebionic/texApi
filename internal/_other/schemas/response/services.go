package response

type Service struct {
	ID       int                `json:"id"`
	Title    TranslateWithoutID `json:"title"`
	Image    *string            `json:"image"`
	ParentID *int               `json:"parent_id"`
	Children []*Service         `json:"children"`
}

type ServiceList struct {
	ID    int                `json:"id"`
	Title TranslateWithoutID `json:"title"`
}
