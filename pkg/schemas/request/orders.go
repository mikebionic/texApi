package request

type CreateOrder struct {
	Description string `json:"description" binding:"omitempty"`
	Address     string `json:"address" binding:"required"`
	Services    []int  `json:"services" binding:"required,dive,gt=0"`
	Date        string `json:"date" binding:"datetime=2006-01-02"`
	Time        string `json:"time" binding:"required"`
	SecretWord  string `json:"secret_word" binding:"required"`
}

type UpdateOrder struct {
	ID       int `json:"id" binding:"gt=0"`
	WorkerID int `json:"worker_id" binding:"gt=0"`
	StatusID int `json:"status_id" binding:"oneof=1 2 3 4 5"`
}

type UpdateOrderTimeDuration struct {
	ID           int `json:"id" binding:"gt=0"`
	TimeDuration int `json:"time_duration" binding:"gt=0"`
}

type UpdateOrderStatusStart struct {
	ID int `json:"id" binding:"gt=0"`
}

type OrderByWorker struct {
	StatusID string `json:"status_id" binding:"oneof=1 2 3 4 5"`
}
