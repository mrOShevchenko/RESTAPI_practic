package repository

import (
	"Nix_trainee_practic/internal/models"
	"fmt"
	"github.com/upper/db/v4"
	"time"
)

const PostTable = "posts"

type posts struct {
	UserID      int64      `db:"user_id"`
	ID          int64      `db:"id,omitempty"`
	Title       string     `db:"title"`
	Body        string     `db:"body"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

//go:generate mockery --dir . --name PostRepo --output ./mock
type PostRepo interface {
	SavePost(post models.Post) (models.Post, error)
	GetPost(id int64) (models.Post, error)
	GetPostsByUser(userID int64) ([]models.Post, error)
	UpdatePost(post models.Post) (models.Post, error)
	DeletePost(id int64) error
}

type postsRepository struct {
	coll db.Collection
}

func NewPostRepo(dbSession db.Session) PostRepo {
	return postsRepository{
		coll: dbSession.Collection(PostTable),
	}
}

func (r postsRepository) SavePost(post models.Post) (models.Post, error) {
	postDB := r.mapPostDBModel(post)
	postDB.CreatedDate = time.Now()
	postDB.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&postDB)
	if err != nil {
		return models.Post{}, fmt.Errorf("post repository save post: %w", err)
	}
	return r.mapPostDbModelToDomain(postDB), nil
}

func (r postsRepository) GetPost(id int64) (models.Post, error) {
	var post posts

	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).One(&post)
	if err != nil {
		return models.Post{}, fmt.Errorf("post repository get post: %w", err)
	}
	return r.mapPostDbModelToDomain(post), nil
}

func (r postsRepository) UpdatePost(post models.Post) (models.Post, error) {
	updatePost := r.mapPostDBModel(post)
	updatePost.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{
		"id": updatePost.ID,
	}).Update(&updatePost)
	if err != nil {
		return models.Post{}, fmt.Errorf("post repository update post: %w", err)
	}
	return r.mapPostDbModelToDomain(updatePost), err
}

func (r postsRepository) DeletePost(id int64) error {
	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		return fmt.Errorf("post repository delete post: %w", err)
	}
	return nil
}

func (r postsRepository) GetPostsByUser(userID int64) ([]models.Post, error) {
	var post []posts

	err := r.coll.Find(db.Cond{"user_id": userID}).All(&post)
	if err != nil {
		return []models.Post{}, fmt.Errorf("post repository get post by user: %w", err)
	}
	return r.mapPostCollection(post), nil

}

func (r postsRepository) mapPostDBModel(p models.Post) posts {
	return posts{
		UserID: p.UserID,
		ID:     p.ID,
		Title:  p.Title,
		Body:   p.Body,
	}
}

func (r postsRepository) mapPostDbModelToDomain(p posts) models.Post {
	return models.Post{
		UserID:      p.UserID,
		ID:          p.ID,
		Title:       p.Title,
		Body:        p.Body,
		CreatedDate: p.CreatedDate,
		UpdatedDate: p.UpdatedDate,
		DeletedDate: p.DeletedDate,
	}
}

func (r postsRepository) mapPostCollection(post []posts) []models.Post {
	var result []models.Post
	for _, coll := range post {
		newPost := r.mapPostDbModelToDomain(coll)
		result = append(result, newPost)
	}
	return result
}
