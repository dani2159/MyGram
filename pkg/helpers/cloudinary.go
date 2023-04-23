package helpers

import (
	"context"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
)

var (
	cld *cloudinary.Cloudinary
	err error
)

func InitCloudinary() {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatal("error connecting to cloudinary:", err.Error())
	}

	log.Println("cloudinary connected successfully")
}

func publicIdPath(fileName string) string {
	return "photos/" + fileName
}

func getFileNameFromUrl(fileUrlString string) string {
	fileUrl, err := url.Parse(fileUrlString)
	if err != nil {
		log.Fatal(err)
	}

	// return last path from url
	// example url: https://res.cloudinary.com/xxx/image/upload/yyy/photos/file-name.png"
	// return value: file-name.png
	fileNameWithExt := path.Base(fileUrl.Path)

	// remove extension to get the file-name
	// from 'file-name.png' to 'file-name'
	return strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))
}

func UploadToCloudinary(file multipart.File) (string, error) {
	ctx := context.Background()
	fileName := uuid.New()

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: publicIdPath(fileName.String()), // folder-name/file-name
	})
	if err != nil {
		log.Printf("error uploading file to cloudinary: %v", err.Error())
	}

	return resp.SecureURL, err
}

func DestroyFromCloudinary(fileUrlString string) (err error) {
	ctx := context.Background()

	fileName := getFileNameFromUrl(fileUrlString)

	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicIdPath(fileName),
	})
	if err != nil {
		log.Printf("error trying to delete file '%v' from cloudinary: %v", fileUrlString, err.Error())
	}
	return
}
