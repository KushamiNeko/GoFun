package agent

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KushamiNeko/GoFun/Trading/config"
	"github.com/KushamiNeko/GoFun/Trading/context"
	"github.com/KushamiNeko/GoFun/Trading/model"
	"github.com/KushamiNeko/GoFun/Utility/database"
	"github.com/KushamiNeko/GoFun/Utility/input"
)

type TradingAgent struct {
	ctx *context.Context

	books   []*model.TradingBook
	reading *model.TradingBook

	bookType  string
	tradingDB string
}

func NewTradingAgentCompact(dbpath, username, book string) (*TradingAgent, error) {
	var err error

	db := database.NewFileDB(
		dbpath,
		database.JsonDB,
	)

	ctx := context.NewContext(db)
	err = ctx.Login(username)
	if err != nil {
		return nil, err
	}

	tradeAgent, err := NewTradingAgent(ctx, false)
	if err != nil {
		return nil, err
	}

	if book != "" {
		books, err := tradeAgent.Books()
		if err != nil {
			return nil, err
		}

		found := false
		for _, b := range books {
			if b.Title() == book {
				err = tradeAgent.SetReading(b)
				if err != nil {
					return nil, err
				}

				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("unknown book: %s", book)
		}
	}

	return tradeAgent, nil
}

func NewTradingAgent(ctx *context.Context, live bool) (*TradingAgent, error) {
	t := new(TradingAgent)
	t.ctx = ctx

	if live {
		t.bookType = "live"
		t.tradingDB = config.DBLiveTrading
	} else {
		t.bookType = "paper"
		t.tradingDB = config.DBPaperTrading
	}

	books, err := t.ctx.DB().Find(
		config.DBTradingBooks,
		t.ctx.User().UID(),
		map[string]string{
			"book_type": t.bookType,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(books) > 0 {

		t.books = make([]*model.TradingBook, len(books))

		for i, b := range books {
			t.books[i], err = model.NewTradingBookFromEntity(b)
			if err != nil {
				return nil, err
			}
		}

		sort.Slice(t.books, func(i, j int) bool {
			return t.books[i].LastModified() > t.books[j].LastModified()
		})

		t.reading = t.books[0]
	}

	return t, nil
}

func (t *TradingAgent) NewBook(inputs map[string]string) error {
	n := input.InputsAbbreviation(inputs, map[string]string{
		"t": "title",
	})

	n["book_type"] = t.bookType

	book, err := model.NewTradingBookFromInputs(n)
	if err != nil {
		return err
	}

	err = t.ctx.DB().Insert(config.DBTradingBooks, t.ctx.User().UID(), book.Entity())
	if err != nil {
		return err
	}

	t.books = append(t.books, book)
	sort.Slice(t.books, func(i, j int) bool {
		return t.books[i].LastModified() > t.books[j].LastModified()
	})

	t.reading = book

	return nil
}

func (t *TradingAgent) UpdateBook() error {

	t.reading.Modified()

	err := t.ctx.DB().Replace(
		config.DBTradingBooks,
		t.ctx.User().UID(),
		map[string]string{
			"record_index": t.reading.RecordIndex(),
		},
		t.reading.Entity(),
	)

	return err
}

func (t *TradingAgent) ChangeBook(inputs map[string]string) error {
	n := input.InputsAbbreviation(inputs, map[string]string{
		"i": "index",
	})

	i, ok := n["index"]
	if !ok {
		return fmt.Errorf("missing index")
	}

	if i == "" {
		return fmt.Errorf("invalid book index")
	}

	for _, b := range t.books {
		if strings.HasPrefix(b.RecordIndex(), i) {
			t.reading = b
			return nil
		}
	}

	return fmt.Errorf("unknown book index: %s", i)
}

func (t *TradingAgent) NewTransaction(inputs map[string]string) error {
	n := input.InputsAbbreviation(inputs, map[string]string{
		"t": "time",
		"s": "symbol",
		"o": "operation",
		"q": "quantity",
		"p": "price",
		"n": "note",
	})

	transaction, err := model.NewFuturesTransactionFromInputs(n)
	if err != nil {
		return err
	}

	err = t.ctx.DB().Insert(
		t.tradingDB,
		t.reading.RecordIndex(),
		transaction.Entity(),
	)
	if err != nil {
		return err
	}

	err = t.UpdateBook()
	if err != nil {
		return err
	}

	return nil
}

func (t *TradingAgent) Positions() ([]*model.FuturesTransaction, error) {
	transactions, err := t.Transactions()
	if err != nil {
		return nil, err
	}

	trades, err := t.processTrades(transactions)
	if err != nil {
		return nil, err
	}

	if len(trades) == 0 {
		return transactions, nil
	}

	tb := make(map[string]*model.FuturesTrade)

	for i := len(trades) - 1; i >= 0; i-- {
		t := trades[i]

		if v, ok := tb[t.Symbol()]; !ok {
			tb[t.Symbol()] = t
		} else {
			if t.CloseTime().Equal(v.CloseTime()) {
				if t.CloseTimeStamp() > v.CloseTimeStamp() {
					tb[t.Symbol()] = t
				}
			} else {
				if t.CloseTime().After(v.CloseTime()) {
					tb[t.Symbol()] = t
				}
			}
		}
	}

	positions := make([]*model.FuturesTransaction, 0)

	for _, t := range transactions {
		if v, ok := tb[t.Symbol()]; !ok {
			positions = append(positions, t)
		} else {
			if t.Time().Equal(v.CloseTime()) {
				if t.TimeStamp() > v.CloseTimeStamp() {
					positions = append(positions, t)
				}
			} else {
				if t.Time().After(v.CloseTime()) {
					positions = append(positions, t)
				}
			}
		}
	}

	return positions, nil
}

func (t *TradingAgent) Books() ([]*model.TradingBook, error) {
	if t.books == nil || len(t.books) == 0 {
		return nil, fmt.Errorf("empty books")
	} else {
		return t.books, nil
	}
}

func (t *TradingAgent) Transactions() ([]*model.FuturesTransaction, error) {
	results, err := t.ctx.DB().Find(
		t.tradingDB,
		t.reading.RecordIndex(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.FuturesTransaction, len(results))

	for i, r := range results {
		t, err := model.NewFuturesTransactionFromEntity(r)
		if err != nil {
			return nil, err
		}

		transactions[i] = t
	}

	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].Time().Equal(transactions[j].Time()) {
			return transactions[i].TimeStamp() < transactions[j].TimeStamp()
		} else {
			return transactions[i].Time().Before(transactions[j].Time())
		}
	})

	return transactions, nil
}

func (t *TradingAgent) Trades() ([]*model.FuturesTrade, error) {
	transactions, err := t.Transactions()
	if err != nil {
		return nil, err
	}

	trades, err := t.processTrades(transactions)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

func (t *TradingAgent) Statistic() (*model.Statistic, error) {
	trades, err := t.Trades()
	if err != nil {
		return nil, err
	}

	return model.NewStatistic(trades)
}

func (t *TradingAgent) processTrades(transactions []*model.FuturesTransaction) ([]*model.FuturesTrade, error) {

	type orders struct {
		t []*model.FuturesTransaction
		p int
	}

	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].Time().Equal(transactions[j].Time()) {
			return transactions[i].TimeStamp() < transactions[j].TimeStamp()
		} else {
			return transactions[i].Time().Before(transactions[j].Time())
		}
	})

	trades := make([]*model.FuturesTrade, 0)
	book := make(map[string]*orders)

	for _, t := range transactions {
		o, ok := book[t.Symbol()]
		if !ok {
			book[t.Symbol()] = &orders{
				t: []*model.FuturesTransaction{t},
				p: t.Action(),
			}
		} else {
			o.p += t.Action()
			o.t = append(book[t.Symbol()].t, t)
		}

		o = book[t.Symbol()]

		if o.p == 0 {
			trade, err := model.NewFuturesTrade(o.t)
			if err != nil {
				return nil, err
			}

			trades = append(trades, trade)
			delete(book, t.Symbol())
		}
	}

	return trades, nil
}

func (t *TradingAgent) Reading() (*model.TradingBook, error) {
	if t.reading != nil {
		return t.reading, nil
	} else {
		return nil, fmt.Errorf("empty books")
	}
}

func (t *TradingAgent) SetReading(book *model.TradingBook) error {
	found := false
	for _, b := range t.books {
		if b.Title() == book.Title() {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid book")
	}

	t.reading = book

	return nil
}
