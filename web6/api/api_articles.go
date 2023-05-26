package api

import (
	"fmt"
	"github.com/materkov/meme9/web6/pkg"
	"log"
	"strconv"
)

//go:generate easyjson -all -lower_camel_case -omit_empty api_articles.go

type Article struct {
	ID    string
	Title string

	Paragraphs []*Paragraph
}

type Paragraph struct {
	Text  *ParagraphText
	Image *ParagraphImage
}

type ParagraphText struct {
	ID   string
	Text string
}

type ParagraphImage struct {
	ID  string
	URL string
}

type articlesListReq struct {
	ID string
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

	wrappedArticle := &Article{
		ID:    strconv.Itoa(article.ID),
		Title: article.Title,
	}

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
		}
	}

	return wrappedArticle, nil
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
