package tests

import (
	"context"
	"testing"

	"github.com/blueprint-uservices/blueprint/runtime/core/registry"
	"github.com/blueprint-uservices/blueprint/runtime/plugins/simplenosqldb"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
	"github.com/stretchr/testify/require"
)

var userProfileServiceRegistry = registry.NewServiceRegistry[workflow.UserProfileService]("user_profile_service")

func init() {
	userProfileServiceRegistry.Register("local", func(ctx context.Context) (workflow.UserProfileService, error) {
		db, err := simplenosqldb.NewSimpleNoSQLDB(ctx)
		if err != nil {
			return nil, err
		}

		return workflow.NewUserProfileServiceImpl(ctx, db)
	})
}

func TestUserProfileServiceGetUserProfile(t *testing.T) {
	ctx := context.Background()
	service, err := userProfileServiceRegistry.Get(ctx)
	require.NoError(t, err)

	fake_username := "fake"
	id := "test_id"
	username := "test_username"
	email := "test_email"
	items := []string{"0", "1", "2"}
	address := "test_address"
	expected_profile := workflow.UserProfile{ID: id, Username: username, Email: email,
		Items: items, Address: address}

	// Add user
	{
		err := service.UpdateUserProfile(ctx, expected_profile)
		require.NoError(t, err)
	}

	// Get user
	{
		actual_profile, err := service.GetUserProfile(ctx, username)
		require.NoError(t, err)
		require.Equal(t, expected_profile, actual_profile)
	}

	// Try to get user that doesn't exist
	{
		_, err := service.GetUserProfile(ctx, fake_username)
		require.Error(t, err)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestUserProfileServiceUpdateUserProfile(t *testing.T) {
	ctx := context.Background()
	service, err := userProfileServiceRegistry.Get(ctx)
	require.NoError(t, err)

	id := "test_id"
	username := "test_username"
	email := "test_email"
	items := []string{"0", "1", "2"}
	address := "test_address"
	expected_profile := workflow.UserProfile{ID: id, Username: username, Email: email,
		Items: items, Address: address}

	// Add user
	{
		err := service.UpdateUserProfile(ctx, expected_profile)
		require.NoError(t, err)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestUserProfileServiceGetUserItemIds(t *testing.T) {
	ctx := context.Background()
	service, err := userProfileServiceRegistry.Get(ctx)
	require.NoError(t, err)

	fake_username := "fake"
	id := "test_id"
	username := "test_username"
	email := "test_email"
	expected_items := []string{"0", "1", "2"}
	address := "test_address"
	profile := workflow.UserProfile{ID: id, Username: username, Email: email,
		Items: expected_items, Address: address}

	// Add user
	{
		err := service.UpdateUserProfile(ctx, profile)
		require.NoError(t, err)
	}

	// Get user items
	{
		actual_items, err := service.GetUserItemIds(ctx, username)
		require.NoError(t, err)
		require.ElementsMatch(t, expected_items, actual_items)
	}

	// Try to get item ids from user that doesn't exist
	{
		_, err := service.GetUserItemIds(ctx, fake_username)
		require.Error(t, err)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}

func TestUserProfileServiceAddItem(t *testing.T) {
	ctx := context.Background()
	service, err := userProfileServiceRegistry.Get(ctx)
	require.NoError(t, err)

	fake_username := "fake"
	id := "test_id"
	username := "test_username"
	email := "test_email"
	items := []string{"0", "1", "2"}
	new_items := []string{"3", "4"}
	exptected_items := append(items, new_items...)
	address := "test_address"
	profile := workflow.UserProfile{ID: id, Username: username, Email: email,
		Items: items, Address: address}

	// Add user
	{
		err := service.UpdateUserProfile(ctx, profile)
		require.NoError(t, err)
	}

	// Add items
	{
		err := service.AddItem(ctx, username, new_items[0])
		require.NoError(t, err)

		err = service.AddItem(ctx, username, new_items[1])
		require.NoError(t, err)
	}

	// Check if items were added properly
	{
		actual_items, err := service.GetUserItemIds(ctx, username)
		require.NoError(t, err)
		require.ElementsMatch(t, exptected_items, actual_items)
	}

	// Try to add items to user that doesn't exist
	{
		err := service.AddItem(ctx, fake_username, items[0])
		require.Error(t, err)
	}

	err = service.Cleanup(ctx)
	require.NoError(t, err)
}
