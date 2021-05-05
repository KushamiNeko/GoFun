package database

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	testDbRoot string = filepath.Join(os.Getenv("HOME"), "Documents/database/testing/json")
)

const (
	dbAdmin        = "admin"
	dbTradingBooks = "trading_books"
	dbLiveTrading  = "live_trading"
	dbPaperTrading = "paper_trading"
	dbWatchList    = "watch_list"

	colUser = "user"

	testDb  = "test"
	testCol = "test"
)

func TestInsert(t *testing.T) {

	db := NewFileDB(testDbRoot, JsonDB)
	e := map[string]string{
		"a": "b",
	}
	err := db.Insert(testDb, testCol, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	db.DropCol(testDb, testCol)
}

func TestReplace(t *testing.T) {
	db := NewFileDB(testDbRoot, JsonDB)
	e := map[string]string{
		"a": "b",
		"b": "b",
	}
	err := db.Insert(testDb, testCol, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	q := map[string]string{
		"a": "b",
	}

	e = map[string]string{
		"a": "c",
	}

	err = db.Replace(testDb, testCol, q, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	n, err := db.Find(testDb, testCol, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if len(n) != 1 {
		t.Errorf("len should be 1 but get %d", len(n))
	}

	if n[0]["a"] != "c" {
		t.Errorf("new value should be c, but get %s", n[0]["a"])
	}

	if n[0]["b"] != "b" {
		t.Errorf("new value should be b, but get %s", n[0]["b"])
	}

	db.DropCol(testDb, testCol)
}

func TestDelete(t *testing.T) {
	db := NewFileDB(testDbRoot, JsonDB)
	es := []map[string]string{
		{
			"a": "b",
			"b": "b",
		},
		{
			"c": "d",
			"b": "b",
		},
	}
	err := db.Insert(testDb, testCol, es...)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	q := map[string]string{
		"b": "b",
	}

	rs, err := db.Find(testDb, testCol, q)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if len(rs) != 2 {
		t.Errorf("len should be 2 but get %d", len(rs))
	}

	dq := map[string]string{
		"a": "b",
	}

	err = db.Delete(testDb, testCol, dq)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	rs, err = db.Find(testDb, testCol, q)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if len(rs) != 1 {
		t.Errorf("len should be 1 but get %d", len(rs))
	}

	db.DropCol(testDb, testCol)
}
