package pkg

type Post struct {
	ID     int
	UserID int
	Date   int
	Text   string
}

func GetPost(id int) (*Post, error) {
	obj := &Post{}
	err := getObject(id, ObjTypePost, obj)
	if err != nil {
		return nil, err
	}

	obj.ID = id
	return obj, nil
}
