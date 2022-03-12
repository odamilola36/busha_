package controllers

import (
	"busha_/dto"
	"busha_/service"
	"busha_/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler interface {
	CreateComment(w http.ResponseWriter, r *http.Request)
	GetAllComments(w http.ResponseWriter, r *http.Request)
	GetComment(w http.ResponseWriter, r *http.Request)
	UpdateComment(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
}

type commentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) CommentHandler {
	return &commentHandler{
		commentService: commentService,
	}
}

func (h *commentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment dto.CreateComment

	if r.Body == nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Please send a request body",
			Error:   "Invalid request body",
			Data:    nil,
		})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			dto.BuildResponse(w, &dto.Response{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
				Error:   err,
				Data:    nil,
			})
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err,
			Data:    nil,
		})
		return
	}

	if len(comment.Body) == 0 || len(comment.Body) > 500 {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Comment body should not be empty and should be less than 500 characters",
			Error:   "Invalid request body",
			Data:    nil,
		})
		return
	}
	if comment.MovieId == 0 {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Movie id is required",
			Error:   "Invalid request body",
			Data:    nil,
		})
		return
	}

	requestIp := strings.Split(r.RemoteAddr, ":")[0]
	commentResp, err := h.commentService.CreateComment(comment, requestIp)

	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Error:   err,
			Data:    nil,
		})
	}

	dto.BuildResponse(w, &dto.Response{
		Status:  http.StatusCreated,
		Message: "Comment created successfully",
		Error:   nil,
		Data:    commentResp,
	})
	return
}

func (h *commentHandler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	movieId, err := strconv.Atoi(pathVars["id"])

	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid movie id",
			Error:   err,
			Data:    nil,
		})
		return
	}

	comments, err := h.commentService.GetCommentsList(movieId)

	if err != nil {
		utils.ErrorLineLogger(err)
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Error getting comments",
			Error:   err,
			Data:    nil,
		})
		return
	}

	if len(comments) == 0 {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusNotFound,
			Message: "No comments found for this movie",
			Error:   nil,
			Data:    comments,
		})
		return
	}

	dto.BuildResponse(w, &dto.Response{
		Status:  http.StatusOK,
		Message: "Comments retrieved successfully",
		Error:   nil,
		Data:    comments,
	})
	return
}

func (h *commentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	commentId, err := strconv.Atoi(pathVars["commentId"])
	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid comment id",
			Error:   err,
			Data:    nil,
		})
		return
	}

	movieId, err := strconv.Atoi(pathVars["movieId"])
	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid movie id",
			Error:   err,
			Data:    nil,
		})
		return
	}

	comment, err := h.commentService.GetComment(movieId, commentId)

	if err != nil {
		utils.ErrorLineLogger(err)
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Error getting comment",
			Error:   err,
			Data:    nil,
		})
		return
	}

	dto.BuildResponse(w, &dto.Response{
		Status:  http.StatusOK,
		Message: "Comment retrieved successfully",
		Error:   nil,
		Data:    comment,
	})
}

func (h *commentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment dto.CreateComment

	s := mux.Vars(r)["commentId"]

	commentId, _ := strconv.Atoi(s)

	if r.Body == nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Please send a request body",
			Error:   "Invalid request body",
			Data:    nil,
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			dto.BuildResponse(w, &dto.Response{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
				Error:   err,
				Data:    nil,
			})
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err,
			Data:    nil,
		})
		return
	}

	if len(comment.Body) == 0 || len(comment.Body) > 500 {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Comment body should not be empty and should be less than 500 characters",
			Error:   "Invalid request body",
			Data:    nil,
		})
		return
	}

	updateComment, err := h.commentService.UpdateComment(comment, commentId)

	if err != nil {
		utils.ErrorLineLogger(err)
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Error updating comment",
			Error:   err,
			Data:    nil,
		})
		return
	}

	dto.BuildResponse(w, &dto.Response{
		Status:  http.StatusOK,
		Message: "Comment updated successfully",
		Error:   nil,
		Data:    updateComment,
	})
}

func (h *commentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	commentId, err := strconv.Atoi(pathVars["commentId"])
	if err != nil {
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid comment id",
			Error:   err,
			Data:    nil,
		})
		return
	}

	err = h.commentService.DeleteComment(commentId)

	if err != nil {
		if err.Error() == "record not found" {
			dto.BuildResponse(w, &dto.Response{
				Status:  http.StatusNotFound,
				Message: "Comment not found",
				Error:   err,
				Data:    nil,
			})
			return
		}
		utils.ErrorLineLogger(err)
		dto.BuildResponse(w, &dto.Response{
			Status:  http.StatusInternalServerError,
			Message: "Error deleting comment",
			Error:   err,
			Data:    nil,
		})
		return
	}

	dto.BuildResponse(w, &dto.Response{
		Status:  http.StatusOK,
		Message: "Comment deleted successfully",
		Error:   nil,
		Data:    nil,
	})
}
