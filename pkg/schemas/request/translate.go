package request

type Translate struct {
	TK string `json:"tk" binding:"required"`
	RU string `json:"ru" binding:"required"`
	EN string `json:"en" binding:"required"`
}
