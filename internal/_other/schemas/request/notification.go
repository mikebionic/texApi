package request

type NotificationForm struct {
	To           string       `json:"to"`
	Notification Notification `json:"notification"`
}

type Notification struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}
