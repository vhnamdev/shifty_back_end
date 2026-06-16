package unit_test

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// ===== MOCK UPLOADER =====
type MockUploader struct{ mock.Mock }

func (m *MockUploader) UploadImage(ctx context.Context, file *multipart.FileHeader, subFolder string) (string, error) {
	args := m.Called(ctx, file, subFolder)
	return args.String(0), args.Error(1)
}

func (m *MockUploader) DeleteImage(ctx context.Context, publicID string) error {
	return m.Called(ctx, publicID).Error(0)
}

func (m *MockUploader) GetPublicIDFromURL(url string) string {
	return m.Called(url).String(0)
}
