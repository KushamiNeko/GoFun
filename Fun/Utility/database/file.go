package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type engine uint8

const (
	JsonDB engine = iota
	//YamlDB
)

type FileDB struct {
	root   string
	engine engine
}

func NewFileDB(root string, engine engine) *FileDB {
	return &FileDB{
		root:   root,
		engine: engine,
	}
}

func (j *FileDB) dbPath(db, col string) string {
	switch j.engine {
	case JsonDB:
		fn := fmt.Sprintf("%s_%s.json", db, col)
		return filepath.Join(j.root, fn)
	//case YamlDB:
	//fn := fmt.Sprintf("%s_%s.yaml", db, col)
	//return filepath.Join(j.root, fn)
	default:
		panic("unknown file db engine")
	}
}

func (j *FileDB) read(db, col string) ([]map[string]string, error) {
	r := make([]map[string]string, 0)

	p := j.dbPath(db, col)
	content, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		} else {
			return nil, err
		}
	}

	if len(content) == 0 {
		return r, nil
	}

	switch j.engine {
	case JsonDB:
		err = json.Unmarshal(content, &r)
		if err != nil {
			return nil, err
		}
	//case YamlDB:
	//err = yaml.Unmarshal(content, &r)
	//if err != nil {
	//return nil, err
	//}
	default:
		panic("unknown file db engine")
	}

	return r, nil
}

func (j *FileDB) write(db, col string, entities []map[string]string) error {
	p := j.dbPath(db, col)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	var b []byte
	switch j.engine {
	case JsonDB:
		b, err = json.MarshalIndent(entities, "", "  ")
		if err != nil {
			return err
		}
	//case YamlDB:
	//b, err = yaml.Marshal(entities)
	//if err != nil {
	//return err
	//}
	default:
		panic("unknown file db engine")
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (j *FileDB) Insert(db, col string, entities ...map[string]string) error {

	r, err := j.read(db, col)
	if err != nil {
		return err
	}

	r = append(r, entities...)

	err = j.write(db, col, r)
	if err != nil {
		return err
	}

	return nil
}

func (j *FileDB) Replace(
	db, col string,
	query map[string]string,
	entity map[string]string) error {

	if query == nil || len(query) == 0 {
		return fmt.Errorf("invalid query")
	}

	if len(entity) == 0 {
		return fmt.Errorf("invalid entity")
	}

	r, err := j.read(db, col)
	if err != nil {
		return err
	}

	updated := false

	for i, e := range r {
		found := true
		for k, v := range query {
			if val, ok := e[k]; !(ok && val == v) {
				found = false
			}
		}

		if found {
			for k, v := range entity {
				r[i][k] = v
			}

			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("no entity match query: %v", query)
	}

	err = j.write(db, col, r)
	if err != nil {
		return err
	}

	return nil
}

func (j *FileDB) Find(
	db, col string,
	query map[string]string) ([]map[string]string, error) {

	r, err := j.read(db, col)
	if err != nil {
		return nil, err
	}

	if query == nil || len(query) == 0 {
		return r, nil
	}

	n := make([]map[string]string, 0)

	for i, e := range r {
		found := true
		for k, v := range query {
			if val, ok := e[k]; !(ok && val == v) {
				found = false
			}
		}

		if found {
			n = append(n, r[i])
		}
	}

	return n, nil
}

func (j *FileDB) Delete(
	db, col string,
	query map[string]string) error {

	if query == nil || len(query) == 0 {
		return fmt.Errorf("invalid query or entity")
	}

	r, err := j.read(db, col)
	if err != nil {
		return err
	}

	n := make([]map[string]string, 0, len(r)-1)

	for _, e := range r {
		found := true
		for k, v := range query {
			if val, ok := e[k]; !(ok && val == v) {
				found = false
			}
		}

		if found {
			continue
		} else {
			n = append(n, e)
		}
	}

	err = j.write(db, col, n)
	if err != nil {
		return err
	}

	return nil
}

func (j *FileDB) DropCol(db, col string) error {
	p := j.dbPath(db, col)

	err := os.Remove(p)
	if err != nil {
		return err
	}

	return nil
}
