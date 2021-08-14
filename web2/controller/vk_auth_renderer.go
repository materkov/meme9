package controller

import (
	"fmt"
	"github.com/materkov/meme9/web2/lib"
	"net/url"
)

type VkAuthRenderer struct {
	URL string
}

func (v *VkAuthRenderer) Render() string {
	requestScheme := lib.DefaultConfig.RequestScheme
	requestHost := lib.DefaultConfig.RequestHost
	vkAppID := lib.DefaultConfig.VkAppID
	redirectURL := fmt.Sprintf("%s://%s/vk-callback", requestScheme, requestHost)
	redirectURL = url.QueryEscape(redirectURL)
	vkURL := fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%d&response_type=code&redirect_uri=%s", vkAppID, redirectURL)

	return fmt.Sprintf("<a href=\"%s\">Авторизоваться через ВК</a>", vkURL)
}
