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

func GetToken(id int) (*Token, error) {
	obj := &Token{}
	err := getObject(id, ObjTypeToken, obj)
	if err != nil {
		return nil, err
	}
	obj.ID = id
	return obj, err
}

func AddToken(obj *Token) error {
	id, err := AddObject(ObjTypeToken, obj)
	if err != nil {
		return err
	}
	obj.ID = id
	return nil
}
