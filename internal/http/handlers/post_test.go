package handlers

import (
	"Nix_trainee_practic/internal/http/handlers/test_case"
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/service"
	"Nix_trainee_practic/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/upper/db/v4"
	"net/http"
	"testing"
)

const (
	postID      = "1"
	postIDError = "a"
)

var requestPostMock = requests.PostRequest{
	Title: "title",
	Body:  "body",
}
var returnModelsPostMock = models.Post{
	UserID: 1,
	ID:     1,
	Title:  "title",
	Body:   "body",
}

var requestGet = test_case.Request{
	Method: http.MethodGet,
	Url:    "/post/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postID,
	},
}

var requestSave = test_case.Request{
	Method: http.MethodPost,
	Url:    "/save",
}

var requestUpdate = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postID,
	},
}

var requestDelete = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postID,
	},
}

var requestGetError = test_case.Request{
	Method: http.MethodGet,
	Url:    "/post/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postIDError,
	},
}

var requestUpdateError = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postIDError,
	},
}

var requestDeleteError = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postIDError,
	},
}

func TestPost_GetPost(t *testing.T) {
	postGetOK := func(c echo.Context) error {
		mock := func(id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("GetPost", id).Return(returnModelsPostMock, nil).Times(1)
			return mock
		}(1)
		mockComment := func(id int64, offset int) service.CommentService {
			mockComment := mocks.NewCommentService(t)
			mockComment.On("GetCommentsByPostID", id, offset).Return([]models.Comment{}, nil).Times(1)
			return mockComment
		}(1, 0)
		return NewPost(mock, mockComment).GetPost(c)
	}

	postGetError := func(c echo.Context) error {
		mock := func() service.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).GetPost(c)
	}

	postGetNotFound := func(c echo.Context) error {
		mock := func(id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("GetPost", id).Return(models.Post{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).GetPost(c)
	}

	postGetInternalServerError := func(c echo.Context) error {
		mock := func(id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("GetPost", id).Return(models.Post{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).GetPost(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "GetPost OK",
			Request:     requestGet,
			RequestBody: "",
			HandlerFunc: postGetOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\",\"comments\":null}\n"},
		},
		{
			TestName:    "GetPost parse path param Error",
			Request:     requestGetError,
			RequestBody: "",
			HandlerFunc: postGetError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			TestName:    "GetPost NoMoreRows",
			Request:     requestGet,
			RequestBody: "",
			HandlerFunc: postGetNotFound,
			Expected:    test_case.ExpectedResponse{StatusCode: 404, BodyPart: "{\"code\":404,\"error\":\"Could not get post: upper: no more rows in this result set\"}\n"},
		},
		{
			TestName:    "GetPost InternalServerError",
			Request:     requestGet,
			RequestBody: "",
			HandlerFunc: postGetInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get post: upper: collection does not exist\"}\n"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(tt)
			c.Set("user", test_case.Token())

			if assert.NoError(t, tt.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), tt.Expected.BodyPart)
				assert.Equal(t, recorder.Code, tt.Expected.StatusCode)
			}
		})
	}

}

func TestPost_SavePost(t *testing.T) {
	postSaveOK := func(c echo.Context) error {
		mock := func(r requests.PostRequest, token *jwt.Token) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("SavePost", requestPostMock, token).Return(returnModelsPostMock, nil).Times(1)
			return mock
		}(requestPostMock, test_case.Token())
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).SavePost(c)
	}

	postSaveError := func(c echo.Context) error {
		mock := func() service.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).SavePost(c)
	}

	postSaveInternalServerError := func(c echo.Context) error {
		mock := func(r requests.PostRequest, token *jwt.Token) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("SavePost", requestPostMock, token).Return(models.Post{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestPostMock, test_case.Token())
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).SavePost(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "SavePost Success",
			Request:     requestSave,
			RequestBody: requestPostMock,
			HandlerFunc: postSaveOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\",\"comments\":null}\n"},
		},
		{
			TestName:    "SavePost post data error",
			Request:     requestSave,
			RequestBody: "",
			HandlerFunc: postSaveError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode post data\"}\n"},
		},
		{
			TestName: "SavePost validate error",
			Request:  requestSave,
			RequestBody: requests.PostRequest{
				Title: "title",
			},
			HandlerFunc: postSaveError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate post data\"}\n"},
		},
		{
			TestName:    "SavePost InternalServerError",
			Request:     requestSave,
			RequestBody: requestPostMock,
			HandlerFunc: postSaveInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not save new post: upper: collection does not exist\"}\n"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(tt)
			c.Set("user", test_case.Token())

			if assert.NoError(t, tt.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), tt.Expected.BodyPart)
				assert.Equal(t, recorder.Code, tt.Expected.StatusCode)
			}
		})
	}

}

