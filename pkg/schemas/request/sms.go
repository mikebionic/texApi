package request

type SMS struct {
	To       string `json:"to"`
	Text     string `json:"text"`
	TextType string `json:"text_type"`
	ApiKey   string `json:"api_key"`
}
