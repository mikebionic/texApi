package chat

type APIHandler struct {
	repository *Repository
	hub        *Hub
	jwtSecret  []byte
}

func NewAPIHandler(repository *Repository, hub *Hub, jwtSecret []byte) *APIHandler {
	return &APIHandler{
		repository: repository,
		hub:        hub,
		jwtSecret:  jwtSecret,
	}
}

type ApiError struct {
	Message string
	Status  int
}

func (e *ApiError) Error() string {
	return e.Message
}
