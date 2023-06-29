package response

type PostResponse struct {
	ID       int64             `json:"id" example:"1"`
	UserID   int64             `json:"user_id" example:"1"`
	Title    string            `json:"title" example:"Lorem ipsum"`
	Body     string            `json:"body" example:"Lorem ipsum"`
	Comments []CommentResponse `json:"comments"`
}
