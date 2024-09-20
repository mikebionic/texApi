package request

type CreateWorker struct {
	Fullname  string `json:"fullname" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address" binding:"required"`
	AboutSelf string `json:"about_self" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Services  []int  `json:"services" binding:"required,dive,gt=0"`
}

type UpdateWorker struct {
	ID        int                   `json:"id" binding:"gt=0"`
	Fullname  string                `json:"fullname" binding:"required"`
	Phone     string                `json:"phone" binding:"required"`
	Address   string                `json:"address" binding:"required"`
	AboutSelf string                `json:"about_self" binding:"required"`
	Password  string                `json:"password" binding:"omitempty"`
	Services  []UpdateWorkerService `json:"services" binding:"required,dive,gt=0"`
}

type UpdateWorkerService struct {
	PrevID int `json:"prev_id" binding:"gt=0"`
	NextID int `json:"next_id" binding:"gt=0"`
}
