package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/materkov/meme9/rss2/pb/github.com/materkov/meme9/api"
	"github.com/twitchtv/twirp"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Config struct {
	Login    string
	Password string
}

var postsClient = api.NewPostsProtobufClient("http://localhost:8002", http.DefaultClient)
var authClient = api.NewAuthProtobufClient("http://localhost:8002", http.DefaultClient)
var config Config
var authToken string
var authId string

func main() {
	err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = auth()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized, user id: %s", authId)

	for {
		stats, err := doRequest()
		if err != nil {
			log.Printf("Error: %s", err)
		} else {
			log.Printf("Result: %+v", stats)
		}

		time.Sleep(time.Second * 60)
	}
}

func parseConfig() error {
	file, err := os.ReadFile("/Users/m.materkov/projects/meme9/configs/imgproxy.json")
	if err != nil {
		file, err = os.ReadFile("/apps/meme9-config/imgproxy.json")
	}
	if err != nil {
		return fmt.Errorf("config not found")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	return nil
}

func auth() error {
	resp, err := authClient.Login(context.Background(), &api.EmailReq{
		Email:    config.Login,
		Password: config.Password,
	})
	if err != nil {
		return err
	}

	authToken = resp.Token
	authId = resp.UserId
	return nil
}

type parseResult struct {
	Skipped []string
	Added   []string
}

func doRequest() (parseResult, error) {
	result := parseResult{}

	resp, err := postsClient.List(context.Background(), &api.ListReq{ByUserId: authId})
	if err != nil {
		return result, fmt.Errorf("list error: %w", err)
	}

	parsedUrls := map[string]bool{}
	urlRegex := regexp.MustCompile("(https?://[^\\s]+)")
	for _, post := range resp.Items {
		foundURL := urlRegex.FindString(post.Text)
		if foundURL != "" {
			parsedUrls[foundURL] = true
		}
	}

	rssText, err := fetchRss()
	if err != nil {
		return result, fmt.Errorf("error getting rss: %w", err)
	}

	rss := rss{}
	err = xml.Unmarshal([]byte(rssText), &rss)
	if err != nil {
		return result, fmt.Errorf("error unmarshaling rss: %w", err)
	}

	if len(rss.Channel.Items) > 5 {
		rss.Channel.Items = rss.Channel.Items[:5]
	}

	ctx, _ := twirp.WithHTTPRequestHeaders(context.Background(), http.Header{
		"authorization": []string{"Bearer " + authToken},
	})

	for _, item := range rss.Channel.Items {
		if len(item.Link) != 2 {
			return result, fmt.Errorf("rss item should contain 2 links")
		}
		if item.Link[0].URL == "" {
			return result, fmt.Errorf("rss first link is empty")
		}

		url := item.Link[0].URL
		if parsedUrls[url] {
			result.Skipped = append(result.Skipped, url)
			continue
		}

		text := item.Title + ".\n" + item.Description + "\n\n" + url

		_, err := postsClient.Add(ctx, &api.AddReq{
			Text: text,
		})
		if err != nil {
			return result, fmt.Errorf("error adding post: %w", err)
		} else {
			result.Added = append(result.Added, url)
		}
	}

	return result, nil
}

type rss struct {
	Channel channel `xml:"channel"`
}

type channel struct {
	Items []item `xml:"item"`
}

type item struct {
	Link        []link `xml:"link"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

type link struct {
	URL  string `xml:",innerxml"`
	Href string `xml:"href,attr"`
}

func fetchRss() (string, error) {
	resp, err := http.DefaultClient.Get("https://rss.nytimes.com/services/xml/rss/nyt/World.xml")
	if err != nil {
		return "", fmt.Errorf("error doing http request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("bad http code: %d, %s", resp.StatusCode, body)
	}

	return string(body), err
}
