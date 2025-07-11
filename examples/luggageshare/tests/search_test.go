package tests

import (
	"context"
	"testing"

	"github.com/blueprint-uservices/blueprint/runtime/core/registry"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
	"github.com/stretchr/testify/require"
)

var searchServiceRegistry = registry.NewServiceRegistry[workflow.SearchService]("search_service")

func init() {
	// Create search service
	searchServiceRegistry.Register("local", func(ctx context.Context) (workflow.SearchService, error) {
		// Create reservation service
		reservation_service, err := reservationServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		// Create luggage service
		luggage_service, err := luggageServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewSearchServiceImpl(ctx, reservation_service, luggage_service)
	})
}
func TestSearchServiceSearch(t *testing.T) {
	ctx := context.Background()
	search_service, err := searchServiceRegistry.Get(ctx)
	require.NoError(t, err)
	luggage_service, err := luggageServiceRegistry.Get(ctx)
	require.NoError(t, err)
	reservation_service, err := reservationServiceRegistry.Get(ctx)
	require.NoError(t, err)

	unique_color := "red"
	var unique_length int64 = 999
	var unique_breadth int64 = 999
	var unique_height int64 = 999

	id0 := "id0"
	color0 := "blue"
	var length0 int64 = 10
	var breadth0 int64 = 20
	var height0 int64 = 5
	var price0 float64 = 100
	item0 := workflow.LuggageItem{ID: id0, Color: color0, Length: length0,
		Breadth: breadth0, Height: height0, Price: price0}

	id1 := "id1"
	color1 := "blue"
	var length1 int64 = 10
	var breadth1 int64 = 20
	var height1 int64 = 5
	var price1 float64 = 75
	item1 := workflow.LuggageItem{ID: id1, Color: color1, Length: length1,
		Breadth: breadth1, Height: height1, Price: price1}

	id2 := "id2"
	color2 := "blue"
	var length2 int64 = 10
	var breadth2 int64 = 20
	var height2 int64 = 5
	var price2 float64 = 200
	item2 := workflow.LuggageItem{ID: id2, Color: color2, Length: length2,
		Breadth: breadth2, Height: height2, Price: price2}

	id3 := "id3"
	color3 := "green"
	var length3 int64 = 10
	var breadth3 int64 = 20
	var height3 int64 = 5
	var price3 float64 = 55
	item3 := workflow.LuggageItem{ID: id3, Color: color3, Length: length3,
		Breadth: breadth3, Height: height3, Price: price3}

	// Add luggage items
	{
		luggage_service.AddItem(ctx, item0)
		luggage_service.AddItem(ctx, item1)
		luggage_service.AddItem(ctx, item2)
		luggage_service.AddItem(ctx, item3)
	}

	// Search for some items
	{
		expected := []workflow.LuggageItem{item0, item1}
		startDate := "2001-01-01"
		endDate := "2001-02-01"
		actual, err := search_service.Search(ctx, color0, length0, breadth0, height0, price0, startDate, endDate)
		require.NoError(t, err)
		require.ElementsMatch(t, expected, actual)
	}

	// Search colour doesn't match
	{
		expected := []workflow.LuggageItem{}
		startDate := "2001-01-01"
		endDate := "2001-02-01"
		actual, err := search_service.Search(ctx, unique_color, length0, breadth0, height0, price0, startDate, endDate)
		require.Error(t, err)
		require.ElementsMatch(t, expected, actual)
	}

	// Search length doesn't match
	{
		expected := []workflow.LuggageItem{}
		startDate := "2001-01-01"
		endDate := "2001-02-01"
		actual, err := search_service.Search(ctx, color0, unique_length, breadth0, height0, price0, startDate, endDate)
		require.Error(t, err)
		require.ElementsMatch(t, expected, actual)
	}

	// Search breadth doesn't match
	{
		expected := []workflow.LuggageItem{}
		startDate := "2001-01-01"
		endDate := "2001-02-01"
		actual, err := search_service.Search(ctx, color0, length0, unique_breadth, height0, price0, startDate, endDate)
		require.Error(t, err)
		require.ElementsMatch(t, expected, actual)
	}

	// Search height doesn't match
	{
		expected := []workflow.LuggageItem{}
		startDate := "2001-01-01"
		endDate := "2001-02-01"
		actual, err := search_service.Search(ctx, color0, length0, breadth0, unique_height, price0, startDate, endDate)
		require.Error(t, err)
		require.ElementsMatch(t, expected, actual)
	}

	// Reserve an item
	{
		startDate := "2001-01-01"
		endDate := "2001-02-01"
		reservation_service.MakeReservation(ctx, id1, "", startDate, endDate)
	}

	// Search for items with reservations filtered
	{
		expected := []workflow.LuggageItem{item0}
		startDate := "2001-01-01"
		endDate := "2001-01-10"
		actual, err := search_service.Search(ctx, color0, length0, breadth0, height0, price0, startDate, endDate)
		require.Len(t, actual, 1)
		require.NoError(t, err)
		require.ElementsMatch(t, expected, actual)
	}

}
