package handlers

import (
	"Nix_trainee_practic/internal/http/handlers/test_case"
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/http/response"
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
	commentIDError = "a"
	commentID      = "2"
)

var requestCommentMock = requests.CommentRequest{
	Body: "Test body",
}

var returnModelsCommentMock = models.Comment{
	ID:     2,
	PostID: 1,
	Name:   "testName",
	Email:  "test@mail.com",
	Body:   "testBody",
}

var requestSaveCommentMock = test_case.Request{
	Method: http.MethodGet,
	Url:    "/save/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "post_id",
		Value: postID,
	},
}

var requestGetCommentMock = test_case.Request{
	Method: http.MethodGet,
	Url:    "/comment/" + commentID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentID,
	},
}

var requestUpdateComment = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + commentID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentID,
	},
}

var requestDeleteComment = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + commentID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentID,
	},
}

var requestSaveCommentError = test_case.Request{
	Method: http.MethodGet,
	Url:    "/save/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "post_id",
		Value: postIDError,
	},
}

var requestGetCommentError = test_case.Request{
	Method: http.MethodGet,
	Url:    "/comment/" + commentIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentIDError,
	},
}

var requestUpdateCommentError = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + commentIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentIDError,
	},
}

var requestDeleteCommentError = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + commentIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentIDError,
	},
}

func TestComment_GetComment(t *testing.T) {
	commentGetOK := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(returnModelsCommentMock, nil).Times(1)
			return mock
		}(2)
		return NewComment(mock).GetComment(c)
	}

	commentGetErrorPath := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(models.Comment{}, response.ErrorResponse(c, http.StatusBadRequest, "Could not parse post ID")).Times(1)
			return mock
		}(2)

		return NewComment(mock).GetComment(c)

	}

	commentGetNotFound := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(models.Comment{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(2)
		return NewComment(mock).GetComment(c)
	}

	commentGetInternalServerError := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(models.Comment{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(2)
		return NewComment(mock).GetComment(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "OK",
			Request:     requestGetCommentMock,
			RequestBody: "",
			HandlerFunc: commentGetOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"testName\",\"email\":\"test@mail.com\",\"body\":\"testBody\"}\n",
			},
		},
		{
			TestName:    "Error GetComment path param",
			Request:     requestGetCommentMock,
			RequestBody: "",
			HandlerFunc: commentGetErrorPath,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n{\"id\":0,\"post_id\":0,\"name\":\"\",\"email\":\"\",\"body\":\"\"}\n",
			},
		},
		{
			TestName:    "GetComment NoMoreRows",
			Request:     requestGetCommentMock,
			RequestBody: "",
			HandlerFunc: commentGetNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not get comment: upper: no more rows in this result set\"}\n",
			},
		},
		{
			TestName:    "GetComment InternalServerError",
			Request:     requestGetCommentMock,
			RequestBody: "",
			HandlerFunc: commentGetInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get comment: upper: collection does not exist\"}\n"},
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

func TestComment_SaveComment(t *testing.T) {
	commentSaveOK := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, token *jwt.Token, postID int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("SaveComment", r, postID, token).Return(returnModelsCommentMock, nil).Times(1)
			return mock
		}(requestCommentMock, test_case.Token(), 1)
		return NewComment(mock).SaveComment(c)
	}

	commentSaveError := func(c echo.Context) error {
		mock := func() service.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return NewComment(mock).SaveComment(c)
	}

	commentSaveNotFound := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, p int64, token *jwt.Token) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("SaveComment", r, p, token).Return(models.Comment{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(requestCommentMock, 1, test_case.Token())
		return NewComment(mock).SaveComment(c)
	}

	commentSaveInternalServerError := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, p int64, token *jwt.Token) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("SaveComment", r, p, token).Return(models.Comment{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestCommentMock, 1, test_case.Token())
		return NewComment(mock).SaveComment(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "OK SaveComment",
			Request:     requestSaveCommentMock,
			RequestBody: requestCommentMock,
			HandlerFunc: commentSaveOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"testName\",\"email\":\"test@mail.com\",\"body\":\"testBody\"}\n"},
		},
		{
			TestName:    "SaveComment Error",
			Request:     requestSaveCommentMock,
			RequestBody: "",
			HandlerFunc: commentSaveError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode comment data\"}\n"},
		},
		{
			TestName: "SaveComment validate Error",
			Request:  requestSaveCommentMock,
			RequestBody: requests.PostRequest{
				Title: "title",
			},
			HandlerFunc: commentSaveError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate comment data\"}\n"},
		},
		{
			TestName:    "SaveComment parse path param Error",
			Request:     requestSaveCommentError,
			RequestBody: requestCommentMock,
			HandlerFunc: commentSaveError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			TestName:    "SaveComment NoMoreRows",
			Request:     requestSaveCommentMock,
			RequestBody: requestCommentMock,
			HandlerFunc: commentSaveNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not save new comment: upper: no more rows in this result set\"}\n"},
		},
		{
			TestName:    "SaveComment InternalServerError",
			Request:     requestSaveCommentMock,
			RequestBody: requestCommentMock,
			HandlerFunc: commentSaveInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not save new comment: upper: collection does not exist\"}\n"},
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

