package response

type RegisterOrder struct {
	OrderID string `json:"orderId"`
}

type Submit struct {
	Message string `json:"message"`
}

type CheckStatus struct {
	ErrorMessage string `json:"errorMessage"`
	OrderStatus  int    `json:"orderStatus"`
}
