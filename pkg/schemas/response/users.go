package response

type User struct {
	ID             int    `json:"id"`
	Fullname       string `json:"fullname"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	IsVerified     bool   `json:"is_verified"`
	SubscriptionID *int   `json:"subscription_id"`
}

type UserExist struct {
	Phone      string `json:"phone"`
	IsVerified bool   `json:"is_verified"`
}
