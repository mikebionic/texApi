package response

type Translate struct {
	ID *int    `json:"id"`
	TK *string `json:"tk"`
	RU *string `json:"ru"`
	EN *string `json:"en"`
}

type TranslateWithoutID struct {
	TK *string `json:"tk"`
	RU *string `json:"ru"`
	EN *string `json:"en"`
}
