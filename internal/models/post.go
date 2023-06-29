package models

import (
	"Nix_trainee_practic/internal/http/response"
	"time"
)

type Post struct {
	UserID      int64
	ID          int64
	Title       string
	Body        string
	Comments    []response.CommentResponse
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (p Post) DomainToResponse() response.PostResponse {
	resp := response.PostResponse{
		ID:     p.ID,
		UserID: p.UserID,
		Title:  p.Title,
		Body:   p.Body,
	}
	if len(p.Comments) != 0 {
		resp.Comments = p.Comments
	}
	return resp
}
