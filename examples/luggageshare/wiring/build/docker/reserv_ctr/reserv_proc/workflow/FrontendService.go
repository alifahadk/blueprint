package workflow

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type FrontendService interface {
	Search(ctx context.Context, color string, length int64, breadth int64, height int64, price float64, startDate string, endDate string) ([]LuggageInfo, error)
	CreateUserProfile(ctx context.Context, username string, email string, address string) error
	AddItem(ctx context.Context, color string, length int64, breadth int64, height int64, price float64, username string) (string, error)
	MakeReservation(ctx context.Context, luggageid string, username string, startdate string, enddate string) (bool, error)
	WriteReview(ctx context.Context, luggageid string, text string, username string, rating int64) error
}

type FrontendServiceImpl struct {
	searchService SearchService
	resService    ReservationService
	reviewService ReviewService
	userService   UserProfileService
	lugService    LuggageService
}

func NewFrontendServiceImpl(ctx context.Context, searchService SearchService, resService ReservationService, reviewService ReviewService, userService UserProfileService, lugService LuggageService) (FrontendService, error) {
	return &FrontendServiceImpl{searchService: searchService, resService: resService, reviewService: reviewService, userService: userService, lugService: lugService}, nil
}

func gen_uuid(prefix string) string {
	return prefix + uuid.NewString()
}

func (f *FrontendServiceImpl) Search(ctx context.Context, color string, length int64, breadth int64, height int64, price float64, startDate string, endDate string) ([]LuggageInfo, error) {
	// TODO: Implement

	// Step 1: Call SearchService to find possible items
	items, err := f.searchService.Search(ctx, color, length, breadth, height, price, startDate, endDate)
	if err != nil {
		return []LuggageInfo{}, err
	}

	// Step 2: Call ReviewService to get reviews for each item

	// Step 3: For each item, construct an info object
	// This could be made non-blocking waiting
	var infos []LuggageInfo
	for _, item := range items {
		reviews, err := f.reviewService.GetReviewsForItem(ctx, item.ID)
		if err != nil {
			return infos, err
		}
		info := LuggageInfo{Item: item, Reviews: reviews}
		infos = append(infos, info)
	}

	// Step 4: Return the list of objects

	return infos, nil
}

func (f *FrontendServiceImpl) CreateUserProfile(ctx context.Context, username string, email string, address string) error {
	// TODO: Implement

	// Step 1: Check if the user exists, return an error if the user exists
	_, err := f.userService.GetUserProfile(ctx, username)
	if err == nil {
		return errors.New("Username already in use")
	}

	// Step 2: Generate unique id
	user_id := gen_uuid("user")

	// Step 3: Create a new user profile and contact the userprofileservice to add the new profile to the database
	up := UserProfile{Username: username, ID: user_id, Email: email, Address: address}

	return f.userService.UpdateUserProfile(ctx, up)
}

func (f *FrontendServiceImpl) AddItem(ctx context.Context, color string, length int64, breadth int64, height int64, price float64, username string) (string, error) {
	// TODO: Implement

	// Step 1: Generate unique id
	item_id := gen_uuid("item")

	// Step 2: Create a new Item object
	item := LuggageItem{ID: item_id, Color: color, Length: length, Breadth: breadth, Height: height, Price: price}

	// Step 3: Add item to the user profile
	err := f.userService.AddItem(ctx, username, item.ID)
	if err != nil {
		return "", err
	}

	// Step 4: Add item to the luggage database
	err = f.lugService.AddItem(ctx, item)
	if err != nil {
		return "", err
	}

	// Return the luggage id

	return item_id, nil
}

func (f *FrontendServiceImpl) MakeReservation(ctx context.Context, luggageid string, username string, startdate string, enddate string) (bool, error) {
	// TODO: Implement

	// Step 1: Check that the luggage with the provided id exists
	_, err := f.lugService.GetItemById(ctx, luggageid)
	if err != nil {
		return false, err
	}

	// Step 2: Check that the username of the reservation maker exists
	_, err = f.userService.GetUserProfile(ctx, username)
	if err != nil {
		return false, err
	}

	// Step 3: Contact the reservation service to make a new reservation for the luggage item

	return f.resService.MakeReservation(ctx, luggageid, username, startdate, enddate)
}

func (f *FrontendServiceImpl) WriteReview(ctx context.Context, luggageid string, text string, username string, rating int64) error {
	// TODO: Implement

	// Step 1: Generate a unique review id
	review_id := gen_uuid("review")

	// Step 2: Check that the luggage with the provided id exists
	_, err := f.lugService.GetItemById(ctx, luggageid)
	if err != nil {
		return err
	}

	// Step 3: Check that the username of the review writer exists
	_, err = f.userService.GetUserProfile(ctx, username)
	if err != nil {
		return err
	}

	// Step 3: Call the review service to add a new review to the database

	review := Review{ID: review_id, LuggageID: luggageid, User: username, Text: text, Rating: rating}

	return f.reviewService.AddReview(ctx, review)
}
