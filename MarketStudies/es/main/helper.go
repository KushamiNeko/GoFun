package main

import (
	"os"
	"path/filepath"

	"github.com/KushamiNeko/go_fun/trading/agent"
	"github.com/KushamiNeko/go_fun/trading/context"
	"github.com/KushamiNeko/go_fun/utils/database"
)

func newTradingAgent(book string) (*agent.TradingAgent, error) {
	var err error

	db := database.NewFileDB(
		filepath.Join(
			os.Getenv("HOME"),
			"Documents/database/filedb/futures_wizards",
		),
		database.JsonEngine,
	)

	ctx := context.NewContext(db)
	err = ctx.Login("aa")
	if err != nil {
		return nil, err
	}

	tradeAgent, err := agent.NewTradingAgent(ctx, false)
	if err != nil {
		return nil, err
	}

	books, err := tradeAgent.Books()
	if err != nil {
		return nil, err
	}

	for _, b := range books {
		if b.Title() == book {
			err = tradeAgent.SetReading(b)
			if err != nil {
				return nil, err
			}

			break
		}
	}

	return tradeAgent, nil
}
