package request

type ContentRequest struct {
	LangID        int    `json:"lang_id" binding:"required"`
	ContentTypeID int    `json:"content_type_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	VideoURL      string `json:"video_url"`
	Step          int    `json:"step"`
}
