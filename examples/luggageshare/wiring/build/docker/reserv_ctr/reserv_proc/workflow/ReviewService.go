package workflow

import (
	"context"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"go.mongodb.org/mongo-driver/bson"
)

type ReviewService interface {
	AddReview(ctx context.Context, review Review) error
	GetReviewsForItem(ctx context.Context, id string) ([]Review, error)
	GetReviewsForUser(ctx context.Context, userid string) ([]Review, error)

	// Cleanup method only used by tests
	Cleanup(ctx context.Context) error
}

type ReviewServiceImpl struct {
	reviewDB backend.NoSQLDatabase
}

func NewReviewServiceImpl(ctx context.Context, db backend.NoSQLDatabase) (ReviewService, error) {
	review_service := &ReviewServiceImpl{reviewDB: db}

	return review_service, nil
}

func (r *ReviewServiceImpl) AddReview(ctx context.Context, review Review) error {
	// TODO: Implement

	// Step 1: Add review to the review database
	coll, err := r.reviewDB.GetCollection(ctx, "reviews", "reviews")
	if err != nil {
		return err
	}

	return coll.InsertOne(ctx, review)
}

func (r *ReviewServiceImpl) GetReviewsForItem(ctx context.Context, id string) ([]Review, error) {
	// TODO: Implement

	// Step 1: Search the database for reviews based on item id (LuggageID) and return the found reviews
	var reviews []Review
	coll, err := r.reviewDB.GetCollection(ctx, "reviews", "reviews")
	if err != nil {
		return reviews, err
	}
	query := bson.D{{"luggageid", id}}
	res, err := coll.FindMany(ctx, query)
	if err != nil {
		return reviews, err
	}
	err = res.All(ctx, &reviews)
	if err != nil {
		return reviews, err
	}

	return reviews, nil
}

func (r *ReviewServiceImpl) GetReviewsForUser(ctx context.Context, user_id string) ([]Review, error) {
	// TODO: Implement

	// Step 1: Search the database for reviews based on user id and return the found reviews

	var reviews []Review
	coll, err := r.reviewDB.GetCollection(ctx, "reviews", "reviews")
	if err != nil {
		return reviews, err
	}
	query := bson.D{{"user", user_id}}
	res, err := coll.FindMany(ctx, query)
	if err != nil {
		return reviews, err
	}
	err = res.All(ctx, &reviews)
	if err != nil {
		return reviews, err
	}

	return reviews, nil
}

func (r *ReviewServiceImpl) Cleanup(ctx context.Context) error {
	coll, err := r.reviewDB.GetCollection(ctx, "reviews", "reviews")
	if err != nil {
		return err
	}

	err = coll.DeleteMany(ctx, bson.D{})
	return err
} 
