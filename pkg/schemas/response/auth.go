package response

type Login struct {
	ID           int           `json:"id"`
	Fullname     string        `json:"fullname"`
	Phone        string        `json:"phone"`
	Address      string        `json:"address"`
	Password     string        `json:"-"`
	Subscription *Subscription `json:"subscription"`
}
