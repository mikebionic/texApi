package request

type LoginForm struct {
	Phone             string `json:"phone" binding:"required"`
	Password          string `json:"password" binding:"required"`
	NotificationToken string `json:"token" binding:"omitempty"`
}

type UserVerify struct {
	SmsID    string `json:"id" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type NotificationToken struct {
	Token string `json:"token" binding:"required"`
}

type ForgetPasswordForm struct {
	Phone string `json:"phone" binding:"required"`
}

type UserNewPassword struct {
	ID       string `json:"id" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	OTP      string `json:"otp" binding:"required"`
	Password string `json:"password" binding:"required"`
}
