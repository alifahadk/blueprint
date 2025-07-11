package workflow

import (
	"context"
	"errors"
)

type SearchService interface {
	Search(ctx context.Context, color string, length int64, breadth int64, height int64, price float64, startDate string, endDate string) ([]LuggageItem, error)
}

type SearchServiceImpl struct {
	resService ReservationService
	lugService LuggageService
}

func NewSearchServiceImpl(ctx context.Context, reservationService ReservationService, luggageService LuggageService) (SearchService, error) {
	return &SearchServiceImpl{resService: reservationService, lugService: luggageService}, nil
}

func (s *SearchServiceImpl) Search(ctx context.Context, color string, length int64, breadth int64, height int64, price float64, startDate string, endDate string) ([]LuggageItem, error) {
	// TODO: Implement

	// Step 1: Find items that fit the criteria (call LuggageService)
	items, err := s.lugService.FindItems(ctx, color, length, breadth, height, price)
	if err != nil {
		return []LuggageItem{}, err
	}

	var item_ids []string
	for _, item := range items {
		item_ids = append(item_ids, item.ID)
	}

	// Step 2: Filter and return the found items based on their availability (call ReservationService)
	filtered_ids, err := s.resService.CheckAvailability(ctx, item_ids, startDate, endDate)
	if err != nil {
		return []LuggageItem{}, err
	}

	// Step 3: If no items found return empty slice and an error
	if len(filtered_ids) == 0 {
		return []LuggageItem{}, errors.New("Found no items matching the criteria")
	}

	m := make(map[string]bool)
	for _, fid := range filtered_ids {
		m[fid] = true
	}

	var final_items []LuggageItem
	for _, item := range items {
		if _, ok := m[item.ID]; ok {
			final_items = append(final_items, item)
		}
	}

	// Note: Return items that have item price <= price and exactly
	// match the colour, length, breadth, and height.

	return final_items, nil
}
