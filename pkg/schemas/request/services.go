package request

type CreateService struct {
	Title    Translate `json:"title" binding:"required"`
	ParentID int       `json:"parent_id" binding:"omitempty,gt=0"`
}

type UpdateService struct {
	ID       int       `json:"id" binding:"gt=0"`
	Title    Translate `json:"title" binding:"required"`
	ParentID int       `json:"parent_id" binding:"omitempty,gt=0"`
}
