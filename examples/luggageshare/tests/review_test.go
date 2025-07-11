package tests

import (
	"context"
	"testing"

	"github.com/blueprint-uservices/blueprint/runtime/core/registry"
	"github.com/blueprint-uservices/blueprint/runtime/plugins/simplenosqldb"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
	"github.com/stretchr/testify/require"
)

var reviewServiceRegistry = registry.NewServiceRegistry[workflow.ReviewService]("review_service")

func init() {
	reviewServiceRegistry.Register("local", func(ctx context.Context) (workflow.ReviewService, error) {
		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewReviewServiceImpl(ctx, db)
	})
}

func TestReviewServiceAdd(t *testing.T) {
	ctx := context.Background()
	service, err := reviewServiceRegistry.Get(ctx)
	require.NoError(t, err)

	id := "test"
	luggage_id := "test_review"
	text := "test text"
	var rating int64 = 10
	user := "test_user"
	review := workflow.Review{ID: id, LuggageID: luggage_id, Text: text,
		Rating: rating, User: user}

	// Try to add review
	{
		err := service.AddReview(ctx, review)
		require.NoError(t, err)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestReviewServiceGetReviewsForItem(t *testing.T) {
	ctx := context.Background()
	service, err := reviewServiceRegistry.Get(ctx)
	require.NoError(t, err)

	id := "test"
	luggage_id := "test_review"
	text := "test text"
	var rating int64 = 10
	user := "test_user"
	expected_review := workflow.Review{ID: id, LuggageID: luggage_id, Text: text,
		Rating: rating, User: user}

	// Try to add review
	{
		err := service.AddReview(ctx, expected_review)
		require.NoError(t, err)
	}

	// Get review
	{
		actual_review, err := service.GetReviewsForItem(ctx, luggage_id)
		require.NoError(t, err)
		require.Equal(t, expected_review, actual_review[0])
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestReviewServiceGetReviewsForUser(t *testing.T) {
	ctx := context.Background()
	service, err := reviewServiceRegistry.Get(ctx)
	require.NoError(t, err)

	id := "test"
	luggage_id := "test_review"
	text := "test text"
	var rating int64 = 10
	user := "test_user"
	expected_review := workflow.Review{ID: id, LuggageID: luggage_id, Text: text,
		Rating: rating, User: user}

	// Try to add review
	{
		err := service.AddReview(ctx, expected_review)
		require.NoError(t, err)
	}

	// Get review
	{
		actual_review, err := service.GetReviewsForUser(ctx, user)
		require.NoError(t, err)
		require.Equal(t, expected_review, actual_review[0])
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}
