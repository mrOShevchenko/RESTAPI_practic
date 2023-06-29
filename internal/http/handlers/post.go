package handlers

import (
	"Nix_trainee_practic/internal/http/requests"
	"Nix_trainee_practic/internal/http/response"
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/service"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	service        service.PostService
	commentService service.CommentService
}

func NewPost(s service.PostService, c service.CommentService) Post {
	return Post{
		service:        s,
		commentService: c,
	}
}

// SavePost 		godoc
// @Summary 		Save Post
// @Description 	Save Post
// @Tags			Posts Actions
// @Accept 			json
// @Produce 		json
// @Param			input body requests.PostRequest true "comment info"
// @Success 		201 {object} response.PostResponse
// @Failure			400 {object} response.Error
// @Failure 		422 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/save [post]
func (p Post) SavePost(ctx echo.Context) error {
	var postRequest requests.PostRequest
	err := ctx.Bind(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode post data")
	}
	err = ctx.Validate(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate post data")
	}
	token := ctx.Get("user").(*jwt.Token)
	post, err := p.service.SavePost(postRequest, token)
	if err != nil {
		log.Print(err)
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not save new post: %s", err))
	}
	postResponse := models.Post.DomainToResponse(post)
	return response.Response(ctx, http.StatusCreated, postResponse)
}

// GetPost  		godoc
// @Summary 		Get Post
// @Description 	Get Post
// @Tags			Posts Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Param			offset path int true "Offset"
// @Success 		200 {object} response.PostResponse
// @Failure 		400 {object} response.Error
// @Failure 		404 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/post/{id} [get]
func (p Post) GetPost(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	offset, err := strconv.ParseInt(ctx.QueryParam("offset"), 10, 0)
	if err != nil {
		offset = 0
	}
	post, err := p.service.GetPost(id)
	if err != nil {
		log.Print(err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get post: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get post: %s", err))
		}
	}
	postComments, err := p.commentService.GetCommentsByPostID(id, int(offset))
	if err != nil {
		log.Print("There are no comments for this post")
	} else {
		dom := models.Comment{}
		post.Comments = dom.AllCommentsDomainToResponse(postComments)
	}
	postResponse := models.Post.DomainToResponse(post)
	return response.Response(ctx, http.StatusOK, postResponse)
}

// UpdatePost  		godoc
// @Summary 		Update Post
// @Description 	Update Post
// @Tags			Posts Actions
// @Accept 			json
// @Produce 		json
// @Param			id path int true "ID"
// @Param			input body requests.PostRequest true "post info"
// @Success 		200 {object} response.PostResponse
// @Failure 		400 {object} response.Error
// @Failure 		422 {object} response.Error
// @Failure 		400 {object} response.Error
// @Failure 		404 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/update/{id} [put]
func (p Post) UpdatePost(ctx echo.Context) error {
	var postRequest requests.PostRequest
	err := ctx.Bind(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode post data")
	}
	err = ctx.Validate(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate post data")
	}
	postID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	post, err := p.service.UpdatePost(postRequest, postID)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get post: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get post: %s", err))
		}
	}
	postResponse := models.Post.DomainToResponse(post)
	return response.Response(ctx, http.StatusOK, postResponse)
}

// DeletePost  		godoc
// @Summary 		Delete Post
// @Description 	Delete Post
// @Tags			Posts Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} response.Data
// @Failure 		400 {object} response.Error
// @Failure 		404 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/delete/{id} [delete]
func (p Post) DeletePost(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	err = p.service.DeletePost(id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get post: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get post: %s", err))
		}
	}
	return response.MessageResponse(ctx, http.StatusOK, "Post successfully delete")
}
