package store

var GlobalStore *SqlStore

func GetPost(id int) (*Post, error) {
	obj := &Post{}
	err := GlobalStore.getObject(id, ObjTypePost, obj)
	if err != nil {
		return nil, err
	}
	obj.ID = id
	return obj, nil
}

func GetUser(id int) (*User, error) {
	obj := &User{}
	err := GlobalStore.getObject(id, ObjTypeUser, obj)
	if err != nil {
		return nil, err
	}
	obj.ID = id
	return obj, err
}

func GetConfig() (*Config, error) {
	obj := &Config{}
	// TODO think about 5
	err := GlobalStore.getObject(FakeObjConfig, ObjTypeConfig, obj)
	if err != nil {
		return nil, err
	}
	//obj.ID = id
	return obj, err
}

func GetToken(id int) (*Token, error) {
	obj := &Token{}
	err := GlobalStore.getObject(id, ObjTypeToken, obj)
	if err != nil {
		return nil, err
	}
	obj.ID = id
	return obj, err
}

func AddToken(obj *Token) error {
	id, err := GlobalStore.AddObject(ObjTypeToken, obj)
	if err != nil {
		return err
	}
	obj.ID = id
	return nil
}
