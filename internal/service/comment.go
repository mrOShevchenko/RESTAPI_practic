package service

import (
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/repository"
	"fmt"
	"github.com/golang-jwt/jwt"
)

//go:generate mockery --dir . --name CommentService --output ./mocks
type CommentService interface {
	SaveComment(commentRequest requests.CommentRequest, postID int64, token *jwt.Token) (models.Comment, error)
	GetComment(id int64) (models.Comment, error)
	UpdateComment(commentRequest requests.CommentRequest, id int64) (models.Comment, error)
	DeleteComment(id int64) error
	GetCommentsByPostID(postID int64, offset int) ([]models.Comment, error)
}

type commentService struct {
	repo repository.CommentRepo
	us   UserService
	ps   PostService
}

func NewComment(repo repository.CommentRepo, us UserService, ps PostService) CommentService {
	return commentService{
		repo: repo,
		us:   us,
		ps:   ps,
	}
}

func (s commentService) SaveComment(commentRequest requests.CommentRequest, postID int64, token *jwt.Token) (models.Comment, error) {
	claims := token.Claims.(*JWTClaim)
	_, err := s.ps.GetPost(postID)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service error save comment: %w", err)
	}
	user, err := s.us.FindByID(claims.ID)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service error save comment: %w", err)
	}
	domainComment := models.Comment{
		PostID: postID,
		Name:   user.Name,
		Email:  user.Email,
		Body:   commentRequest.Body,
	}
	comment, err := s.repo.SaveComment(domainComment)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service error save comment: %w", err)
	}
	return comment, nil
}

func (s commentService) GetComment(id int64) (models.Comment, error) {
	comment, err := s.repo.GetComment(id)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service error get comment: %w", err)
	}
	return comment, nil
}

func (s commentService) UpdateComment(commentRequest requests.CommentRequest, id int64) (models.Comment, error) {
	comment, err := s.repo.GetComment(id)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service error update comment: %w", err)
	}
	comment.Body = commentRequest.Body
	comment, err = s.repo.UpdateComment(comment)
	if err != nil {
		return models.Comment{}, fmt.Errorf("service error update comment: %w", err)
	}
	return comment, nil
}

func (s commentService) DeleteComment(id int64) error {
	_, err := s.repo.GetComment(id)
	if err != nil {
		return fmt.Errorf("service error delete comment: %w", err)
	}
	err = s.repo.DeleteComment(id)
	if err != nil {
		return fmt.Errorf("service error delete comment: %w", err)
	}
	return nil
}

func (s commentService) GetCommentsByPostID(postID int64, offset int) ([]models.Comment, error) {
	comments, err := s.repo.GetCommentsByPostID(postID, offset)
	if err != nil {
		return []models.Comment{}, fmt.Errorf("service error get all comments by postID: %w", err)
	}
	return comments, nil
}
