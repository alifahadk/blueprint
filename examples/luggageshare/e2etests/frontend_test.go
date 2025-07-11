package e2etests

import (
	"context"
	"testing"

	"github.com/blueprint-uservices/blueprint/runtime/core/registry"
	"github.com/blueprint-uservices/blueprint/runtime/plugins/simplenosqldb"
	"github.com/blueprint-uservices/blueprint/runtime/plugins/sqlitereldb"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var luggageServiceRegistry = registry.NewServiceRegistry[workflow.LuggageService]("luggage_service")
var reservationServiceRegistry = registry.NewServiceRegistry[workflow.ReservationService]("reservation_service")
var reviewServiceRegistry = registry.NewServiceRegistry[workflow.ReviewService]("review_service")
var searchServiceRegistry = registry.NewServiceRegistry[workflow.SearchService]("search_service")
var frontendServiceRegistry = registry.NewServiceRegistry[workflow.FrontendService]("frontend_service")
var userProfileServiceRegistry = registry.NewServiceRegistry[workflow.UserProfileService]("user_profile_service")

func init() {

	reviewServiceRegistry.Register("local", func(ctx context.Context) (workflow.ReviewService, error) {
		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewReviewServiceImpl(ctx, db)
	})

	userProfileServiceRegistry.Register("local", func(ctx context.Context) (workflow.UserProfileService, error) {
		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewUserProfileServiceImpl(ctx, db)
	})

	luggageServiceRegistry.Register("local", func(ctx context.Context) (workflow.LuggageService, error) {
		db, err := sqlitereldb.NewSqliteRelDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewLuggageServiceImpl(ctx, db)
	})

	reservationServiceRegistry.Register("local", func(ctx context.Context) (workflow.ReservationService, error) {
		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewReservationServiceImpl(ctx, db)
	})

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

	frontendServiceRegistry.Register("local", func(ctx context.Context) (workflow.FrontendService, error) {
		searchService, err := searchServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		resService, err := reservationServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		reviewService, err := reviewServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		userService, err := userProfileServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		luggageService, err := luggageServiceRegistry.Get(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewFrontendServiceImpl(ctx, searchService, resService, reviewService, userService, luggageService)
	})
}

func TestFrontendAddItem(t *testing.T) {
	ctx := context.Background()
	service, err := frontendServiceRegistry.Get(ctx)
	require.NoError(t, err)

	// Create a new user profile here
	err = service.CreateUserProfile(ctx, "addItemUser", "addItemUser@cldrel25.com", "MPI-SWS Building")
	require.NoError(t, err)

	// Add item to an existing user
	_, err = service.AddItem(ctx, "blue", 10, 10, 10, 1.50, "addItemUser")
	require.NoError(t, err)

	// Add item to a non-existing user
	_, err = service.AddItem(ctx, "blue", 5, 5, 5, 2.50, "noExistUser")
	require.Error(t, err)
}

func TestFrontendUpdateUserProfile(t *testing.T) {
	ctx := context.Background()
	service, err := frontendServiceRegistry.Get(ctx)
	require.NoError(t, err)

	// Test insertion of 1 user works
	username := "frontendTestUser"
	email := username + "@cldrel25.com"
	address := "MPI-SWS Building"
	err = service.CreateUserProfile(ctx, username, email, address)
	require.NoError(t, err)

	// Test duplicate user doesn't work
	err = service.CreateUserProfile(ctx, username, "newemail@cldrel25.com", address)
	require.Error(t, err)
}

func TestFrontendWriteReview(t *testing.T) {
	ctx := context.Background()
	service, err := frontendServiceRegistry.Get(ctx)
	require.NoError(t, err)

	// Test insertion of review for a luggage id that doesn't exist!
	err = service.WriteReview(ctx, "adgfilsadgfi", "Good luggage", "badUsername", 10)
	require.Error(t, err)

	// Insert a new user
	err = service.CreateUserProfile(ctx, "writeReviewUser", "writeReviewUser@cldrel25.com", "MPI-SWS Building")
	require.NoError(t, err)

	// Insert a new luggage item for the user
	lugid, err := service.AddItem(ctx, "blue", 10, 10, 10, 1.50, "writeReviewUser")
	require.NoError(t, err)

	// Write a review for a luggage id and for a user that does exist
	err = service.WriteReview(ctx, lugid, "Good luggage item. Will use again", "writeReviewUser", 10)
	require.NoError(t, err)

	// Try to write a review for a luggage id that exists but from a user that does not exist
	err = service.WriteReview(ctx, lugid, "Bad luggage item but hey I am a bot", "botUser", 1)
	require.Error(t, err)
}

func TestFrontendMakeReservation(t *testing.T) {
	ctx := context.Background()
	service, err := frontendServiceRegistry.Get(ctx)
	require.NoError(t, err)

	// Try to make a reservation for a luggage id that doesn't exist!
	success, err := service.MakeReservation(ctx, "ksdjbfgikusbgfi", "botUser", "2025-07-01", "2025-07-10")
	assert.False(t, success)
	require.Error(t, err)

	// Insert a new user
	err = service.CreateUserProfile(ctx, "makeResUser1", "makeResUser1@cldrel25.com", "MPI-SWS Building")
	require.NoError(t, err)

	// Insert a new luggage item for the user
	lugid, err := service.AddItem(ctx, "blue", 10, 10, 10, 1.50, "makeResUser1")
	require.NoError(t, err)

	// Insert a 2nd user
	err = service.CreateUserProfile(ctx, "makeResUser2", "makeResUser2@cldrel25.com", "MPI-SWS Building")
	require.NoError(t, err)

	// Make reservation for the luggage id
	success, err = service.MakeReservation(ctx, lugid, "makeResUser2", "2025-07-01", "2025-07-10")
	require.NoError(t, err)
	assert.True(t, success)

	// Try to make an overlapping reservation
	success, err = service.MakeReservation(ctx, lugid, "makeResUser2", "2025-07-04", "2025-07-14")
	require.Error(t, err)
	assert.False(t, success)

}

func TestFrontendSearch(t *testing.T) {
	ctx := context.Background()
	service, err := frontendServiceRegistry.Get(ctx)
	require.NoError(t, err)

	color := "pink"
	length := int64(30)
	width := int64(30)
	height := int64(20)
	price := float64(3.00)
	startDate := "2025-07-01"
	endDate := "2025-07-10"
	startdate_reserved := "2025-07-03"
	enddate_reserved := "2025-07-07"

	infos, err := service.Search(ctx, color, length, width, height, price, startDate, endDate)
	require.Error(t, err)

	// Insert a new user
	err = service.CreateUserProfile(ctx, "searchUser", "searchUser@cldrel25.com", "MPI-SWS Building")
	require.NoError(t, err)

	// Add a few items!
	added_items := []string{}
	for i := 0; i < 5; i++ {
		id, err := service.AddItem(ctx, color, length, width, height, price, "searchUser")
		require.NoError(t, err)
		added_items = append(added_items, id)
	}
	// Add a different type of item
	extra_id, err := service.AddItem(ctx, "blue", 5, 5, 5, 4.5, "searchUser")

	infos, err = service.Search(ctx, color, length, width, height, price, startDate, endDate)
	require.NoError(t, err)
	assert.Equal(t, 5, len(infos))

	for _, info := range infos {
		// Make sure we didn't get the different type of item in the returned list
		assert.NotEqual(t, info.Item.ID, extra_id)
		assert.Equal(t, info.Item.Color, color)
		assert.Equal(t, info.Item.Length, length)
		assert.Equal(t, info.Item.Breadth, width)
		assert.Equal(t, info.Item.Height, height)
	}

	// Reserve all of the pink items
	for _, id := range added_items {
		success, err := service.MakeReservation(ctx, id, "searchUser", startDate, endDate)
		require.NoError(t, err)
		assert.True(t, success)
	}

	infos, err = service.Search(ctx, color, length, width, height, price, startdate_reserved, enddate_reserved)
	require.Error(t, err)
	assert.Equal(t, 0, len(infos))
}
