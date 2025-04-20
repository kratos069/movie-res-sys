package util

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
}

// Initialize a CloudinaryService instance
func NewCloudinaryService() (*CloudinaryService, error) {
	config, err := LoadConfig("../")
	if err != nil {
		log.Fatalln("cannot load config file:", err)
	}

	cld, err := cloudinary.NewFromParams(config.CloudName, config.CloudApiKey, config.CloudApiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &CloudinaryService{cloudinary: cld}, nil
}

// Upload an image to Cloudinary from a local file path
func (cs *CloudinaryService) UploadImage(ctx context.Context,
	file *multipart.FileHeader) (string, error) {
	// Remove from local
	defer func() {
		os.Remove("uploads/" + file.Filename)
	}()

	// Set the upload parameters (you can customize these as needed)
	uploadParams := uploader.UploadParams{
		Folder: "insta", // specify folder in Cloudinary
	}

	// Perform the upload
	uploadResult, err := cs.cloudinary.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Cloudinary: %w", err)
	}

	// Print the entire upload result for debugging
	fmt.Printf("Upload result: %+v\n", uploadResult)

	// Ensure the SecureURL is part of the result, and then return it
	if uploadResult.SecureURL == "" {
		return "", fmt.Errorf("uploaded image URL is empty")
	}

	// Log the URL of the uploaded image for debugging
	fmt.Printf("Uploaded image URL: %s\n", uploadResult.SecureURL)

	// Return the secure URL of the uploaded image
	return uploadResult.SecureURL, nil
}

// Delete an image from Cloudinary by its Public ID
func (cs *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
    result, err := cs.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
        PublicID: publicID,
    })
    if err != nil {
        return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
    }
    if result.Result == "not found" {
        return fmt.Errorf("image with publicID %s not found", publicID)
    }
    return nil
}
