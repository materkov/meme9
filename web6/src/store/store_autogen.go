package store

func GetPost(id int) (*Post, error) {
	obj := &Post{}
	err := getObject(id, ObjTypePost, obj)
	if err != nil {
		return nil, err
	}
	obj.ID = id
	return obj, nil
}

func GetUser(id int) (*User, error) {
	obj := &User{}
	err := getObject(id, ObjTypeUser, obj)
	if err != nil {
		return nil, err
	}
	obj.ID = id
	return obj, err
}

func GetConfig() (*Config, error) {
	obj := &Config{}
	err := getObject(5, ObjTypeConfig, obj)
	if err != nil {
		return nil, err
	}
	//obj.ID = id
	return obj, err
}
