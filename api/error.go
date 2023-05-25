package api

type errorResponse struct {
	Message string `json:"message"`
	Debug   string `json:"debug"`
}

func newBadRequest(err error) *errorResponse {
	return &errorResponse{
		Message: "リクエストに何らかの間違いがあります。",
		Debug:   err.Error(),
	}
}

func newInternalServerError(err error) *errorResponse {
	return &errorResponse{
		Message: "サーバで何らかのエラーが発生しました。もう一度お試しください。",
		Debug:   err.Error(),
	}
}
