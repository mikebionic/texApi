package response

type Worker struct {
	ID        int          `json:"id"`
	Fullname  string       `json:"fullname"`
	Phone     string       `json:"phone"`
	Address   string       `json:"address"`
	Photo     *string      `json:"photo"`
	AboutSelf string       `json:"about_self"`
	Services  *[]Translate `json:"services"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}