func TestComment_UpdateComment(t *testing.T) {
	commentUpdateOK := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			returnModelsCommentMock.Body = r.Body
			mock.On("UpdateComment", r, id).Return(returnModelsCommentMock, nil).Times(1)
			return mock
		}(requests.CommentRequest{Body: "Update body"}, 2)
		return NewComment(mock).UpdateComment(c)
	}

	commentUpdateError := func(c echo.Context) error {
		mock := func() service.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return NewComment(mock).UpdateComment(c)
	}

	commentUpdateNotFound := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("UpdateComment", r, id).Return(models.Comment{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(requestCommentMock, 2)
		return NewComment(mock).UpdateComment(c)
	}

	commentUpdateInternalServerError := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("UpdateComment", r, id).Return(models.Comment{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestCommentMock, 2)
		return NewComment(mock).UpdateComment(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "UpdateComment OK",
			Request:     requestUpdateComment,
			RequestBody: requests.CommentRequest{Body: "Update body"},
			HandlerFunc: commentUpdateOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"testName\",\"email\":\"test@mail.com\",\"body\":\"Update body\"}\n"},
		},
		{
			TestName:    "UpdateComment comment data error",
			Request:     requestUpdateComment,
			RequestBody: "",
			HandlerFunc: commentUpdateError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode comment data\"}\n"},
		},
		{
			TestName: "UpdateComment validate error",
			Request:  requestUpdateComment,
			RequestBody: requests.PostRequest{
				Title: "title",
			},
			HandlerFunc: commentUpdateError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate comment data\"}\n"},
		},
		{
			TestName:    "UpdateComment parse path param Error",
			Request:     requestUpdateCommentError,
			RequestBody: requestPostMock,
			HandlerFunc: commentUpdateError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse comment ID\"}\n"},
		},
		{
			TestName:    "UpdateComment NoMoreRows",
			Request:     requestUpdateComment,
			RequestBody: requestCommentMock,
			HandlerFunc: commentUpdateNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not update comment: upper: no more rows in this result set\"}\n"},
		}, {
			TestName:    "UpdateComment InternalServerError",
			Request:     requestUpdateComment,
			RequestBody: requestCommentMock,
			HandlerFunc: commentUpdateInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not update comment: upper: collection does not exist\"}\n"},
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

func TestComment_DeleteComment(t *testing.T) {
	commentDeleteOK := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("DeleteComment", id).Return(nil).Times(1)
			return mock
		}(2)
		return NewComment(mock).DeleteComment(c)
	}

	commentDelete := func(c echo.Context) error {
		mock := func() service.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return NewComment(mock).DeleteComment(c)
	}

	commentDeleteNotFound := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("DeleteComment", id).Return(db.ErrNoMoreRows).Times(1)
			return mock
		}(2)
		return NewComment(mock).DeleteComment(c)
	}

	commentDeleteInternalServerError := func(c echo.Context) error {
		mock := func(id int64) service.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("DeleteComment", id).Return(db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(2)
		return NewComment(mock).DeleteComment(c)
	}

	testTable := []test_case.TestCase{
		{
			TestName:    "DeleteComment OK",
			Request:     requestDeleteComment,
			RequestBody: "",
			HandlerFunc: commentDeleteOK,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
			},
		},
		{
			TestName:    "DeleteComment parse path param Error",
			Request:     requestDeleteCommentError,
			RequestBody: "",
			HandlerFunc: commentDelete,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
			},
		},
		{
			TestName:    "DeleteComment NoMoreRows",
			Request:     requestDeleteComment,
			RequestBody: "",
			HandlerFunc: commentDeleteNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
			},
		},
		{
			TestName:    "DeleteComment InternalServerError",
			Request:     requestDeleteComment,
			RequestBody: "",
			HandlerFunc: commentDeleteInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
			},
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
