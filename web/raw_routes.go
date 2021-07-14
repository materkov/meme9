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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/materkov/meme9/web/pb"
	"github.com/materkov/meme9/web/store"
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

	assocType := store.Assoc_VK_ID + strconv.Itoa(body.UserID)
	assocs, err := objectStore.AssocRange(0, assocType, 1)
	if err != nil {
		return "", fmt.Errorf("error selecting by vk id: %w", err)
	}

	userID := 0
	var user *store.User
	if len(assocs) > 0 {
		userID = assocs[0].VkID.ID2
		obj, err := objectStore.ObjGet(userID)
		if err != nil {
			return "", fmt.Errorf("error getting users: %w", err)
		} else if obj == nil || obj.User == nil {
			return "", fmt.Errorf("nil or not user object")
		}

		user = obj.User
	} else {
		userID, err = objectStore.GenerateNextID()
		if err != nil {
			return "", fmt.Errorf("error generating user id: %w", err)
		}

		user = &store.User{
			ID: userID,
		}

		err = objectStore.ObjAdd(&store.StoredObject{ID: userID, User: &store.User{
			ID: userID,
		}})
		if err != nil {
			return "", fmt.Errorf("error saving obj: %w", err)
		}

		assocType := store.Assoc_VK_ID + strconv.Itoa(body.UserID)
		err = objectStore.AssocAdd(0, userID, assocType, &store.StoredAssoc{VkID: &store.VkID{
			ID1:  0,
			ID2:  userID,
			Type: assocType,
		}})
		if err != nil {
			return "", fmt.Errorf("error saving assoc: %w", err)
		}
	}

	vkName, vkAvatar, err := fetchVkData(body.UserID, body.AccessToken)
	if err != nil {
		log.Printf("Error getting vk data: %s", err)
	} else {
		user.Name = vkName
		user.VkAvatar = vkAvatar
		err = objectStore.ObjUpdate(&store.StoredObject{ID: user.ID, User: user})
		if err != nil {
			return "", fmt.Errorf("failed updating name and avatar: %w", err)
		}
	}

	objectID, err := objectStore.GenerateNextID()
	if err != nil {
		return "", fmt.Errorf("failed generating object id: %w", err)
	}

	token := fmt.Sprintf("%d-%s", objectID, RandString(40))
	err = objectStore.ObjAdd(&store.StoredObject{ID: objectID, Token: &store.Token{
		ID:     objectID,
		Token:  token,
		UserID: userID,
	}})
	if err != nil {
		return "", fmt.Errorf("failed saving token: %w", err)
	}

	return token, nil
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

	objectID, err := objectStore.GenerateNextID()
	if err != nil {
		fmt.Fprintf(w, "cannot upload file")
		return
	}

	photo := store.Photo{
		ID:     objectID,
		UserID: viewer.UserID,
		Path:   filePath,
	}

	err = objectStore.ObjAdd(&store.StoredObject{ID: photo.ID, Photo: &photo})
	if err != nil {
		fmt.Fprintf(w, "cannot save photo")
		return
	}

	_, _ = fmt.Fprintf(w, "%d", photo.ID)
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

func handleAPI(w http.ResponseWriter, request *http.Request) {
	method := request.URL.Query().Get("method")
	tracer := NewTracer("api " + method)
	defer tracer.Stop()
	w.Header().Set("x-request-id", fmt.Sprintf("%x", tracer.ID))

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writeAPIError(w, fmt.Errorf("failed reading body"))
		return
	}

	ctx := request.Context()
	viewer := GetViewerFromContext(ctx)
	m := protojson.UnmarshalOptions{DiscardUnknown: true}

	var resp proto.Message

	switch method {
	case "meme.Feed.GetHeader":
		req := &pb.FeedGetHeaderRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = feedSrv.GetHeader(ctx, req)
	case "meme.Posts.Add":
		req := &pb.PostsAddRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = postsSrv.Add(ctx, req)
	case "meme.Posts.ToggleLike":
		req := &pb.ToggleLikeRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = postsSrv.ToggleLike(ctx, req)
	case "meme.Posts.AddComment":
		req := &pb.AddCommentRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = postsSrv.AddComment(ctx, req)
	case "meme.Utils.ResolveRoute":
		req := &pb.ResolveRouteRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = utilsSrv.ResolveRoute(ctx, req)
	case "meme.Relations.Follow":
		req := &pb.RelationsFollowRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = relationsSrv.Follow(ctx, req)
	case "meme.Relations.Unfollow":
		req := &pb.RelationsUnfollowRequest{}
		if err := m.Unmarshal(body, req); err != nil {
			writeAPIError(w, fmt.Errorf("failed unmarshaling request"))
			return
		}
		resp, err = relationsSrv.Unfollow(ctx, req)
	default:
		err = fmt.Errorf("unknown method")
	}

	if err != nil {
		writeAPIError(w, err)
		return
	}

	marshaller := &protojson.MarshalOptions{}
	respBytes, _ := marshaller.Marshal(resp)
	_, _ = w.Write(respBytes)

	defer tracer.StartChild("api log").Stop()
	objectID, err := objectStore.GenerateNextID()
	if err != nil {
		return
	}

	err = objectStore.ObjAdd(&store.StoredObject{ID: objectID, APILog: &store.APILog{
		ID:       objectID,
		UserID:   viewer.UserID,
		Method:   method,
		Request:  string(body),
		Response: string(respBytes),
	}})
	if err != nil {
		log.Printf("Error saving api logs: %s", err)
	}
}
