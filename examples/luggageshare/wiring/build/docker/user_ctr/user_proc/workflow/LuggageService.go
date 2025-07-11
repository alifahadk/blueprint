package workflow

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

type LuggageService interface {
	GetItemById(ctx context.Context, id string) (LuggageItem, error)
	AddItem(ctx context.Context, item LuggageItem) error
	FindItems(ctx context.Context, color string, length int64, breadth int64, height int64, price float64) ([]LuggageItem, error)

	// Cleanup method only used by tests
	Cleanup(ctx context.Context) error
}

var CREATE_TABLE_QUERY = `CREATE TABLE IF NOT EXISTS luggage (
	id VARCHAR(80) NOT NULL,
	color VARCHAR(1024) NOT NULL,
	length INT NOT NULL,
	breadth INT NOT NULL,
	height INT NOT NULL,
	price DECIMAL NOT NULL,
	PRIMARY KEY(id));
`

type LuggageServiceImpl struct {
	luggageDb backend.RelationalDB
}

func NewLuggageServiceImpl(ctx context.Context, db backend.RelationalDB) (LuggageService, error) {
	_, err := db.Exec(ctx, CREATE_TABLE_QUERY)
	if err != nil {
		return nil, err
	}
	return &LuggageServiceImpl{luggageDb: db}, nil
}

func (l *LuggageServiceImpl) GetItemById(ctx context.Context, id string) (LuggageItem, error) {
	// TODO: Implement

	// Step 1: Connect to the database to find the item by id
	findQuery := `SELECT luggage.id, luggage.color, luggage.length, luggage.breadth, luggage.height, luggage.price FROM luggage where luggage.id =? LIMIT 1;`
	var item LuggageItem
	err := l.luggageDb.Get(ctx, &item, findQuery, id)
	if err != nil {
		return LuggageItem{}, err
	}

	// Step 2: Return the item if found

	// Step 3: Return an error if item is not found

	return item, nil
}

func (l *LuggageServiceImpl) AddItem(ctx context.Context, item LuggageItem) error {
	// TODO: Implement

	// Step 1: ADd the item to the database
	insert_query := "INSERT INTO luggage (id, color, length, breadth, height, price) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := l.luggageDb.Exec(ctx, insert_query, item.ID, item.Color, item.Length, item.Breadth, item.Height, item.Price)
	if err != nil {
		return err
	}

	return nil
}

func (l *LuggageServiceImpl) FindItems(ctx context.Context, color string, length int64, breadth int64, height int64, price float64) ([]LuggageItem, error) {
	// TODO: Implement

	// Step 1: Search the database with the constraints to find the item!
	select_query := `SELECT luggage.id, luggage.color, luggage.length, luggage.breadth, luggage.height, luggage.price FROM luggage WHERE luggage.color =? AND luggage.length =? AND luggage.breadth =? AND luggage.height =? AND luggage.price <=?;`

	var items []LuggageItem

	err := l.luggageDb.Select(ctx, &items, select_query, color, length, breadth, height, price)
	if err != nil {
		return []LuggageItem{}, err
	}

	return items, nil
}

func (l *LuggageServiceImpl) Cleanup(ctx context.Context) error {
	_, err := l.luggageDb.Exec(ctx, `DELETE FROM luggage;`)
	return err
}
