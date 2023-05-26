package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Article struct {
	ID     int
	Title  string
	UserID int
	Date   int

	UpdatedAt int

	Paragraphs []Paragraph
}

const (
	objTypeArticle = 1
	objTypeConfig  = 2
	objTypeUser    = 3
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

func getObject(id int, objType int, obj interface{}) error {
	var data []byte
	err := SqlClient.QueryRow("select data from objects where id = ? and obj_type = ?", id, objType).Scan(&data)
	if err == sql.ErrNoRows {
		return ErrObjectNotFound
	} else if err != nil {
		return fmt.Errorf("error selecting database: %w", err)
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("error unmarshaling object %d: %w", objType, err)
	}

	return nil
}

func GetArticle(id int) (*Article, error) {
	article := &Article{}
	err := getObject(id, objTypeArticle, article)
	return article, err
}

func GetConfig() (*Config, error) {
	obj := &Config{}
	err := getObject(5, objTypeConfig, obj)
	return obj, err
}

func UpdateObject(object interface{}, id int) error {
	data, _ := json.Marshal(object)
	_, err := SqlClient.Exec("update objects set data = ? where id = ?", data, id)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}

	return nil
}
