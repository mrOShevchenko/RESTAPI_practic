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

type Comment struct {
	service service.CommentService
}

func NewComment(s service.CommentService) Comment {
	return Comment{
		service: s,
	}
}

// SaveComment 		godoc
// @Summary 		Save Comment
// @Description 	Save Comment
// @Tags			Comments Actions
// @Accept 			json
// @Produce 		json
// @Param			post_id path int true "PostID"
// @Param			input body requests.CommentRequest true "comment info"
// @Success 		201 {object} response.CommentResponse
// @Failure			400 {object} response.Error
// @Failure 		422 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/comments/save/{post_id} [post]
func (c Comment) SaveComment(ctx echo.Context) error {
	var commentRequest requests.CommentRequest
	err := ctx.Bind(&commentRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode comment data")
	}
	err = ctx.Validate(&commentRequest)
	if err != nil {
		log.Print(err)
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate comment data")
	}
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	token := ctx.Get("user").(*jwt.Token)
	comment, err := c.service.SaveComment(commentRequest, postID, token)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not save new comment: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not save new comment: %s", err))
		}
	}
	commentResponse := models.Comment.DomainToResponse(comment)
	return response.Response(ctx, http.StatusCreated, commentResponse)
}

// GetComment 		godoc
// @Summary 		Get Comment
// @Description 	Get Comment
// @Tags			Comments Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} response.CommentResponse
// @Failure			400 {object} response.Error
// @Failure			404 {object} response.Error
// @Failure			500 {object} response.Error
// @Security 		ApiKeyAuth
// @Router			/api/v1/comments/comment/{id} [get]
func (c Comment) GetComment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse comment ID")
	}
	comment, err := c.service.GetComment(id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get comment: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get comment: %s", err))
		}
	}
	commentResponse := models.Comment.DomainToResponse(comment)
	return response.Response(ctx, http.StatusOK, commentResponse)
}

// UpdateComment 	godoc
// @Summary 		Update Comment
// @Description 	Update Comment
// @Tags			Comments Actions
// @Accept 			json
// @Produce 		json
// @Param			id path int true "ID"
// @Param			input body requests.CommentRequest true "comment info"
// @Success 		200 {object} response.CommentResponse
// @Failure			400 {object} response.Error
// @Failure			422 {object} response.Error
// @Failure			404 {object} response.Error
// @Security 		ApiKeyAuth
// @Router			/api/v1/comments/update/{id} [put]
func (c Comment) UpdateComment(ctx echo.Context) error {
	var commentRequest requests.CommentRequest
	err := ctx.Bind(&commentRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode comment data")
	}
	err = ctx.Validate(&commentRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate comment data")
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse comment ID")
	}
	comment, err := c.service.UpdateComment(commentRequest, id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not update comment: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not update comment: %s", err))
		}
	}
	commentResponse := models.Comment.DomainToResponse(comment)
	return response.Response(ctx, http.StatusOK, commentResponse)
}

// DeleteComment	godoc
// @Summary 		Delete Comment
// @Description 	Delete Comment
// @Tags			Comments Actions
// @Param			id path int true "ID"
// @Success 		200 {object} response.Data
// @Failure			400	{object} response.Error
// @Failure			404 {object} response.Error
// @Failure			500 {object} response.Error
// @Security 		ApiKeyAuth
// @Router			/api/v1/comments/delete/{id} [delete]
func (c Comment) DeleteComment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse comment ID")
	}
	err = c.service.DeleteComment(id)
	if err != nil {
		log.Print(err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not delete comment: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not delete comment: %s", err))
		}
	}
	return response.MessageResponse(ctx, http.StatusOK, "Comment Delete")
}
