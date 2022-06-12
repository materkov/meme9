package store

import (
	"encoding/json"
	"fmt"
)

type storeItem struct {
	id      int
	objType int
	objData []byte
}

type listItem struct {
	Object1, Object2 int
	ListType         int
}

type MockStore struct {
	Objects []storeItem
	Lists   []listItem
}

func (m *MockStore) ObjGet(ids []int) (map[int]Object, error) {
	result := map[int]Object{}

	for _, id := range ids {
		for _, object := range m.Objects {
			if object.id == id {
				result[id], _ = parseObject(object.id, object.objType, object.objData)
			}
		}
	}

	return result, nil
}

func (m *MockStore) ListGet(objectID int, listType int) ([]int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStore) ObjAdd(objectID int, objectType int, obj interface{}) error {
	for _, object := range m.Objects {
		if object.id == objectID {
			return fmt.Errorf("id %d already exists", objectID)
		}
	}

	objData, _ := json.Marshal(obj)

	m.Objects = append(m.Objects, storeItem{
		id:      objectID,
		objType: objectType,
		objData: objData,
	})

	return nil
}

func (m *MockStore) ListAdd(object1, listType, object2 int) error {
	for _, item := range m.Lists {
		if item.Object1 == object1 && item.ListType == listType {
			return fmt.Errorf("assoc %d->%d already exists", object1, listType)
		}
	}

	m.Lists = append(m.Lists, listItem{
		Object1:  object1,
		Object2:  object2,
		ListType: listType,
	})

	return nil
}

func (m *MockStore) GetMapping(keyType int, key string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStore) SaveMapping(keyType int, key string, objectID int) error {
	//TODO implement me
	panic("implement me")
}
