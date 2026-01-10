package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/materkov/meme9/photos-service/api/mocks"
	"github.com/materkov/meme9/photos-service/processor"
)

func TestService_HandleUpload(t *testing.T) {
	t.Run("successful upload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProcessor := mocks.NewMockProcessor(ctrl)
		mockUploader := mocks.NewMockUploader(ctrl)
		service := New(mockProcessor, mockUploader)

		originalImage := []byte("original image data")
		processedImage := []byte("processed image data")
		uploadURL := "https://example.com/image.jpg"

		mockProcessor.EXPECT().
			Process(gomock.Any(), originalImage).
			Return(processedImage, nil)

		mockUploader.EXPECT().
			Upload(gomock.Any(), processedImage).
			Return(uploadURL, nil)

		req := httptest.NewRequest(http.MethodPost, "/twirp/meme.photos.Photos/upload", bytes.NewReader(originalImage))
		ctx := context.WithValue(req.Context(), UserIDKey, "user123")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		service.HandleUpload(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, uploadURL, w.Body.String())
	})

	t.Run("no auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProcessor := mocks.NewMockProcessor(ctrl)
		mockUploader := mocks.NewMockUploader(ctrl)
		service := New(mockProcessor, mockUploader)

		req := httptest.NewRequest(http.MethodPost, "/twirp/meme.photos.Photos/upload", bytes.NewReader([]byte("image data")))
		w := httptest.NewRecorder()

		service.HandleUpload(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		require.Equal(t, "auth_required\n", w.Body.String())
	})

	t.Run("invalid image", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProcessor := mocks.NewMockProcessor(ctrl)
		mockUploader := mocks.NewMockUploader(ctrl)
		service := New(mockProcessor, mockUploader)

		originalImage := []byte("invalid image data")

		mockProcessor.EXPECT().
			Process(gomock.Any(), originalImage).
			Return(nil, processor.ErrInvalidImage)

		req := httptest.NewRequest(http.MethodPost, "/twirp/meme.photos.Photos/upload", bytes.NewReader(originalImage))
		ctx := context.WithValue(req.Context(), UserIDKey, "user123")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		service.HandleUpload(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "Invalid image\n", w.Body.String())
	})

	t.Run("processing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProcessor := mocks.NewMockProcessor(ctrl)
		mockUploader := mocks.NewMockUploader(ctrl)
		service := New(mockProcessor, mockUploader)

		originalImage := []byte("image data")

		mockProcessor.EXPECT().
			Process(gomock.Any(), originalImage).
			Return(nil, errors.New("processing error"))

		req := httptest.NewRequest(http.MethodPost, "/twirp/meme.photos.Photos/upload", bytes.NewReader(originalImage))
		ctx := context.WithValue(req.Context(), UserIDKey, "user123")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		service.HandleUpload(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Equal(t, "error\n", w.Body.String())
	})

	t.Run("upload error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProcessor := mocks.NewMockProcessor(ctrl)
		mockUploader := mocks.NewMockUploader(ctrl)
		service := New(mockProcessor, mockUploader)

		originalImage := []byte("image data")
		processedImage := []byte("processed image data")

		mockProcessor.EXPECT().
			Process(gomock.Any(), originalImage).
			Return(processedImage, nil)

		mockUploader.EXPECT().
			Upload(gomock.Any(), processedImage).
			Return("", errors.New("upload error"))

		req := httptest.NewRequest(http.MethodPost, "/twirp/meme.photos.Photos/upload", bytes.NewReader(originalImage))
		ctx := context.WithValue(req.Context(), UserIDKey, "user123")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		service.HandleUpload(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Equal(t, "error\n", w.Body.String())
	})
}
