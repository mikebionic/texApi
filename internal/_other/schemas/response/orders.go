package response

type Orders struct {
	Total  int     `json:"total"`
	Orders []Order `json:"orders"`
}

type Order struct {
	ID           int            `json:"id"`
	OrderNumber  string         `json:"order_number"`
	FilePaths    *[]*string     `json:"file_paths"`
	UserID       int            `json:"user_id"`
	WorkerID     *int           `json:"worker_id"`
	Address      string         `json:"address"`
	Date         string         `json:"date"`
	Time         string         `json:"time"`
	TimeDuration *int           `json:"time_duration"`
	Status       Translate      `json:"status"`
	Description  *string        `json:"description"`
	Services     []OrderService `json:"services"`
	ReadByAdmin  bool           `json:"read_by_admin"`
}

type OrderService struct {
	ID    *int    `json:"id"`
	Image *string `json:"image"`
	TK    *string `json:"tk"`
	RU    *string `json:"ru"`
	EN    *string `json:"en"`
}

type OrderByWorker struct {
	ID           int          `json:"id"`
	OrderNumber  string       `json:"order_number"`
	FilePaths    *[]*string   `json:"file_paths"`
	User         UserByWorker `json:"user"`
	WorkerID     *int         `json:"worker_id"`
	Address      string       `json:"address"`
	Date         string       `json:"date"`
	Time         string       `json:"time"`
	TimeDuration *int         `json:"time_duration"`
	Status       Translate    `json:"status"`
	Description  *string      `json:"description"`
	Services     []Translate  `json:"services"`
}

type OrderByUser struct {
	ID           int           `json:"id"`
	OrderNumber  string        `json:"order_number"`
	FilePaths    *[]*string    `json:"file_paths"`
	Worker       *WorkerByUser `json:"worker"`
	Address      string        `json:"address"`
	Date         string        `json:"date"`
	Time         string        `json:"time"`
	TimeDuration *int          `json:"time_duration"`
	SecretWord   string        `json:"secret_word"`
	Status       Translate     `json:"status"`
	Description  *string       `json:"description"`
	Services     []Translate   `json:"services"`
}

type UserByWorker struct {
	Fullname string `json:"fullname"`
	Phone    string `json:"phone"`
}

type WorkerByUser struct {
	Fullname  string  `json:"fullname"`
	Phone     string  `json:"phone"`
	Photo     *string `json:"photo"`
	AboutSelf string  `json:"about_self"`
}
