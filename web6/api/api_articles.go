package api

import (
	"fmt"
	"github.com/materkov/meme9/web6/pkg"
	"log"
	"strconv"
	"time"
)

//go:generate easyjson -all -lower_camel_case -omit_empty api_articles.go

type Article struct {
	ID    string
	Title string
	User  *User

	CreatedAt string

	Paragraphs []*Paragraph
}

type Paragraph struct {
	Text  *ParagraphText
	Image *ParagraphImage
	List  *ParagraphList
}

type ParagraphText struct {
	ID   string
	Text string
}

type ParagraphImage struct {
	ID  string
	URL string
}

type ParagraphList struct {
	ID    string
	Items []string
	Type  ListType
}

type ListType string

const (
	ListTypeUnknown   ListType = ""
	ListTypeOrdered   ListType = "ORDERED"
	ListTypeUnordered ListType = "UNORDERED"
)

type articlesListReq struct {
	ID string
}

func transformArticle(articleId string, article *pkg.Article) *Article {
	wrappedArticle := &Article{
		ID: articleId,
	}
	if article == nil {
		return wrappedArticle
	}

	wrappedArticle.Title = article.Title
	wrappedArticle.CreatedAt = transformDate(article.Date)

	user, err := pkg.GetUser(article.UserID)
	if err != nil {
		log.Printf("[ERROR] Error loading user: %s", err)
	}

	wrappedArticle.User = transformUser(article.UserID, user)

	wrappedArticle.Paragraphs = make([]*Paragraph, len(article.Paragraphs))
	for i, p := range article.Paragraphs {
		if p.ParagraphText != nil {
			wrappedArticle.Paragraphs[i] = &Paragraph{Text: &ParagraphText{
				ID:   strconv.Itoa(p.ID),
				Text: p.ParagraphText.Text,
			}}
		} else if p.ParagraphImage != nil {
			wrappedArticle.Paragraphs[i] = &Paragraph{Image: &ParagraphImage{
				ID:  strconv.Itoa(p.ID),
				URL: p.ParagraphImage.URL,
			}}
		} else if p.ParagraphList != nil {
			wrappedArticle.Paragraphs[i] = &Paragraph{List: &ParagraphList{
				ID:    strconv.Itoa(p.ID),
				Type:  transformListType(p.ParagraphList.Type),
				Items: p.ParagraphList.Items,
			}}
		}
	}

	return wrappedArticle
}

func transformListType(t pkg.ListType) ListType {
	if t == pkg.ListTypeOrdered {
		return ListTypeOrdered
	} else if t == pkg.ListTypeUnordered {
		return ListTypeUnordered
	} else {
		return ListTypeUnknown
	}
}

func (a *API) ArticlesList(r *articlesListReq) (*Article, error) {
	id, _ := strconv.Atoi(r.ID)
	log.Printf("Article %d", id)
	article, err := pkg.GetArticle(id)
	if err == pkg.ErrObjectNotFound {
		return nil, &Error{Code: 404, Message: "article not found"}
	} else if err != nil {
		return nil, err
	}

	wrappedArticle := transformArticle(r.ID, article)
	return wrappedArticle, nil
}

func (a *API) ArticlesLastPosted(r *Void) ([]*Article, error) {
	articleIds, err := pkg.GetEdges(pkg.FakeObjPosted, pkg.EdgeTypeLastPosted)
	if err != nil {
		log.Printf("[ERROR] Error getting last posted: %s", err)
	}

	result := make([]*Article, len(articleIds))
	for i, id := range articleIds {
		article, err := pkg.GetArticle(id)
		if err != nil {
			log.Printf("[ERROR] Error getting article: %s", err)
		}

		result[i] = transformArticle(strconv.Itoa(id), article)
	}

	return result, err
}

type InputArticle struct {
	ID         string
	Title      string
	Paragraphs []*InputParagraph
}

type InputParagraph struct {
	InputParagraphText  *InputParagraphText
	InputParagraphImage *InputParagraphImage
}

type InputParagraphText struct {
	Text string
}

type InputParagraphImage struct {
	URL string
}

type Void struct{}

func (a *API) ArticlesSave(r *InputArticle) (*Void, error) {
	article := pkg.Article{
		ID:         0,
		Title:      "",
		Paragraphs: nil,
	}
	id, _ := strconv.Atoi(r.ID)
	if id <= 0 {
		return nil, &Error{Code: 404, Message: "article not found"}
	}

	article.ID = id
	article.Title = r.Title
	article.UpdatedAt = int(time.Now().Unix())

	paragraphID := 1
	for _, paragraph := range r.Paragraphs {
		if paragraph.InputParagraphText != nil {
			article.Paragraphs = append(article.Paragraphs, pkg.Paragraph{
				ID: paragraphID,
				ParagraphText: &pkg.ParagraphText{
					Text: paragraph.InputParagraphText.Text,
				},
			})
		} else if paragraph.InputParagraphImage != nil {
			article.Paragraphs = append(article.Paragraphs, pkg.Paragraph{
				ID: paragraphID,
				ParagraphImage: &pkg.ParagraphImage{
					URL: paragraph.InputParagraphImage.URL,
				},
			})
		} else {
			return nil, fmt.Errorf("incorrect paragraph type")
		}

		paragraphID++
	}

	err := pkg.UpdateObject(&article, article.ID)
	if err != nil {
		return nil, err
	}

	return &Void{}, err
}

type ListPostedByUserReq struct {
	UserId string
}

func (a *API) listPostedByUser(r *ListPostedByUserReq) ([]*Article, error) {
	userID, _ := strconv.Atoi(r.UserId)
	articleIds, err := pkg.GetEdges(userID, pkg.EdgeTypePosted)
	if err != nil {
		return nil, err
	}

	var result []*Article
	for _, articleId := range articleIds {
		article, err := pkg.GetArticle(articleId)
		if err != nil {
			continue
		}
		wrappedArticle := transformArticle(strconv.Itoa(articleId), article)
		result = append(result, wrappedArticle)
	}

	return result, nil
}

type UsersListReq struct {
	UserIds []string
}

func (a *API) usersList(r *UsersListReq) ([]*User, error) {
	result := make([]*User, len(r.UserIds))
	for i, userIdStr := range r.UserIds {
		userId, _ := strconv.Atoi(userIdStr)
		user, err := pkg.GetUser(userId)
		if err != nil {
			log.Printf("[ERROR] Error loading user: %s", err)
		}

		result[i] = transformUser(userId, user)
	}

	return result, nil
}
