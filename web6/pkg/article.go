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
	FakeObjPosted     = 8
	FakeObjPostedPost = -1
	FakeObjVkAuth     = -2

	ObjTypeArticle = 1
	ObjTypeConfig  = 2
	ObjTypeUser    = 3
	ObjTypePost    = 4

	EdgeTypePosted     = 1
	EdgeTypeLastPosted = 2
	EdgeTypePostedPost = 3
	EdgeTypeVkAuth     = 4
)

const (
	paragraphText  = 1
	paragraphImage = 2
)

type Paragraph struct {
	ID int

	ParagraphText  *ParagraphText
	ParagraphImage *ParagraphImage
	ParagraphList  *ParagraphList
}

type ParagraphText struct {
	Text string
}

type ParagraphImage struct {
	URL string
}

type ParagraphList struct {
	Items []string
	Type  ListType
}

type ListType int

const (
	ListTypeUnknown ListType = iota
	ListTypeOrdered
	ListTypeUnordered
)

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
	err := getObject(id, ObjTypeArticle, article)
	return article, err
}

func GetConfig() (*Config, error) {
	obj := &Config{}
	err := getObject(5, ObjTypeConfig, obj)
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

func AddObject(objType int, object interface{}) (int, error) {
	data, _ := json.Marshal(object)
	res, err := SqlClient.Exec("insert into objects(obj_type, data) values (?, ?)", objType, data)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %w", err)
	}

	objId, _ := res.LastInsertId()

	return int(objId), nil
}
