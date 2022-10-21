package files

import "fmt"

func GetURL(fileSha string) string {
	return fmt.Sprintf("https://689809.selcdn.ru/meme-files/avatars/%s", fileSha)
}
