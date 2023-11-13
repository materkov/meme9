package store

import (
	"encoding/json"
	"fmt"
	"time"
)

type mockObject struct {
	ID   int
	Type int
	Data []byte
}

type mockEdge struct {
	From, To int
	Type     int
	Date     int
}

type MockStore struct {
	objects   []mockObject
	edges     []mockEdge
	idCounter int
}

func (m *MockStore) getObject(id int, objType int, obj interface{}) error {
	for _, object := range m.objects {
		if object.ID == id && object.Type == objType {
			err := json.Unmarshal(object.Data, obj)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrObjectNotFound
}

func (m *MockStore) AddObject(objType int, object interface{}) (int, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return 0, fmt.Errorf("error marshaling to json: %w", err)
	}

	m.idCounter++

	m.objects = append(m.objects, mockObject{
		ID:   m.idCounter,
		Type: objType,
		Data: data,
	})

	return m.idCounter, nil
}

func (m *MockStore) AddEdge(fromID, toID, edgeType int) error {
	for _, edge := range m.edges {
		if edge.From == fromID && edge.To == toID && edge.Type == edgeType {
			return ErrDuplicateEdge
		}
	}

	m.edges = append(m.edges, mockEdge{
		From: fromID,
		To:   toID,
		Type: edgeType,
		Date: int(time.Now().Unix()),
	})

	return nil
}

func (m *MockStore) GetEdge(fromID, toID, edgeType int) (*Edge, error) {
	for _, edge := range m.edges {
		if edge.From == fromID && edge.To == toID && edge.Type == edgeType {
			return &Edge{
				FromID: edge.From,
				ToID:   edge.To,
				Date:   edge.Date,
			}, nil
		}
	}

	return nil, ErrNoEdge
}

func (m *MockStore) CountEdges(fromID, edgeType int) (int, error) {
	count := 0
	for _, edge := range m.edges {
		if edge.From == fromID && edge.Type == edgeType {
			count++
		}
	}

	return count, nil
}

func (m *MockStore) GetEdges(fromID int, edgeType int) ([]Edge, error) {
	var result []Edge
	for _, edge := range m.edges {
		if edge.From == fromID && edge.Type == edgeType {
			result = append(result, Edge{
				FromID: fromID,
				ToID:   edge.To,
				Date:   edge.Date,
			})
		}
	}

	return result, nil
}

func (m *MockStore) DelEdge(fromID, toID, edgeType int) error {
	for i, edge := range m.edges {
		if edge.From == fromID && edge.To == toID && edge.Type == edgeType {
			m.edges = append(m.edges[:i], m.edges[i+1:]...)
			break
		}
	}

	return nil
}

func (m *MockStore) UpdateObject(object interface{}, id int) error {
	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %w", err)
	}

	for i, obj := range m.objects {
		if obj.ID == id {
			m.objects[i].Data = data
		}
	}

	return nil
}
