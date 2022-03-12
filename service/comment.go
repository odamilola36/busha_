package service

import (
	"busha_/dto"
	"busha_/models"
	"busha_/utils"
	"gorm.io/gorm"
)

type CommentService interface {
	CreateComment(comment dto.CreateComment, ip string) (models.Comment, error)
	GetCommentsList(movieId int) ([]models.Comment, error)
	GetComment(movieId, commentId int) (models.Comment, error)
	DeleteComment(commentId int) error
	UpdateComment(comment dto.CreateComment, commentId int) (models.Comment, error)
}

type commentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) CommentService {
	return &commentService{
		db: db,
	}
}

func (s *commentService) CreateComment(comment dto.CreateComment, ip string) (models.Comment, error) {
	var newComment models.Comment

	newComment.MovieId = comment.MovieId
	newComment.Body = comment.Body
	newComment.CommenterIP = ip

	if err := s.db.Create(&newComment).Error; err != nil {
		return newComment, err
	}

	return newComment, nil
}

func (s *commentService) GetCommentsList(movieId int) ([]models.Comment, error) {

	//get comments list from the database order by creation date
	var comments []models.Comment

	err := s.db.Where("movie_id = ?", movieId).Order("created_at desc").Find(&comments).Error

	if err != nil {
		if err.Error() == "record not found" {
			return comments, nil
		}
		return comments, err
	}

	return comments, nil
}

func (s *commentService) GetComment(movieId, commentId int) (models.Comment, error) {
	var comment models.Comment

	err := s.db.Where("movie_id = ? AND id = ?", movieId, commentId).First(&comment).Error

	if err != nil {
		if err.Error() == "record not found" {
			utils.ErrorLineLogger(err)
			return comment, nil
		}
		utils.ErrorLineLogger(err)
		return comment, err
	}

	return comment, nil

}

func (s *commentService) DeleteComment(commentId int) error {
	var comment models.Comment

	err := s.db.Where("id = ?", commentId).First(&comment).Error

	if err != nil {
		if err.Error() == "record not found" {
			utils.ErrorLineLogger(err)
			return err
		}
		utils.ErrorLineLogger(err)
		return err
	}

	err = s.db.Delete(&comment).Error

	if err != nil {
		utils.ErrorLineLogger(err)
		return err
	}

	return nil
}

func (s *commentService) UpdateComment(comment dto.CreateComment, commentId int) (models.Comment, error) {
	var commentToUpdate models.Comment

	err := s.db.Where("id = ?", commentId).First(&commentToUpdate).Error

	if err != nil {
		if err.Error() == "record not found" {
			utils.ErrorLineLogger(err)
			return commentToUpdate, err
		}
		utils.ErrorLineLogger(err)
		return commentToUpdate, err
	}

	commentToUpdate.Body = comment.Body

	err = s.db.Save(&commentToUpdate).Error

	if err != nil {
		utils.ErrorLineLogger(err)
		return commentToUpdate, err
	}

	return commentToUpdate, nil
}
