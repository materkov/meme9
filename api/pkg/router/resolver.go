package router

import (
	"regexp"

	"github.com/gogo/protobuf/proto"
	"github.com/materkov/meme9/api/pb"
)

var GlobalJs = []string{
	"/static/React.js",
	"/static/Global.js",
}

func ResolveRoute(url string) ResolvedRoute {
	if match, _ := regexp.MatchString(`^/users/([0-9]+)`, url); match {
		return ResolvedRoute{
			Js: []string{
				"/static/UserPage.js",
			},
			RootComponent: "UserPage",
			ApiMethod:     "meme.API.UserPage",
			ApiArgs: &pb.UserPageRequest{
				UserId: url[7:],
			},
		}
	}

	if match, _ := regexp.MatchString(`^/posts/([0-9]+)`, url); match {
		return ResolvedRoute{
			Js: []string{
				"/static/PostPage.js",
			},
			RootComponent: "PostPage",
			ApiMethod:     "meme.API.PostPage",
			ApiArgs: &pb.PostPageRequest{
				PostId: url[7:],
			},
		}
	}

	if match, _ := regexp.MatchString(`^/login`, url); match {
		return ResolvedRoute{
			Js: []string{
				"/static/LoginPage.js",
			},
			RootComponent: "LoginPage",
			ApiMethod:     "meme.API.LoginPage",
			ApiArgs:       &pb.LoginPageRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/composer`, url); match {
		return ResolvedRoute{
			Js: []string{
				"/static/Composer.js",
			},
			RootComponent: "Composer",
			ApiMethod:     "meme.API.Composer",
			ApiArgs:       &pb.ComposerRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/feed`, url); match {
		return ResolvedRoute{
			Js: []string{
				"/static/Feed.js",
			},
			RootComponent: "Feed",
			ApiMethod:     "meme.API.GetFeed",
			ApiArgs:       &pb.GetFeedRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/vk-callback`, url); match {
		return ResolvedRoute{
			Js:            []string{},
			RootComponent: "",
			ApiMethod:     "meme.API.VKCallback",
			ApiArgs:       &pb.VKCallbackRequest{},
		}
	}

	if match, _ := regexp.MatchString(`^/$`, url); match {
		return ResolvedRoute{
			Js: []string{
				"/static/Index.js",
			},
			RootComponent: "Index",
			ApiMethod:     "meme.API.Index",
			ApiArgs:       &pb.IndexRequest{},
		}
	}

	return ResolvedRoute{}
}

type ResolvedRoute struct {
	Js            []string
	RootComponent string
	ApiMethod     string
	ApiArgs       proto.Message
}
