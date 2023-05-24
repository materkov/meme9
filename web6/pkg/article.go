package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Article struct {
	ID    int
	Title string

	Paragraphs []Paragraph
}

const (
	objTypeArticle = 1
)

const (
	paragraphText  = 1
	paragraphImage = 2
)

type Paragraph struct {
	ID int

	ParagraphText  *ParagraphText
	ParagraphImage *ParagraphImage
}

type ParagraphText struct {
	Text string
}

type ParagraphImage struct {
	URL string
}

var SqlClient *sql.DB

var ErrObjectNotFound = fmt.Errorf("object not found")

func GetArticle(id int) (*Article, error) {
	var data []byte
	err := SqlClient.QueryRow("select data from objects where id = ? and obj_type = 1", id).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error selecting database: %w", err)
	}

	article := Article{}
	err = json.Unmarshal(data, &article)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %w", err)
	}

	return &article, nil
}

func GetConfig() error {
	var data []byte
	err := SqlClient.QueryRow("select data from objects where id = 5").Scan(&data)
	if err != nil {
		return err
	}

	_ = json.Unmarshal(data, &GlobalConfig)
	return nil
}

func SaveArticle(article *Article) error {
	data, _ := json.Marshal(article)
	_, err := SqlClient.Exec("update objects set data = ? where id = ?", data, article.ID)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}
	return nil
}
