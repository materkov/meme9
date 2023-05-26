package pkg

type User struct {
	ID   int
	Name string
}

func GetUser(id int) (*User, error) {
	obj := &User{}
	err := getObject(id, objTypeUser, obj)
	return obj, err
}
