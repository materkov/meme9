package sqlmock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type DummyDriver struct{}

func (d *DummyDriver) Open(_ string) (driver.Conn, error) {
	return &DummyConn{
		uniques: map[string]int{},
	}, nil
}

func init() {
	sql.Register("sqlmock", &DummyDriver{})
}

type DummyConn struct {
	idx     int
	objects [][]string
	edges   [][]string
	uniques map[string]int
}

func (d *DummyConn) Close() error {
	return nil
}

func (d *DummyConn) Begin() (driver.Tx, error) {
	panic("implement me")
}

type mockStmt struct {
	query string
	conn  *DummyConn
}

func (m *mockStmt) Close() error {
	return nil
}

func (m *mockStmt) NumInput() int {
	return strings.Count(m.query, "?")
}

type result struct {
	id int
}

func (r *result) LastInsertId() (int64, error) {
	return int64(r.id), nil
}

func (r *result) RowsAffected() (int64, error) {
	panic("implement me")
}

func (m *mockStmt) Exec(args []driver.Value) (driver.Result, error) {
	if m.query == "insert into objects(obj_type, data) values (?, ?)" {
		m.conn.idx++
		id := m.conn.idx
		m.conn.objects = append(m.conn.objects, []string{
			strconv.Itoa(id),
			strconv.Itoa(int(args[0].(int64))),
			string(args[1].([]byte)),
		})

		return &result{id: id}, nil
	} else if m.query == "insert into edges(from_id, to_id, edge_type, date) values (?, ?, ?, ?)" {
		m.conn.idx++
		id := m.conn.idx

		m.conn.edges = append(m.conn.edges, []string{
			strconv.Itoa(id),
			strconv.Itoa(int(args[0].(int64))),
			strconv.Itoa(int(args[1].(int64))),
			strconv.Itoa(int(args[2].(int64))),
			strconv.Itoa(int(time.Now().Unix())),
		})
		return &result{id: id}, nil
	} else if m.query == "update objects set data = ? where id = ?" {
		for idx, obj := range m.conn.objects {
			if obj[0] == strconv.Itoa(int(args[1].(int64))) {
				m.conn.objects[idx][2] = string(args[0].([]byte))
			}
		}
		return &result{id: 0}, nil
	} else if m.query == "delete from edges where from_id = ? and edge_type = ? and to_id = ?" {
		for idx, obj := range m.conn.edges {
			if obj[1] == strconv.Itoa(int(args[0].(int64))) &&
				obj[2] == strconv.Itoa(int(args[2].(int64))) &&
				obj[3] == strconv.Itoa(int(args[1].(int64))) {
				m.conn.edges = append(m.conn.edges[:idx], m.conn.edges[idx+1:]...)
			}
		}
		return &result{id: 0}, nil
	} else if m.query == "insert into uniques(type, `key`, object_id) values (?, ?, ?)" {
		m.conn.uniques[fmt.Sprintf("%d:%s", args[0].(int64), args[1].(string))] = int(args[2].(int64))
		return &result{id: 0}, nil
	} else {
		panic("implement me")
	}
}

func (m *mockStmt) Query(args []driver.Value) (driver.Rows, error) {
	panic("implement me")
}

func (c *DummyConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{query: query, conn: c}, nil
}
func (c *DummyConn) Rollback() error {
	panic("implement me")
}

func (c *DummyConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if query == "select data from objects where id = ? and obj_type = ?" {
		for _, object := range c.objects {
			if object[0] == strconv.Itoa(int(args[0].(int64))) && object[1] == strconv.Itoa(int(args[1].(int64))) {
				return &results{
					columns: []string{"data"},
					result: [][]string{
						{object[2]},
					},
				}, nil
			}
		}
	} else if query == "select date from edges where from_id = ? and to_id = ? and edge_type = ?" {
		for _, edge := range c.edges {
			if edge[1] == strconv.Itoa(int(args[0].(int64))) &&
				edge[2] == strconv.Itoa(int(args[1].(int64))) &&
				edge[3] == strconv.Itoa(int(args[2].(int64))) {
				return &results{
					columns: []string{"date"},
					result: [][]string{
						{edge[4]},
					},
				}, nil
			}
		}
	} else if query == "select count(*) from edges where from_id = ? and edge_type = ?" {
		cnt := 0
		for _, edge := range c.edges {
			if edge[1] == strconv.Itoa(int(args[0].(int64))) &&
				edge[3] == strconv.Itoa(int(args[1].(int64))) {
				cnt++
			}
		}

		return &results{
			columns: []string{"count(*)"},
			result: [][]string{
				{strconv.Itoa(cnt)},
			},
		}, nil
	} else if query == "select to_id, date from edges where from_id = ? and edge_type = ? order by id desc" {
		var localResults [][]string
		for _, edge := range c.edges {
			if edge[1] == strconv.Itoa(int(args[0].(int64))) &&
				edge[3] == strconv.Itoa(int(args[1].(int64))) {
				localResults = append(localResults, []string{
					edge[2],
					edge[4],
				})
			}
		}

		return &results{
			columns: []string{"to_id", "date"},
			result:  localResults,
		}, nil
	} else if query == "select object_id from uniques where type = ? and `key` = ?" {
		var items [][]string
		objectID, _ := c.uniques[fmt.Sprintf("%d:%s", args[0].(int64), args[1].(string))]
		if objectID != 0 {
			items = [][]string{
				{strconv.Itoa(objectID)},
			}
		}
		return &results{
			columns: []string{"object_id"},
			result:  items,
		}, nil
	} else {
		panic("not implemented")
	}

	return &results{}, nil
}

type results struct {
	columns []string
	result  [][]string
	idx     int
}

func (r *results) Columns() []string {
	return r.columns
}

func (r *results) Close() error {
	return nil
}

func (r *results) Next(dest []driver.Value) error {
	if r.idx >= len(r.result) {
		return io.EOF
	}

	for i := 0; i < len(r.columns); i++ {
		if r.columns[0] == "data" {
			// str columns
			dest[i] = driver.Value([]byte(r.result[r.idx][i]))
		} else {
			intVal, _ := strconv.Atoi(r.result[r.idx][i])
			dest[i] = driver.Value(int64(intVal))
		}
	}

	r.idx++

	return nil
}
