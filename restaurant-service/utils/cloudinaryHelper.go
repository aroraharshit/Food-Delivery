package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func IntiCloudinary() {
	var err error
	cld, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal(err)
	}
}

func UploadImageToCloudinary(ctx context.Context, file *multipart.FileHeader) (string, error) {
	imgFile, err := file.Open()
	if err != nil {
		return "", err
	}

	defer imgFile.Close()

	fileExtention := strings.ToLower(file.Filename[len(file.Filename)-4:])

	if !(fileExtention == ".jpg" || fileExtention == ".png" || fileExtention == ".jpeg") {
		return "", fmt.Errorf("only .jpg,.jpeg,.png files are allowed")
	}

	uploadResponse, err := cld.Upload.Upload(ctx, imgFile, uploader.UploadParams{})

	if err != nil {
		return "", fmt.Errorf("error uploading image to cloadinary %w", err)
	}

	return uploadResponse.SecureURL, nil
}
