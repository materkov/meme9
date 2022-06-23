package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/materkov/web3/store"
)

func UpdateAvatar(file []byte, userID int) error {
	token, err := StorageAuthorize()
	if err != nil {
		return fmt.Errorf("error getting storage auth token: %w", err)
	}

	hash := sha256.New()
	hash.Write(file)
	fileHash := hex.EncodeToString(hash.Sum(nil))
	fileName := fmt.Sprintf("avatars/%s", fileHash)

	err = StorageUpload(fileName, file, token)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	objects, err := GlobalStore.ObjGet([]int{userID})
	if err != nil {
		return fmt.Errorf("error loading user: %w", err)
	}

	user, ok := objects[userID].(*store.User)
	if !ok {
		return fmt.Errorf("error loading user")
	}

	user.AvatarFile = fileHash

	err = GlobalStore.ObjUpdate(user)
	if err != nil {
		return fmt.Errorf("failed saving user: %w", err)
	}

	return nil
}
