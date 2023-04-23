package services

import (
	"MyGramAPI/app/entity"
	"MyGramAPI/pkg/database"
	"MyGramAPI/pkg/helpers"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetAllComment godoc
// @Summary Get all comments
// @Description User can retrieve all comments and no need to login
// @Tags comments
// @Consumes ({mpfd,json})
// @Produce json
// @Success 200 {object} entity.Response "Will send all comments"
// @Failure 404  {object}  entity.Response "If there is no comment, error will appear"
// @Router /api/v1/comments [GET]
func GetAllComment(c *gin.Context) {
	db, _ := database.Connect()
	Comment := []entity.Comment{}
	err := db.Find(&Comment).Error

	if err != nil {
		c.JSON(http.StatusNotFound, entity.Response{
			Success: false,
			Message: "There's no comment found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Comments has been loaded successfully",
		Data:    Comment,
	})
}

// GetComment godoc
// @Summary Get one comment
// @Description User can retrieve a comment and no need to login
// @Tags comments
// @Consumes ({mpfd,json})
// @Produce json
// @Param id path int true "comment id"
// @Success 200 {object} entity.Response "If a comment's id matches with the parameter"
// @Failure 404  {object}  entity.Response "If the comments's id doesn't match with the parameter, error will appear"
// @Router /api/v1/comments/{id} [GET]
func GetComment(c *gin.Context) {
	db, _ := database.Connect()
	contentType := helpers.GetContentType(c)
	Comment := entity.Comment{}

	//get parameter
	commentID, _ := strconv.Atoi(c.Param("id"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	//query select * from comment where id = param
	err := db.First(&Comment, "id = ?", commentID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, entity.Response{
			Success: false,
			Message: "Comment not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Comment has been loaded successfully",
		Data:    Comment,
	})
}

// CreateComment godoc
// @Summary Create a comment
// @Description User can create a comment.
// @Tags comments
// @Consumes ({mpfd,json})
// @Produce json
// @Param photo_id formData int true "photo id"
// @Param message formData string true "your comment"
// @Success 201 {object} entity.Response "If all of the parameters filled and you're login"
// @Failure 404 {object} entity.Response "If photo id's not found"
// @Failure 401  {object}  entity.Response "If you are not login or some parameters not filled, error will appear"
// @Security Bearer
// @Router /api/v1/comments [POST]
func CreateComment(c *gin.Context) {
	db, _ := database.Connect()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := entity.Comment{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	err := db.Debug().Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, entity.Response{
		Success: true,
		Message: "Comment has been created successfully",
		Data:    Comment,
	})
}

// UpdateComment godoc
// @Summary Edit a comment
// @Description User can edit their own comment.
// @Tags comments
// @Consumes ({mpfd,json})
// @Produce json
// @Param id path int true "comment id"
// @Param message formData string true "your comment"
// @Success 200 {object} entity.Response "If all the parameters are valid"
// @Failure 404  {object}  entity.Response "If there is something wrong, error will appear"
// @Security Bearer
// @Router /api/v1/comments/{id} [PUT]
func UpdateComment(c *gin.Context) {
	db, _ := database.Connect()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Comment := entity.Comment{}

	commentID, _ := strconv.Atoi(c.Param("id"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.ID = uint(commentID)

	err := db.Model(&Comment).Where("id = ?", commentID).Updates(entity.Comment{Message: Comment.Message}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Comment has been updated successfully",
		Data:    Comment,
	})
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description User can delete their own comment.
// @Tags comments
// @Consumes ({mpfd,json})
// @Produce json
// @Param id path int true "comment id"
// @Success 200 {object} entity.Response "If comment is exist and it's your own comment"
// @Failure 400  {object}  entity.Response "If the comment's id is not your own and if the comment doesn't exist, error will appear"
// @Security Bearer
// @Router /api/v1/comments/{id} [DELETE]
func DeleteComment(c *gin.Context) {
	db, _ := database.Connect()
	contentType := helpers.GetContentType(c)
	Comment := entity.Comment{}

	//get parameter
	commentID, _ := strconv.Atoi(c.Param("id"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	err := db.Where("id = ?", commentID).Delete(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Comment has been deleted successfully",
		Data:    nil,
	})
}
