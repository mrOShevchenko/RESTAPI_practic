package response

type CommentResponse struct {
	ID     int64  `json:"id" example:"1"`
	PostID int64  `json:"post_id" example:"1"`
	Name   string `json:"name" example:"Bob"`
	Email  string `json:"email" example:"example@email.com"`
	Body   string `json:"body" example:"lorem ipsum"`
}
