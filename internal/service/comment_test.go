package service

import (
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/repository"
	"Nix_trainee_practic/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/upper/db/v4"
	"testing"
)

func TestCommentService_GetComment(t *testing.T) {
	testTable := []struct {
		name      string
		id        int64
		repo      func(id int64) repository.CommentRepo
		expected  models.Comment
		expectErr bool
	}{
		{
			"OK GetComment",
			2,
			func(id int64) repository.CommentRepo {
				mock := mocks.NewCommentRepo(t)
				mock.On("GetComment", id).
					Return(models.Comment{ID: 2}, nil)
				return mock
			},
			models.Comment{ID: 2},
			false,
		},
		{
			"Error GetComment",
			2,
			func(id int64) repository.CommentRepo {
				mock := mocks.NewCommentRepo(t)
				mock.On("GetComment", id).
					Return(models.Comment{}, db.ErrNoMoreRows)
				return mock
			},
			models.Comment{}, // Ожидаемое значение ID изменено на 2
			true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			s := commentService{
				repo: tt.repo(tt.id),
			}
			comment, err := NewComment(s.repo, s.us, s.ps).GetComment(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, comment)
			}
		})
	}

}

func TestCommentService_SaveComment(t *testing.T) {
	testTable := []struct {
		name           string
		commentRequest requests.CommentRequest
		postID         int64
		token          *jwt.Token
		ps             func(postID int64) PostService
		us             func(id int64) UserService
		repo           func(commentRequest requests.CommentRequest) repository.CommentRepo
		expect         models.Comment
		expectErr      bool
	}{
		{
			name:           "OK create comment",
			commentRequest: requests.CommentRequest{Body: "body"},
			postID:         2,
			token:          token(),
			ps: func(postID int64) PostService {
				mock := NewMockPostService(t)
				mock.On("GetPost", postID).
					Return(models.Post{}, nil).Times(1)
				return mock
			},
			us: func(id int64) UserService {
				mock := NewMockUserService(t)
				mock.On("FindByID", id).
					Return(models.User{
						ID:    id,
						Email: "test@mail.com",
						Name:  "Name",
					}, nil).Times(1)
				return mock
			},
			repo: func(commentRequest requests.CommentRequest) repository.CommentRepo {
				mock := mocks.NewCommentRepo(t)
				comment := models.Comment{
					PostID: 2,
					Name:   "Name",
					Email:  "test@mail.com",
					Body:   "body",
				}
				mock.On("SaveComment", comment).
					Return(comment, nil)
				return mock
			},
			expect: models.Comment{
				PostID: 2,
				Name:   "Name",
				Email:  "test@mail.com",
				Body:   "body",
			},
			expectErr: false,
		},
		{
			name:           "ERROR post not exist",
			commentRequest: requests.CommentRequest{Body: "body"},
			postID:         2,
			token:          token(),
			ps: func(postID int64) PostService {
				mock := NewMockPostService(t)
				mock.On("GetPost", postID).
					Return(models.Post{}, db.ErrNoMoreRows).Times(1)
				return mock
			},
			us: func(id int64) UserService {
				mock := NewMockUserService(t)
				return mock
			},
			repo: func(commentRequest requests.CommentRequest) repository.CommentRepo {
				mock := mocks.NewCommentRepo(t)
				return mock
			},
			expect:    models.Comment{},
			expectErr: true,
		},
		{
			name:           "ERROR user not exist",
			commentRequest: requests.CommentRequest{Body: "body"},
			postID:         2,
			token:          token(),
			ps: func(postID int64) PostService {
				mock := NewMockPostService(t)
				mock.On("GetPost", postID).
					Return(models.Post{}, nil).Times(1)
				return mock
			},
			us: func(id int64) UserService {
				mock := NewMockUserService(t)
				mock.On("FindByID", id).
					Return(models.User{}, db.ErrNoMoreRows).Times(1)
				return mock
			},
			repo: func(commentRequest requests.CommentRequest) repository.CommentRepo {
				mock := mocks.NewCommentRepo(t)
				return mock
			},
			expect:    models.Comment{},
			expectErr: true,
		},
		{
			name:           "ERROR with repository",
			commentRequest: requests.CommentRequest{Body: "body"},
			postID:         2,
			token:          token(),
			ps: func(postID int64) PostService {
				mock := NewMockPostService(t)
				mock.On("GetPost", postID).
					Return(models.Post{}, nil).Times(1)
				return mock
			},
			us: func(id int64) UserService {
				mock := NewMockUserService(t)
				mock.On("FindByID", id).
					Return(models.User{
						ID:    id,
						Email: "test@mail.com",
						Name:  "Name",
					}, nil).Times(1)
				return mock
			},
			repo: func(commentRequest requests.CommentRequest) repository.CommentRepo {
				mock := mocks.NewCommentRepo(t)
				comment := models.Comment{
					PostID: 2,
					Name:   "Name",
					Email:  "test@mail.com",
					Body:   "body",
				}
				mock.
					On("SaveComment", comment).
					Return(models.Comment{}, db.ErrMissingPrimaryKeys)
				return mock
			},
			expect:    models.Comment{},
			expectErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			s := commentService{
				repo: tt.repo(tt.commentRequest),
				us:   tt.us(tt.token.Claims.(*JWTClaim).ID),
				ps:   tt.ps(tt.postID),
			}
			comment, err := NewComment(s.repo, s.us, s.ps).SaveComment(tt.commentRequest, tt.postID, tt.token)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, comment, tt.expect)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, comment, tt.expect)
			}
		})
	}
}
