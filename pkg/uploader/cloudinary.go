package uploader

import (
	"context"
	"fmt"
	"mime/multipart"
	"shifty-backend/pkg/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, subFolder string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
	GetPublicIDFromURL(url string) string
}
type CloudinaryService struct {
	client     *cloudinary.Cloudinary
	baseFolder string
}

func NewCloudinary(cloudName, apiKey, apiSecret, baseFolder string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{
		client:     cld,
		baseFolder: baseFolder,
	}, nil
}

// Func Upload Image
func (c *CloudinaryService) UploadImage(ctx context.Context, file *multipart.FileHeader, subFolder string) (string, error) {

	// Open the file
	src, err := file.Open()

	// return error
	if err != nil {
		return "", err
	}

	// close file when exist this func
	defer src.Close()
	fullFolder := c.baseFolder
	if subFolder != "" {
		fullFolder = fmt.Sprintf("%s/%s", c.baseFolder, subFolder)
	}
	// Upload file into folder
	resp, err := c.client.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: fullFolder,
	})

	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

// Func Delete Image
func (c *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {

	// Destroy url in cloudinary with publicID
	_, err := c.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	return err
}

// GetPublicIDFromURL extracts the Public ID from a Cloudinary URL.
func (c *CloudinaryService) GetPublicIDFromURL(url string) string {
	// Regex matches the string after "/upload/" (and optional version like v123/)
	// up to the first dot (file extension).

	match := utils.CloudinaryPublicIDRegex.FindStringSubmatch(url)

	if len(match) > 1 {
		return match[1]
	}
	return ""
}
