package requests

type CommentRequest struct {
	Body string `json:"body" example:"lorem ipsum" validate:"required"`
}
