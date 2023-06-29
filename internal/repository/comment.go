package repository

import (
	"Nix_trainee_practic/internal/models"
	"fmt"
	"github.com/upper/db/v4"
	"time"
)

const CommentTable = "commentses"

type comments struct {
	ID          int64      `db:"id,omitempty"`
	PostID      int64      `db:"post_id"`
	Name        string     `db:"name"`
	Email       string     `db:"email"`
	Body        string     `db:"body"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

//go:generate mockery --dir . --name CommentRepo --output ./mock
type CommentRepo interface {
	SaveComment(comment models.Comment) (models.Comment, error)
	GetComment(id int64) (models.Comment, error)
	UpdateComment(comment models.Comment) (models.Comment, error)
	DeleteComment(id int64) error
	GetCommentsByPostID(postID int64, offset int) ([]models.Comment, error)
}

type commentsRepository struct {
	coll db.Collection
}

func NewCommentRepo(dbSession db.Session) CommentRepo {
	return commentsRepository{
		coll: dbSession.Collection(CommentTable),
	}
}

func (r commentsRepository) SaveComment(comment models.Comment) (models.Comment, error) {
	commentsDB := r.mapCommentDBModel(comment)
	commentsDB.CreatedDate = time.Now()
	commentsDB.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&commentsDB)
	if err != nil {
		return models.Comment{}, fmt.Errorf("comment repository save comment: %w", err)
	}
	return r.mapCommentDbModelToDomain(commentsDB), nil
}

func (r commentsRepository) GetComment(id int64) (models.Comment, error) {
	var comment comments

	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).One(&comment)
	if err != nil {
		return models.Comment{}, fmt.Errorf("comment repository get comment: %w", err)
	}
	return r.mapCommentDbModelToDomain(comment), nil
}

func (r commentsRepository) UpdateComment(comment models.Comment) (models.Comment, error) {
	updateComment := r.mapCommentDBModel(comment)
	updateComment.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{
		"id": updateComment.ID,
	}).Update(&updateComment)
	if err != nil {
		return models.Comment{}, fmt.Errorf("comment repository update comment: %w", err)
	}
	return r.mapCommentDbModelToDomain(updateComment), err
}

func (r commentsRepository) DeleteComment(id int64) error {
	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		return fmt.Errorf("comment repository delete comment: %w", err)
	}
	return nil
}

func (r commentsRepository) GetCommentsByPostID(postID int64, offset int) ([]models.Comment, error) {
	var comment []comments

	err := r.coll.Find(db.Cond{"post_id": postID}).Offset(offset).Limit(10).All(&comment)
	if err != nil {
		return []models.Comment{}, fmt.Errorf("comment repository GetCommentsByPostID: %w", err)
	}
	return r.mapCommentCollection(comment), nil
}

func (r commentsRepository) mapCommentDBModel(comment models.Comment) comments {
	return comments{
		ID:     comment.ID,
		PostID: comment.PostID,
		Name:   comment.Name,
		Email:  comment.Email,
		Body:   comment.Body,
	}
}

func (r commentsRepository) mapCommentDbModelToDomain(comment comments) models.Comment {
	return models.Comment{
		ID:          comment.ID,
		PostID:      comment.PostID,
		Name:        comment.Name,
		Email:       comment.Email,
		Body:        comment.Body,
		CreatedDate: comment.CreatedDate,
		DeletedDate: comment.DeletedDate,
		UpdatedDate: comment.UpdatedDate,
	}
}

func (r commentsRepository) mapCommentCollection(comment []comments) []models.Comment {
	var result []models.Comment
	for _, coll := range comment {
		newComment := r.mapCommentDbModelToDomain(coll)
		result = append(result, newComment)
	}
	return result
}
