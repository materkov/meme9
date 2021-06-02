package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/materkov/meme9/web/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func handleVKCallback(w http.ResponseWriter, r *http.Request) {
	viewer := GetViewerFromContext(r.Context())
	accessToken, err := doVKCallback(r.URL.Query().Get("code"), viewer)
	if err != nil {
		log.Printf("Error: %s", err)
		_, _ = fmt.Fprint(w, "Failed to authorize")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func doVKCallback(code string, viewer *Viewer) (string, error) {
	if code == "" {
		return "", fmt.Errorf("empty VK code")
	}

	redirectURI := fmt.Sprintf("%s://%s/vk-callback", viewer.RequestScheme, viewer.RequestHost)

	resp, err := http.PostForm("https://oauth.vk.com/access_token", url.Values{
		"client_id":     []string{strconv.Itoa(config.VKAppID)},
		"client_secret": []string{config.VKAppSecret},
		"redirect_uri":  []string{redirectURI},
		"code":          []string{code},
	})
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	body := struct {
		AccessToken string `json:"access_token"`
		UserID      int    `json:"user_id"`
	}{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return "", fmt.Errorf("error parsing json %s: %s", bodyBytes, err)
	} else if body.AccessToken == "" {
		return "", fmt.Errorf("empty access token: %s", bodyBytes)
	}

	userID, err := store.GetByVkID(body.UserID)
	if err != nil {
		return "", fmt.Errorf("error selecting by vk id: %w", err)
	}

	var user *User
	users, err := store.GetUsers([]int{userID})
	if err != nil {
		return "", fmt.Errorf("error getting users: %w", err)
	} else if len(users) == 1 {
		user = users[0]
	} else {
		userID, err = store.GenerateNextID(ObjectTypeUser)
		if err != nil {
			return "", fmt.Errorf("error generating user id: %w", err)
		}

		user = &User{
			ID:   userID,
			VkID: body.UserID,
		}

		err = store.AddUserByVK(user)
		if err != nil {
			return "", fmt.Errorf("error saving user: %w", err)
		}
	}

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err != nil {
		log.Printf("Error getting vk data: %s", err)
	} else {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		err = store.UpdateNameAvatar(user)
		if err != nil {
			return "", fmt.Errorf("failed updating name and avatar: %w", err)
		}
	}

	objectID, err := store.GenerateNextID(ObjectTypeToken)
	if err != nil {
		return "", fmt.Errorf("failed generating object id: %w", err)
	}

	token := Token{
		ID:     objectID,
		Token:  RandString(50),
		UserID: userID,
	}
	err = store.AddToken(&token)
	if err != nil {
		return "", fmt.Errorf("failed saving token: %w", err)
	}

	return token.Token, nil
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	respRoute, _ := utilsSrv.ResolveRoute(r.Context(), &pb.ResolveRouteRequest{Url: r.URL.Path})
	resp, _ := feedSrv.GetHeader(r.Context(), nil)

	initialDataHeader, _ := protojson.Marshal(resp)
	initialData, _ := protojson.Marshal(respRoute)

	_, _ = fmt.Fprintf(w, `
<!DOCTYPE html>
<html lang="ru">
<head>
    <title>meme</title>
    <meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<script>
    window.initialDataHeader = %s;
    window.initialData = %s;
</script>
<div id="root"></div>
<script src="/static/App.js"></script>
</body>
</html>
`,
		initialDataHeader, initialData,
	)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	viewer := GetViewerFromContext(r.Context())
	if viewer.UserID == 0 {
		fmt.Fprintf(w, "no auth")
		return
	}

	file, err := ioutil.ReadAll(r.Body)
	if err != nil || len(file) == 0 {
		fmt.Fprintf(w, "no file")
		return
	}

	filePath := RandString(20)

	_, err = s3manager.NewUploader(awsSession).Upload(&s3manager.UploadInput{
		Bucket:      aws.String("meme-files"),
		Key:         aws.String("photos/" + filePath + ".jpg"),
		Body:        bytes.NewReader(file),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		fmt.Fprintf(w, "cannot upload file")
		return
	}

	objectID, err := store.GenerateNextID(ObjectTypePhoto)
	if err != nil {
		fmt.Fprintf(w, "cannot upload file")
		return
	}

	photo := Photo{
		ID:     objectID,
		UserID: viewer.UserID,
		Path:   filePath,
	}
	err = store.AddPhoto(&photo)
	if err != nil {
		fmt.Fprintf(w, "cannot save photo")
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("%d", photo.ID))
}

func writeAPIError(w http.ResponseWriter, err error) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	w.WriteHeader(400)
	_ = json.NewEncoder(w).Encode(response)
}

func handleAPI(w http.ResponseWriter, req *http.Request) {
	method := req.URL.Query().Get("method")
	method = strings.Replace(method, "/", ".", -1)

	body, err := ioutil.ReadAll(req.Body)
	c := req.Context()
	viewer := GetViewerFromContext(c)

	m := protojson.UnmarshalOptions{DiscardUnknown: true}

	var resp proto.Message
	switch method {
	case "meme.Feed.GetHeader":
		r := &pb.FeedGetHeaderRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = feedSrv.GetHeader(c, r)
	case "meme.Profile.Get":
		r := &pb.ProfileGetRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = profileSrv.Get(c, r)
	case "meme.Posts.Add":
		r := &pb.PostsAddRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = postsSrv.Add(c, r)
	case "meme.Posts.ToggleLike":
		r := &pb.ToggleLikeRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = postsSrv.ToggleLike(c, r)
	case "meme.Posts.AddComment":
		r := &pb.AddCommentRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = postsSrv.AddComment(c, r)
	case "meme.Utils.ResolveRoute":
		r := &pb.ResolveRouteRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = utilsSrv.ResolveRoute(c, r)
	case "meme.Relations.Follow":
		r := &pb.RelationsFollowRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = relationsSrv.Follow(c, r)
	case "meme.Relations.Unfollow":
		r := &pb.RelationsUnfollowRequest{}
		err := m.Unmarshal(body, r)
		if err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = relationsSrv.Unfollow(c, r)
	default:
		writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
		return
	}

	if err != nil {
		writeAPIError(w, err)
		return
	} else {
		marshaller := &protojson.MarshalOptions{}
		respBytes, _ := marshaller.Marshal(resp)
		_, _ = w.Write(respBytes)

		_ = store.AddAPILog(viewer.UserID, method, body, respBytes)
	}
}
