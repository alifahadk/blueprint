
package e2etests

import (
	"context"
	"blueprint/testclients/clients"
	"github.com/blueprint-uservices/cldrel_course/luggageshare/workflow"
)

// Auto-generated code by the Blueprint gotests plugin.
func init() {
	// Initialize the clientlib early so that it can pick up command-line flags
	clientlib := clients.NewClientLibrary("tests")

	
	reviewServiceRegistry.Register("review_service", func(ctx context.Context) (workflow.ReviewService, error) {
		// Build the client library
		namespace, err := clientlib.Build(ctx)
		if err != nil {
			return nil, err
		}

		// Get and return the client
		var client workflow.ReviewService
		err = namespace.Get("review_service.client", &client)
		return client, err
	})
	
	reservationServiceRegistry.Register("reserv_service", func(ctx context.Context) (workflow.ReservationService, error) {
		// Build the client library
		namespace, err := clientlib.Build(ctx)
		if err != nil {
			return nil, err
		}

		// Get and return the client
		var client workflow.ReservationService
		err = namespace.Get("reserv_service.client", &client)
		return client, err
	})
	
	searchServiceRegistry.Register("search_service", func(ctx context.Context) (workflow.SearchService, error) {
		// Build the client library
		namespace, err := clientlib.Build(ctx)
		if err != nil {
			return nil, err
		}

		// Get and return the client
		var client workflow.SearchService
		err = namespace.Get("search_service.client", &client)
		return client, err
	})
	
	frontendServiceRegistry.Register("frontend_service", func(ctx context.Context) (workflow.FrontendService, error) {
		// Build the client library
		namespace, err := clientlib.Build(ctx)
		if err != nil {
			return nil, err
		}

		// Get and return the client
		var client workflow.FrontendService
		err = namespace.Get("frontend_service.client", &client)
		return client, err
	})
	
	luggageServiceRegistry.Register("luggage_service", func(ctx context.Context) (workflow.LuggageService, error) {
		// Build the client library
		namespace, err := clientlib.Build(ctx)
		if err != nil {
			return nil, err
		}

		// Get and return the client
		var client workflow.LuggageService
		err = namespace.Get("luggage_service.client", &client)
		return client, err
	})
	
	userProfileServiceRegistry.Register("user_service", func(ctx context.Context) (workflow.UserProfileService, error) {
		// Build the client library
		namespace, err := clientlib.Build(ctx)
		if err != nil {
			return nil, err
		}

		// Get and return the client
		var client workflow.UserProfileService
		err = namespace.Get("user_service.client", &client)
		return client, err
	})
	
}
