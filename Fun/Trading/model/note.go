package model

import (
	"fmt"
	"regexp"
	"time"
)

type TradingNote struct {
	time string
	note string
}

func NewTradingNoteFromEntity(entity map[string]string) (*TradingNote, error) {
	dt, ok := entity["time"]
	if !ok {
		return nil, fmt.Errorf("missing time")
	}

	regex := regexp.MustCompile(`^(\d{8}|\d{14})$`)
	if !regex.MatchString(dt) {
		return nil, fmt.Errorf("invalid time format")
	}

	note, ok := entity["note"]
	if !ok {
		return nil, fmt.Errorf("missing note")
	}

	return &TradingNote{
		time: dt,
		note: note,
	}, nil
}

func (t *TradingNote) Time() time.Time {
	var dt time.Time
	var err error

	regex := regexp.MustCompile(`^\d{8}$`)
	if regex.MatchString(t.time) {
		dt, err = time.Parse("20060102", t.time)
		if err != nil {
			panic(err)
		}
	} else {
		dt, err = time.Parse("20060102150405", t.time)
		if err != nil {
			panic(err)
		}
	}

	return dt
}

func (t *TradingNote) Note() string {
	return t.note
}
