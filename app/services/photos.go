package services

import (
	"MyGramAPI/app/entity"
	"MyGramAPI/pkg/database"
	"MyGramAPI/pkg/helpers"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetAllPhotos godoc
// @Summary Get all photos
// @Description User can retrieve all photos and no need to login
// @Tags photos
// @Consumes ({mpfd,json})
// @Produce json
// @Success 200 {object} entity.Response "Will send all photos"
// @Failure 404  {object}  entity.Response "If there is no photos, error will appear"
// @Router /api/v1/photos [GET]

func GetAllPhoto(c *gin.Context) {
	db, _ := database.Connect()

	Photo := []entity.Photo{}
	User := entity.User{}

	ResData := []entity.DataPhoto{}
	err := db.Order("created_at desc").Find(&Photo).Error
	for _, photo := range Photo {
		var username string
		db.Select("username").First(&User, int(photo.UserID)).Scan(&username)

		Comment := []entity.Comment{}
		ResComment := []entity.DataComment{}
		db.Find(&Comment, "my_gram_photo_id", int(photo.ID))
		for _, comment := range Comment {
			var uname string
			db.Select("username").First(&User, int(comment.UserID)).Scan(&uname)

			ResComment = append(ResComment, entity.DataComment{
				ID:        comment.ID,
				Message:   comment.Message,
				Username:  uname,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			})
		}

		ResData = append(ResData, entity.DataPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			UserID:    photo.UserID,
			Username:  username,
			Photo_URL: photo.Photo_URL,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			Comment:   ResComment,
		})
	}

	if err != nil {
		c.JSON(http.StatusNotFound, entity.Response{
			Success: false,
			Message: "There's no photo found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Photos has been loaded successfully",
		Data:    ResData,
	},
	)
}

// GetPhoto godoc
// @Summary Get one photo
// @Description User can retrieve a photo and no need to login
// @Tags photos
// @Consumes ({mpfd,json})
// @Produce json
// @Param id path int true "photo id"
// @Success 200 {object} entity.Response "If a photo's id matches with the parameter"
// @Failure 404  {object}  entity.Response "If the photo's id doesn't match with the parameter, error will appear"
// @Router /api/v1/photos/{id} [GET]
func GetPhoto(c *gin.Context) {
	db, _ := database.Connect()
	contentType := helpers.GetContentType(c)
	Photo := entity.Photo{}

	//get parameter
	photoID, _ := strconv.Atoi(c.Param("id"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	//query select * from photo where id = param
	err := db.First(&Photo, "id = ?", photoID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, entity.Response{
			Success: false,
			Message: "Photo not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, entity.Response{
		Success: true,
		Message: "Photo has been loaded successfully",
		Data:    Photo,
	})
}

// CreatePhoto godoc
// @Summary Upload a photo
// @Description User can upload a photo.
// @Tags photos
// @Consumes ({mpfd,json})
// @access-control-allow-origin *
// @Produce json
// @Param title formData string true "photo title"
// @Param caption formData string true "photo caption"
// @Param photo_url formData file true "photo url"
// @Param Authorization header string true "format: Bearer token-here"
// @Success 201 {object} entity.Response "If all of the parameters filled and you're logged in"
// @Failure 404  {object}  entity.Response "If you are not login or some parameters not filled, error will appear"
// @Security Bearer
// @Router /api/v1/photos [POST]
func CreatePhoto(c *gin.Context) {
	
	var photoFileHeader *multipart.FileHeader
	
	db, _ := database.Connect()
	contentType := helpers.GetContentType(c)
	Photo := entity.Photo{}
	
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	Photo.UserID = userID

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	// photo source, check if photo is uploaded
	photoFileHeader, err := c.FormFile("photo_url")
	if err != nil {
		log.Printf("get form err - %s", err.Error())
		respon := helpers.ApiResponse("No photo file uploaded", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, respon)
		return
	}

	//chech if the file is an image or not
	isPhoto := filepath.Ext(photoFileHeader.Filename)
	if isPhoto != ".jpg" && isPhoto != ".jpeg" && isPhoto != ".png" && isPhoto != ".webp"  {
		respon := helpers.ApiResponse("File uploaded is not an image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, respon)
		return
	}
	
	// open the file and get its content
	photoFile, err := photoFileHeader.Open()
	if err != nil {
		log.Printf("error opening file: %v", err)
		return
	}
	defer photoFile.Close()

	//upload photo to cloudinary
	photoSource, err := helpers.UploadToCloudinary(photoFile)
	if err != nil {
		return
	}	

	Photo = entity.Photo{
		Title:     Photo.Title,
		Caption:   Photo.Caption,
		UserID: userID,
		Photo_URL: photoSource,
	}


	err = db.Debug().Create(&Photo).Error

	if err != nil {
		response := helpers.ApiResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helpers.ApiResponse("Photo has been created successfully", http.StatusCreated, "success", Photo)
	c.JSON(http.StatusCreated, response)
}

// UpdatePhoto godoc
// @Summary Edit a photo
// @Description User can edit their own photo.
// @Tags photos
// @Consumes ({mpfd,json})
// @Produce json
// @Param id path int true "photo id"
// @Param title formData string true "photo title"
// @Param caption formData string true "photo caption"
// @Param photo_url formData string true "photo url"
// @Success 200 {object} entity.Response "If the parameters are valid"
// @Failure 401  {object}  entity.Response "If there is something wrong, error will appear"
// @Security Bearer
// @Router /api/v1/photos/{id} [PUT]
func UpdatePhoto(c *gin.Context) {
	var photoFileHeader *multipart.FileHeader

	db, _ := database.Connect()

	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Photo := entity.Photo{}

	photoID, _ := strconv.Atoi(c.Param("id"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	photoFileHeader, err := c.FormFile("photo_url")
	
	if err != nil {
		log.Printf("get form err - %s", err.Error())
		respon := helpers.ApiResponse("No photo file uploaded", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, respon)
		return
	}

	if photoFileHeader != nil {
		//chech if the file is an image or not
		isPhoto := filepath.Ext(photoFileHeader.Filename)
		if isPhoto != ".jpg" && isPhoto != ".jpeg" && isPhoto != ".png" && isPhoto != ".webp"  {
			respon := helpers.ApiResponse("File uploaded is not an image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, respon)
			return
		}

		// open the file and get its content
		photoFile, err := photoFileHeader.Open()
		if err != nil {
			log.Printf("error opening file: %v", err)
			return
		}
		defer photoFile.Close()

		//upload photo to cloudinary
		photoSource, err := helpers.UploadToCloudinary(photoFile)
		if err != nil {
			return
		}

		Photo.Photo_URL = photoSource
	}

	Photo.UserID = userID
	Photo.ID = uint(photoID)

	err = db.Model(&Photo).Where("id = ?", photoID).Updates(entity.Photo{Title: Photo.Title, Caption: Photo.Caption, Photo_URL: Photo.Photo_URL}).Error

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
		Message: "Photo has been updated successfully",
		Data:    Photo,
	})

}

// DeletePhoto godoc
// @Summary Delete a photo
// @Description User can delete their own photo.
// @Tags photos
// @Consumes ({mpfd,json})
// @Produce json
// @Param id path int true "photo id"
// @Success 200 {object} entity.Response "If photo is exist and it's your own photo, photo will deleted"
// @Failure 400  {object}  entity.Response "If the photo is not your own or if the photo doesn't exist, error will appear"
// @Security Bearer
// @Router /api/v1/photos/{id} [DELETE]
func DeletePhoto(c *gin.Context) {
	db, _ := database.Connect()
	contentType := helpers.GetContentType(c)
	Photo := entity.Photo{}

	//get parameter
	photoID, _ := strconv.Atoi(c.Param("id"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	err := db.Where("id = ?", photoID).Delete(&Photo).Error

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
		Message: "Photo has been deleted successfully",
		Data:    nil,
	})
}
