package tests

import (
	"context"
	"testing"

	"github.com/blueprint-uservices/blueprint/runtime/core/registry"
	"github.com/blueprint-uservices/blueprint/runtime/plugins/sqlitereldb"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
	"github.com/stretchr/testify/require"
)

var luggageServiceRegistry = registry.NewServiceRegistry[workflow.LuggageService]("luggage_service")

func init() {
	luggageServiceRegistry.Register("local", func(ctx context.Context) (workflow.LuggageService, error) {
		db, err := sqlitereldb.NewSqliteRelDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewLuggageServiceImpl(ctx, db)
	})
}

func TestLuggageServiceAddItem(t *testing.T) {
	ctx := context.Background()
	service, err := luggageServiceRegistry.Get(ctx)
	require.NoError(t, err)

	// Add Item
	{
		item := workflow.LuggageItem{}
		err := service.AddItem(ctx, item)
		require.NoError(t, err)
	}
}

func TestLuggageServiceGetItemById(t *testing.T) {
	ctx := context.Background()
	service, err := luggageServiceRegistry.Get(ctx)
	require.NoError(t, err)

	var id string = "test"
	var color string = "blue"
	var length int64 = 5
	var breadth int64 = 10
	var height int64 = 3
	var price float64 = 2
	expected_item := workflow.LuggageItem{ID: id, Color: color, Length: length,
		Breadth: breadth, Height: height, Price: price}

	// Try to get item that doesn't exist
	{
		_, err := service.GetItemById(ctx, id)
		require.Error(t, err)
	}

	// Add item
	{
		err := service.AddItem(ctx, expected_item)
		require.NoError(t, err)
	}

	// Get Item
	{
		actual_item, err := service.GetItemById(ctx, id)
		require.NoError(t, err)
		require.Equal(t, expected_item, actual_item)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestLuggageServiceFindItems(t *testing.T) {
	ctx := context.Background()
	service, err := luggageServiceRegistry.Get(ctx)
	require.NoError(t, err)

	var id string = "test"
	var color string = "blue"
	var length int64 = 5
	var breadth int64 = 10
	var height int64 = 3
	var price float64 = 2
	expected_item := workflow.LuggageItem{ID: id, Color: color, Length: length,
		Breadth: breadth, Height: height, Price: price}

	// Try to find item that doesn't exist
	{
		_, err := service.FindItems(ctx, color, length, breadth, height, price)
		require.NoError(t, err)
	}

	// Add item
	{
		err := service.AddItem(ctx, expected_item)
		require.NoError(t, err)
	}

	// Find item
	{
		actual_item, err := service.FindItems(ctx, color, length, breadth, height, price)
		require.NoError(t, err)
		require.Equal(t, expected_item, actual_item[0])
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}
