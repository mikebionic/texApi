package request

type RegisterOrder struct {
	ApiClient          string `json:"api_client"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Amount             int    `json:"amount"`
	Description        string `json:"description"`
	SessionTimeoutSecs int    `json:"sessionTimeoutSecs"`
}

type SubmitCard struct {
	MDOrder  string `json:"MDORDER"`
	Expiry   string `json:"$EXPIRY"`
	Pan      string `json:"$PAN"`
	Text     string `json:"TEXT"`
	Cvc      string `json:"$CVC"`
	Language string `json:"language"`
}

type SubmitOTP struct {
	MDOrder      string `json:"MDORDER"`
	PasswordEdit string `json:"passwordEdit"`
}

type CheckStatus struct {
	OrderID string `json:"orderId"`
}

type ResendPassword struct {
	OrderID string `json:"order_id"`
}
