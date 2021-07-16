package app

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/materkov/meme9/web/store"
)

func (a *App) UploadPhoto(file []byte, userID int) (*store.Photo, error) {
	filePath := RandString(20)

	_, err := s3manager.NewUploader(awsSession).Upload(&s3manager.UploadInput{
		Bucket:      aws.String("meme-files"),
		Key:         aws.String("photos/" + filePath + ".jpg"),
		Body:        bytes.NewReader(file),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot upload file: %w", err)
	}

	objectID, err := a.Store.GenerateNextID()
	if err != nil {
		return nil, fmt.Errorf("cannot generate id: %w", err)
	}

	photo := store.Photo{
		ID:     objectID,
		UserID: userID,
		Path:   filePath,
	}

	err = a.Store.ObjAdd(&store.StoredObject{ID: photo.ID, Photo: &photo})
	if err != nil {
		return nil, fmt.Errorf("error saving object: %w",err)
	}

	return &photo, nil
}
