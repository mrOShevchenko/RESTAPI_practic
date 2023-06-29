package requests

type PostRequest struct {
	Title string `json:"title" example:"Lorem ipsum" validate:"required"`
	Body  string `json:"body" example:"Lorem ipsum" validate:"required"`
}
