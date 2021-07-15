package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/GoFun/Fun/Trading/config"
	"github.com/KushamiNeko/GoFun/Fun/Utility/hashutils"
)

type TradingBook struct {
	recordIndex string
	noteIndex   string

	lastModified string

	title string

	bookType string
}

func NewTradingBook(title, bookType string, rIndex, nIndex, lastModified string) (*TradingBook, error) {
	tb := new(TradingBook)

	tb.title = title
	tb.bookType = bookType

	if rIndex == "" {
		tb.recordIndex = hashutils.RandString(config.IDLen)
	} else {
		tb.recordIndex = rIndex
	}

	if nIndex == "" {
		tb.noteIndex = hashutils.RandString(config.IDLen)
	} else {
		tb.noteIndex = nIndex
	}

	if lastModified == "" {
		tb.lastModified = strconv.FormatInt(time.Now().UnixNano(), 10)
	} else {
		tb.lastModified = lastModified
	}

	err := tb.validateInput()
	if err != nil {
		return nil, err
	}

	return tb, nil
}

func NewTradingBookFromInputs(entity map[string]string) (*TradingBook, error) {
	t, ok := entity["title"]
	if !ok {
		return nil, fmt.Errorf("missing title")
	}

	b, ok := entity["book_type"]
	if !ok {
		return nil, fmt.Errorf("missing bookType")
	}

	tb, err := NewTradingBook(t, b, "", "", "")
	if err != nil {
		return nil, err
	}

	return tb, nil
}

func NewTradingBookFromEntity(entity map[string]string) (*TradingBook, error) {
	t, ok := entity["title"]
	if !ok {
		return nil, fmt.Errorf("missing title")
	}

	b, ok := entity["book_type"]
	if !ok {
		return nil, fmt.Errorf("missing bookType")
	}

	ri, ok := entity["record_index"]
	if !ok {
		return nil, fmt.Errorf("missing record_index")
	}

	ni, ok := entity["note_index"]
	if !ok {
		return nil, fmt.Errorf("missing note_index")
	}

	l, ok := entity["last_modified"]
	if !ok {
		return nil, fmt.Errorf("missing last_modified")
	}

	tb, err := NewTradingBook(t, b, ri, ni, l)
	if err != nil {
		return nil, err
	}

	return tb, nil
}

func (t *TradingBook) validateInput() error {

	const (
		regexIndex        = `^[0-9a-zA-Z]+$`
		regexLastModified = `^[0-9.]+$`
		regexTitle        = `^[0-9a-zA-Z-_.]+$`
		regexBookType     = `(?:paper|live|PAPER|LIVE)$`
	)

	var re *regexp.Regexp

	re = regexp.MustCompile(regexTitle)
	if !re.MatchString(t.title) {
		return fmt.Errorf("invalid title: %s", t.title)
	}

	re = regexp.MustCompile(regexBookType)
	if !re.MatchString(t.bookType) {
		return fmt.Errorf("invalid bookType: %s", t.bookType)
	}

	re = regexp.MustCompile(regexIndex)
	if !re.MatchString(t.recordIndex) {
		return fmt.Errorf("invalid index: %s", t.recordIndex)
	}

	re = regexp.MustCompile(regexIndex)
	if !re.MatchString(t.noteIndex) {
		return fmt.Errorf("invalid index: %s", t.noteIndex)
	}

	re = regexp.MustCompile(regexLastModified)
	if !re.MatchString(t.lastModified) {
		return fmt.Errorf("invalid lastModified: %s", t.lastModified)
	}

	return nil
}

func (t *TradingBook) RecordIndex() string {
	return t.recordIndex
}

func (t *TradingBook) NoteIndex() string {
	return t.noteIndex
}

func (t *TradingBook) Title() string {
	return t.title
}

func (t *TradingBook) BookType() string {
	return t.bookType
}

func (t *TradingBook) LastModified() int64 {
	d, _ := strconv.ParseInt(t.lastModified, 10, 64)
	return d
}

func (t *TradingBook) Modified() {
	t.lastModified = strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (t *TradingBook) Entity() map[string]string {

	return map[string]string{
		"record_index":  t.recordIndex,
		"note_index":    t.noteIndex,
		"last_modified": t.lastModified,
		"title":         t.title,
		"book_type":     t.bookType,
	}
}

const (
	tradingBookFmtString = "%-[3]*[4]s%-[3]*[5]s%-[1]*[6]s%-[1]*[7]s"
)

func (t *TradingBook) Fmt() string {
	return fmt.Sprintf(
		tradingBookFmtString,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXXL,
		t.recordIndex,
		t.noteIndex,
		t.bookType,
		t.title,
	)
}

func TradingBookFmtLabels() string {
	return fmt.Sprintf(
		tradingBookFmtString,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXXL,
		"Record Index",
		"Note Index",
		"Book Type",
		"Title",
	)
}
