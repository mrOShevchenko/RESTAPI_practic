package models

import (
	"Nix_trainee_practic/internal/http/response"
	"time"
)

type Comment struct {
	ID          int64
	PostID      int64
	Name        string
	Email       string
	Body        string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (c Comment) DomainToResponse() response.CommentResponse {
	return response.CommentResponse{
		ID:     c.ID,
		PostID: c.PostID,
		Name:   c.Name,
		Email:  c.Email,
		Body:   c.Body,
	}
}

func (c Comment) AllCommentsDomainToResponse(comments []Comment) []response.CommentResponse {
	var convertDomainCommentsToResponse []response.CommentResponse
	for _, com := range comments {
		convertDomainCommentsToResponse = append(convertDomainCommentsToResponse, com.DomainToResponse())
	}
	return convertDomainCommentsToResponse
}