func TestPost_UpdatePost(t *testing.T) {
	postUpdateOK := func(c echo.Context) error {
		mock := func(r requests.PostRequest, id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("UpdatePost", requestPostMock, id).Return(returnModelsPostMock, nil).Times(1)
			return mock
		}(requestPostMock, 1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).UpdatePost(c)
	}

	postUpdateError := func(c echo.Context) error {
		mock := func() service.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).UpdatePost(c)
	}

	postUpdateNotFound := func(c echo.Context) error {
		mock := func(r requests.PostRequest, id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("UpdatePost", requestPostMock, id).Return(models.Post{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(requestPostMock, 1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).UpdatePost(c)
	}

	postUpdateInternalServerError := func(c echo.Context) error {
		mock := func(r requests.PostRequest, id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("UpdatePost", requestPostMock, id).Return(models.Post{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestPostMock, 1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).UpdatePost(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "UpdatePost OK",
			Request:     requestUpdate,
			RequestBody: requestPostMock,
			HandlerFunc: postUpdateOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\",\"comments\":null}\n"},
		},
		{
			TestName:    "UpdatePost post data Error",
			Request:     requestUpdate,
			RequestBody: "",
			HandlerFunc: postUpdateError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode post data\"}\n"},
		},
		{
			TestName: "UpdatePost validate Error",
			Request:  requestUpdate,
			RequestBody: requests.PostRequest{
				Title: "title",
			},
			HandlerFunc: postUpdateError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate post data\"}\n"},
		},
		{
			TestName:    "UpdatePost parse path param Error",
			Request:     requestUpdateError,
			RequestBody: requestPostMock,
			HandlerFunc: postUpdateError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			TestName:    "UpdatePost NoMoreRows",
			Request:     requestUpdate,
			RequestBody: requestPostMock,
			HandlerFunc: postUpdateNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not get post: upper: no more rows in this result set\"}\n"},
		},
		{
			TestName:    "UpdatePost InternalServerError",
			Request:     requestUpdate,
			RequestBody: requestPostMock,
			HandlerFunc: postUpdateInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get post: upper: collection does not exist\"}\n"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(tt)
			c.Set("user", test_case.Token())

			if assert.NoError(t, tt.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), tt.Expected.BodyPart)
				assert.Equal(t, recorder.Code, tt.Expected.StatusCode)
			}
		})
	}
}

func TestPost_DeletePost(t *testing.T) {
	postDeleteOK := func(c echo.Context) error {
		mock := func(id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("DeletePost", id).Return(nil).Times(1)
			return mock
		}(1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).DeletePost(c)
	}

	postDeleteError := func(c echo.Context) error {
		mock := func() service.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).DeletePost(c)
	}

	postDeleteNotFound := func(c echo.Context) error {
		mock := func(id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("DeletePost", id).Return(db.ErrNoMoreRows).Times(1)
			return mock
		}(1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).DeletePost(c)
	}

	postDeleteInternalServerError := func(c echo.Context) error {
		mock := func(id int64) service.PostService {
			mock := mocks.NewPostService(t)
			mock.On("DeletePost", id).Return(db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(1)
		mockComment := mocks.NewCommentService(t)
		return NewPost(mock, mockComment).DeletePost(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "DeletePost OK",
			Request:     requestDelete,
			RequestBody: "",
			HandlerFunc: postDeleteOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
			},
		},
		{
			TestName:    "DeletePost parse path param Error",
			Request:     requestDeleteError,
			RequestBody: "",
			HandlerFunc: postDeleteError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
			},
		},
		{
			TestName:    "DeletePost NoMoreRows",
			Request:     requestDelete,
			RequestBody: "",
			HandlerFunc: postDeleteNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
			},
		},
		{
			TestName:    "DeletePost InternalServerError",
			Request:     requestDelete,
			RequestBody: "",
			HandlerFunc: postDeleteInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
			},
		},
	}

	for _, tt := range testTable {
		c, recorder := test_case.PrepareContextFromTestCase(tt)
		c.Set("user", test_case.Token())

		if assert.NoError(t, tt.HandlerFunc(c)) {
			assert.Contains(t, recorder.Body.String(), tt.Expected.BodyPart)
			assert.Equal(t, recorder.Code, tt.Expected.StatusCode)
		}
	}

}
