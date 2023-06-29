package service

import (
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/repository"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
)

//go:generate mockery --dir . --name PostService --output ./mocks
type PostService interface {
	SavePost(postRequest requests.PostRequest, token *jwt.Token) (models.Post, error)
	GetPost(id int64) (models.Post, error)
	UpdatePost(postRequest requests.PostRequest, postID int64) (models.Post, error)
	DeletePost(id int64) error
	GetPostsByUser(userID int64) ([]models.Post, error)
}

type postService struct {
	repo repository.PostRepo
}

func NewPost(repo repository.PostRepo) PostService {
	return postService{
		repo: repo,
	}
}

func (s postService) SavePost(postRequest requests.PostRequest, token *jwt.Token) (models.Post, error) {
	claim := token.Claims.(*JWTClaim)
	userID := claim.ID
	domainPost := models.Post{
		Title:  postRequest.Title,
		Body:   postRequest.Body,
		UserID: userID,
	}
	post, err := s.repo.SavePost(domainPost)
	if err != nil {
		return models.Post{}, fmt.Errorf("service error save post: %w", err)
	}
	return post, nil
}

func (s postService) GetPost(id int64) (models.Post, error) {
	post, err := s.repo.GetPost(id)
	if err != nil {
		return models.Post{}, fmt.Errorf("service error get post: %w", err)
	}
	return post, nil
}

func (s postService) UpdatePost(postRequest requests.PostRequest, postID int64) (models.Post, error) {
	post, err := s.repo.GetPost(postID)
	if err != nil {
		return models.Post{}, fmt.Errorf("service error update post: %w", err)
	}

	post.Body = postRequest.Body
	post.Title = postRequest.Title

	post, err = s.repo.UpdatePost(post)
	if err != nil {
		log.Println(err)
		return models.Post{}, fmt.Errorf("service error update post: %w", err)
	}
	return post, nil
}

func (s postService) DeletePost(id int64) error {
	_, err := s.repo.GetPost(id)
	if err != nil {
		return fmt.Errorf("service error delete post: %w", err)
	}
	err = s.repo.DeletePost(id)
	if err != nil {
		return fmt.Errorf("service error delete post: %w", err)
	}
	return err
}

func (s postService) GetPostsByUser(userID int64) ([]models.Post, error) {
	posts, err := s.repo.GetPostsByUser(userID)
	if err != nil {
		return []models.Post{}, fmt.Errorf("service error get posts by user id: %w", err)
	}
	return posts, nil
}
