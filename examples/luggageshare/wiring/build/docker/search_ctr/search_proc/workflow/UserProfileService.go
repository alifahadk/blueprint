package workflow

import (
	"context"
	"errors"

	"github.com/blueprint-uservices/blueprint/runtime/core/backend"
	"go.mongodb.org/mongo-driver/bson"
)

type UserProfileService interface {
	GetUserProfile(ctx context.Context, username string) (UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile UserProfile) error
	GetUserItemIds(ctx context.Context, username string) ([]string, error)
	AddItem(ctx context.Context, username string, item_id string) error

	// Cleanup method only used by tests
	Cleanup(ctx context.Context) error
}

type UserProfileServiceImpl struct {
	userDB backend.NoSQLDatabase
}

func NewUserProfileServiceImpl(ctx context.Context, db backend.NoSQLDatabase) (UserProfileService, error) {
	user_profile_service := &UserProfileServiceImpl{userDB: db}

	return user_profile_service, nil
}

func (u *UserProfileServiceImpl) GetUserProfile(ctx context.Context, username string) (UserProfile, error) {

	// TODO: Implement

	// Step 1: Check if the user exists
	coll, err := u.userDB.GetCollection(ctx, "user", "user")
	if err != nil {
		return UserProfile{}, err
	}
	query := bson.D{{"username", username}}
	res, err := coll.FindOne(ctx, query)
	if err != nil {
		return UserProfile{}, err
	}
	var up UserProfile
	exists, err := res.One(ctx, &up)
	if err != nil {
		return UserProfile{}, err
	}

	// Step 2: If the user exists, then return the user's profile
	if exists {
		return up, nil
	}

	// Step 3: If the user does not exist, return an error

	return UserProfile{}, errors.New("User does not exist")
}

func (u *UserProfileServiceImpl) UpdateUserProfile(ctx context.Context, profile UserProfile) error {

	// TODO: Implement

	// Step 1: Check if the user profile exists in the database
	id := profile.ID
	coll, err := u.userDB.GetCollection(ctx, "user", "user")
	if err != nil {
		return err
	}
	query := bson.D{{"id", id}}
	res, err := coll.FindOne(ctx, query)
	if err != nil {
		return err
	}
	var up UserProfile
	exists, err := res.One(ctx, &up)
	if err != nil {
		return err
	}

	// Step 2: If the profile exists, overwrite the profile in the database.
	if exists {
		n, err := coll.ReplaceOne(ctx, query, profile)
		if err != nil {
			return err
		}
		if n == 0 {
			return errors.New("Failed to update profile")
		}
	}

	// Step 3: If the profile does not exist, create a new profile in the database.

	return coll.InsertOne(ctx, profile)
}

func (u *UserProfileServiceImpl) GetUserItemIds(ctx context.Context, username string) ([]string, error) {

	// TODO: Implement

	// Step 1: Check if the user exists
	coll, err := u.userDB.GetCollection(ctx, "user", "user")
	if err != nil {
		return []string{}, err
	}
	query := bson.D{{"username", username}}
	res, err := coll.FindOne(ctx, query)
	if err != nil {
		return []string{}, err
	}
	var up UserProfile
	exists, err := res.One(ctx, &up)
	if err != nil {
		return []string{}, err
	}

	// Step 2: If the user exists, then return the items associated with the user
	if exists {
		return up.Items, nil
	}

	// Step 3: If the user does not exist, return an error

	return []string{}, errors.New("User does not exist")
}

func (u *UserProfileServiceImpl) AddItem(ctx context.Context, username string, item_id string) error {

	// TODO: Implement

	// Step 1: Check if the user exists
	coll, err := u.userDB.GetCollection(ctx, "user", "user")
	if err != nil {
		return err
	}
	query := bson.D{{"username", username}}
	res, err := coll.FindOne(ctx, query)
	if err != nil {
		return err
	}
	var up UserProfile
	exists, err := res.One(ctx, &up)
	if err != nil {
		return err
	}

	// Step 2: If the user exists, add the new item to the user's item list
	if exists {
		up.Items = append(up.Items, item_id)
		n, err := coll.ReplaceOne(ctx, query, up)
		if err != nil {
			return err
		}
		if n == 0 {
			return errors.New("Failed to add item to the user")
		}
		return nil
	}

	// Step 3: If the user does not exist, return an error

	return errors.New("User does not exist")
}

func (u *UserProfileServiceImpl) Cleanup(ctx context.Context) error {
	coll, err := u.userDB.GetCollection(ctx, "user", "user")
	if err != nil {
		return err
	}

	err = coll.DeleteMany(ctx, bson.D{})
	return err
}
