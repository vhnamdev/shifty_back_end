package uploader

import (
	"context"
	"mime/multipart"
	"regexp"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	client     *cloudinary.Cloudinary
	folderName string
}

func NewCloudinary(cloudName, apiKey, apiSecret, folderName string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{
		client:     cld,
		folderName: folderName,
	}, nil
}

// Func Upload Image
func (c *CloudinaryService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {

	// Open the file
	src, err := file.Open()

	// return error
	if err != nil {
		return "", err
	}

	// close file when exist this func
	defer src.Close()

	// Upload file into folder
	resp, err := c.client.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: c.folderName,
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
func GetPublicIDFromURL(url string) string {
	// Regex matches the string after "/upload/" (and optional version like v123/)
	// up to the first dot (file extension).
	re := regexp.MustCompile(`/upload/(?:v\d+/)?([^.]+)\.`)

	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		return match[1]
	}
	return ""
}
